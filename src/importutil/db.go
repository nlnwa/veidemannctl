// Copyright © 2019 National Library of Norway
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package importutil

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/pb"
	configV1 "github.com/nlnwa/veidemann-api-go/config/v1"
	log "github.com/sirupsen/logrus"
	"io"
	"net/url"
	"os"
	"sync"
	"time"
)

type ExistsCode int

const (
	ERROR ExistsCode = iota
	NEW
	DUPLICATE_NEW
	EXISTS_VEIDEMANN
	DUPLICATE_VEIDEMANN
)

func (e ExistsCode) ExistsInVeidemann() bool {
	return e == EXISTS_VEIDEMANN || e == DUPLICATE_VEIDEMANN
}

func (e ExistsCode) String() string {
	names := [...]string{
		"ERROR",
		"NEW",
		"DUPLICATE_NEW",
		"EXISTS_VEIDEMANN",
		"DUPLICATE_VEIDEMANN",}
	if e < ERROR || e > DUPLICATE_VEIDEMANN {
		return "UNKNOWN"
	}
	return names[e]
}

type ExistsResponse struct {
	Code     ExistsCode
	KnownIds []string
}

type ImportDb struct {
	db     *badger.DB
	gc     *time.Ticker
	client configV1.ConfigClient
	vmMux  sync.Mutex
}

func NewImportDb(client configV1.ConfigClient, dbDir string, resetDb bool) *ImportDb {
	opts := badger.DefaultOptions
	opts.Dir = dbDir
	opts.ValueDir = dbDir
	opts.Logger = log.StandardLogger()
	if resetDb {
		os.RemoveAll(dbDir)
	}
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}

	d := &ImportDb{
		db:     db,
		client: client,
	}

	d.gc = time.NewTicker(5 * time.Minute)
	go func() {
		for range d.gc.C {
		again:
			err := db.RunValueLogGC(0.7)
			if err == nil {
				goto again
			}
		}
	}()

	return d
}

func (d *ImportDb) ImportExisting() {
	req := &configV1.ListRequest{
		Kind: configV1.Kind_seed,
	}
	r, err := d.client.ListConfigObjects(context.Background(), req)
	if err != nil {
		log.Fatalf("Error reading seeds from Veidemann: %v", err)
	}

	i := 0
	start := time.Now()
	for {
		msg, err := r.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading seed from Veidemann: %v", err)
		}

		uri, err := url.Parse(msg.GetMeta().GetName())
		if err != nil {
			log.Warnf("error parsing uri '%v': %v", uri, err)
			continue
		}
		exists := d.contains(uri, msg.Id, true)
		switch exists.Code {
		case ERROR:
			log.Infof("Failed handling: %v", uri)
		case NEW:
			log.Debugf("New uri '%v' added", uri)
		case EXISTS_VEIDEMANN:
			log.Debugf("Already exists in Veidemann: %v", uri)
		case DUPLICATE_NEW:
			log.Infof("Found new duplicate: %v", uri)
		case DUPLICATE_VEIDEMANN:
			log.Infof("Found duplicate already existing in Veidemann: %v", uri)
		}
		i++
	}
	elapsed := time.Since(start)
	fmt.Printf("Imported %v seeds from Veidemann in %s\n", i, elapsed)
}

func (d *ImportDb) DuplicateReport() error {
	stream := d.db.NewStream()

	// -- Optional settings
	stream.NumGo = 16                     // Set number of goroutines to use for iteration.
	stream.LogPrefix = "Badger.Streaming" // For identifying stream logs. Outputs to Logger.
	// -- End of optional settings.

	// Send is called serially, while Stream.Orchestrate is running.
	stream.Send = func(list *pb.KVList) error {
		if len(list.GetKv()) != 1 {
			for _, kv := range list.GetKv() {
				val := d.bytesToStringArray(kv.Value)
				if len(val) > 1 {
					for idx, id := range val {
						ref := &configV1.ConfigRef{Id: id, Kind: configV1.Kind_seed}
						seed, err := d.client.GetConfigObject(context.Background(), ref)
						if idx == 0 {
							fmt.Println(string(kv.Key))
						}
						if err == nil {
							fmt.Println("  *", seed.GetId(), seed.GetSeed().GetEntityRef(), seed.GetMeta().GetName())
							entity, err := d.client.GetConfigObject(context.Background(), seed.GetSeed().GetEntityRef())
							if err == nil {
								fmt.Println("       -", entity.GetId(), entity.GetMeta().GetName(), entity.GetMeta().GetDescription())
							}
						}
					}
				}
			}
		}
		return nil
	}

	// Run the stream
	if err := stream.Orchestrate(context.Background()); err != nil {
		return err
	}
	return nil
}

func (d *ImportDb) Check(uri string) (*ExistsResponse, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return &ExistsResponse{Code: ERROR}, fmt.Errorf("error parsing uri '%v': %v", uri, err)
	}
	return d.contains(u, "", false), nil
}

func (d *ImportDb) CheckAndUpdate(uri string) (*ExistsResponse, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return &ExistsResponse{Code: ERROR}, fmt.Errorf("error parsing uri '%v': %v", uri, err)
	}
	return d.contains(u, "", true), nil
}

func (d *ImportDb) CheckAndUpdateVeidemann(uri string, data interface{},
	createFunc func(client configV1.ConfigClient, data interface{}) (id string, err error)) (*ExistsResponse, error) {

	d.vmMux.Lock()
	defer d.vmMux.Unlock()

	u, err := url.Parse(uri)
	if err != nil {
		return &ExistsResponse{Code: ERROR}, fmt.Errorf("error parsing uri '%v': %v", uri, err)
	}

	exists := d.contains(u, "", true)
	if !exists.Code.ExistsInVeidemann() {
		id, err := createFunc(d.client, data)
		if err != nil {
			return exists, err
		}
		d.contains(u, id, true)
	}
	return exists, nil
}

func (d *ImportDb) contains(uri *url.URL, id string, update bool) (response *ExistsResponse) {
	response = &ExistsResponse{}
	var err error
	for {
		err = d.db.Update(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(uri.Host))
			if err == badger.ErrKeyNotFound {
				response.Code = NEW
				if id != "" {
					response.KnownIds = []string{id}
				}
				if update {
					v := d.stringArrayToBytes(response.KnownIds)
					txn.Set([]byte(uri.Host), v)
				}
				return nil
			}
			var val []string
			if err == nil {
				err := item.Value(func(v []byte) error {
					val = d.bytesToStringArray(v)
					if len(val) == 0 {val = nil}
					return nil
				})
				if err != nil {
					return nil
				}

				if id != "" && !d.stringArrayContains(val, id) {
					val = append(val, id)
					if update {
						v := d.stringArrayToBytes(val)
						txn.Set([]byte(uri.Host), v)
					}
				}

				switch len(val) {
				case 0:
					response.Code = DUPLICATE_NEW
				case 1:
					if id == "" || d.stringArrayContains(val, id) {
						response.Code = EXISTS_VEIDEMANN
					} else {
						response.Code = DUPLICATE_VEIDEMANN
					}
				default:
					response.Code = DUPLICATE_VEIDEMANN
				}

				response.KnownIds = val
				return nil
			}
			return nil
		})
		if err != badger.ErrConflict {
			break
		}
	}
	if err != nil {
		log.Fatalf("DB err %v", err)
	}

	return response
}

func (d *ImportDb) Close() {
	for {
		err := d.db.RunValueLogGC(0.7)
		if err != nil {
			break
		}
	}

	d.gc.Stop()
	d.db.Close()
}

func (d *ImportDb) stringArrayToBytes(v []string) []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(v)
	return buf.Bytes()
}

func (d *ImportDb) bytesToStringArray(v []byte) []string {
	buf := bytes.NewBuffer(v)
	strs := []string{}
	gob.NewDecoder(buf).Decode(&strs)
	return strs
}

func (d *ImportDb) stringArrayContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
