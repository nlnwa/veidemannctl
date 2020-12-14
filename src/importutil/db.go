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
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v2"
	"github.com/dgraph-io/badger/v2/pb"
	configV1 "github.com/nlnwa/veidemann-api/go/config/v1"
	"github.com/nlnwa/veidemannctl/src/configutil"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"sync"
	"sync/atomic"
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
		"DUPLICATE_VEIDEMANN"}
	if e < ERROR || e > DUPLICATE_VEIDEMANN {
		return "UNKNOWN"
	}
	return names[e]
}

type ExistsResponse struct {
	NormalizedKey string
	Code          ExistsCode
	KnownIds      []string
}

type ImportDb struct {
	db            *badger.DB
	gc            *time.Ticker
	client        configV1.ConfigClient
	kind          configV1.Kind
	keyNormalizer KeyNormalizer
	vmMux         sync.Mutex
}

type KeyNormalizer interface {
	Normalize(key string) (string, error)
}

type NoopKeyNormalizer struct {
}

func (u *NoopKeyNormalizer) Normalize(s string) (key string, err error) {
	return s, nil
}

func NewImportDb(client configV1.ConfigClient, dbDir string, kind configV1.Kind, keyNormalizer KeyNormalizer, resetDb bool) *ImportDb {
	kindName := kind.String()
	dbDir = path.Join(dbDir, configutil.GlobalFlags.Context, kindName)
	if err := os.MkdirAll(dbDir, 0777); err != nil {
		log.Fatal(err)
	}
	opts := badger.DefaultOptions(dbDir)
	opts.Logger = log.StandardLogger()
	if resetDb {
		if err := os.RemoveAll(dbDir); err != nil {
			log.Fatal(err)
		}
	}
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}

	d := &ImportDb{
		db:            db,
		client:        client,
		kind:          kind,
		keyNormalizer: keyNormalizer,
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
		Kind: d.kind,
	}
	r, err := d.client.ListConfigObjects(context.Background(), req)
	if err != nil {
		log.Fatalf("Error reading %s from Veidemann: %v", d.kind.String(), err)
	}

	var stat struct {
		processed int
		imported  int
		err       int
	}

	start := time.Now()
	for {
		msg, err := r.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading %s from Veidemann: %v", d.kind.String(), err)
		}

		metaName := msg.GetMeta().GetName()

		exists := d.contains(metaName, msg.Id, true)
		switch exists.Code {
		case ERROR:
			stat.err++
			log.Infof("Failed handling: %v", metaName)
		case NEW:
			stat.imported++
			log.Debugf("New key '%v' added", exists.NormalizedKey)
		case EXISTS_VEIDEMANN:
			log.Debugf("Already exists in Veidemann: %v", exists.NormalizedKey)
		case DUPLICATE_NEW:
			log.Infof("Found new duplicate: %v", exists.NormalizedKey)
		case DUPLICATE_VEIDEMANN:
			log.Infof("Found duplicate already existing in Veidemann: %v", exists.NormalizedKey)
		}
		stat.processed++
	}
	elapsed := time.Since(start)
	fmt.Printf("Processed %v %s from Veidemann in %s. %v imported, %v errors\n", stat.processed, d.kind.String(), elapsed, stat.imported, stat.err)
}

type SeedDuplicateReportRecord struct {
	Host  string
	Seeds []SeedRecord
}

type SeedRecord struct {
	SeedId            string
	Uri               string
	SeedDescription   string
	EntityId          string
	EntityName        string
	EntityDescription string
}

func (d *ImportDb) SeedDuplicateReport(w io.Writer) error {
	stream := d.db.NewStream()

	// -- Optional settings
	stream.NumGo = 16                     // Set number of goroutines to use for iteration.
	stream.LogPrefix = "Badger.Streaming" // For identifying stream logs. Outputs to Logger.
	// -- End of optional settings.

	var uniqueDuplicatedNames int32
	var duplicatedNames int32

	// Send is called serially, while Stream.Orchestrate is running.
	stream.Send = func(list *pb.KVList) error {
		for _, kv := range list.GetKv() {
			val := d.bytesToStringArray(kv.Value)
			if len(val) > 1 {
				rec := &SeedDuplicateReportRecord{Host: string(kv.Key)}
				for _, id := range val {
					ref := &configV1.ConfigRef{Id: id, Kind: configV1.Kind_seed}
					seed, err := d.client.GetConfigObject(context.Background(), ref)
					if err == nil {
						sr := SeedRecord{
							SeedId:          seed.GetId(),
							Uri:             seed.GetMeta().GetName(),
							SeedDescription: seed.GetMeta().GetDescription(),
							EntityId:        seed.GetSeed().GetEntityRef().GetId(),
						}
						rec.Seeds = append(rec.Seeds, sr)
						entity, err := d.client.GetConfigObject(context.Background(), seed.GetSeed().GetEntityRef())
						if err == nil {
							sr.EntityName = entity.GetMeta().GetName()
							sr.EntityDescription = entity.GetMeta().GetDescription()
						} else {
							log.Warnf("error getting entity from Veidemann: %v", err)
						}
					} else {
						log.Warnf("error getting seed from Veidemann: %v", err)
					}
				}
				atomic.AddInt32(&uniqueDuplicatedNames, 1)
				atomic.AddInt32(&duplicatedNames, int32(len(val)))
				b, err := json.Marshal(rec)
				if err == nil {
					d.vmMux.Lock()
					if _, err := w.Write(b); err != nil {
						log.Warnf("error wirting record: %v", err)
					}
					if _, err := w.Write([]byte{'\n'}); err != nil {
						log.Warnf("error wirting record: %v", err)
					}
					d.vmMux.Unlock()
				} else {
					log.Warnf("error formatting json: %v", err)
				}
			}
		}
		return nil
	}

	// Run the stream
	if err := stream.Orchestrate(context.Background()); err != nil {
		return err
	}
	fmt.Printf("%v seed urls exist more than once. Total duplicates: %v.\n", uniqueDuplicatedNames, duplicatedNames)
	return nil
}

type CrawlEntityDuplicateReportRecord struct {
	Name          string
	CrawlEntities []CrawlEntityRecord
}

type CrawlEntityRecord struct {
	Id string
}

func (d *ImportDb) CrawlEntityDuplicateReport(w io.Writer) error {
	stream := d.db.NewStream()

	// -- Optional settings
	stream.NumGo = 16                     // Set number of goroutines to use for iteration.
	stream.LogPrefix = "Badger.Streaming" // For identifying stream logs. Outputs to Logger.
	// -- End of optional settings.

	var uniqueDuplicatedNames int32
	var duplicatedNames int32

	// Send is called serially, while Stream.Orchestrate is running.
	stream.Send = func(list *pb.KVList) error {
		for _, kv := range list.GetKv() {
			val := d.bytesToStringArray(kv.Value)
			if len(val) > 1 {
				rec := &CrawlEntityDuplicateReportRecord{Name: string(kv.Key)}
				for _, id := range val {
					sr := CrawlEntityRecord{
						Id: id,
					}
					rec.CrawlEntities = append(rec.CrawlEntities, sr)
				}
				atomic.AddInt32(&uniqueDuplicatedNames, 1)
				atomic.AddInt32(&duplicatedNames, int32(len(val)))
				b, err := json.Marshal(rec)
				if err == nil {
					d.vmMux.Lock()
					if _, err := w.Write(b); err != nil {
						log.Warnf("error wirting record: %v", err)
					}
					if _, err := w.Write([]byte{'\n'}); err != nil {
						log.Warnf("error wirting record: %v", err)
					}
					d.vmMux.Unlock()
				} else {
					log.Warnf("error formatting json: %v", err)
				}
			}
		}
		return nil
	}

	// Run the stream
	if err := stream.Orchestrate(context.Background()); err != nil {
		return err
	}
	fmt.Printf("%v crawl entities names exist more than once. Total duplicates: %v.\n", uniqueDuplicatedNames, duplicatedNames)
	return nil
}

func (d *ImportDb) Check(metaName string) (*ExistsResponse, error) {
	return d.contains(metaName, "", false), nil
}

func (d *ImportDb) CheckAndUpdate(metaName string) (*ExistsResponse, error) {
	return d.contains(metaName, "", true), nil
}

func (d *ImportDb) CheckAndUpdateVeidemann(metaName string, data interface{},
	createFunc func(client configV1.ConfigClient, data interface{}) (id string, err error)) (*ExistsResponse, error) {

	d.vmMux.Lock()
	defer d.vmMux.Unlock()

	exists := d.contains(metaName, "", true)
	if !exists.Code.ExistsInVeidemann() {
		id, err := createFunc(d.client, data)
		if err != nil {
			return exists, err
		}
		d.contains(metaName, id, true)
	}
	return exists, nil
}

func (d *ImportDb) contains(metaName, id string, update bool) (response *ExistsResponse) {
	response = &ExistsResponse{}

	var err error
	response.NormalizedKey, err = d.keyNormalizer.Normalize(metaName)
	if err != nil {
		response.Code = ERROR
		return
	}

	for {
		err = d.db.Update(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(response.NormalizedKey))
			if err == badger.ErrKeyNotFound {
				response.Code = NEW
				if id != "" {
					response.KnownIds = []string{id}
				}
				if update {
					v := d.stringArrayToBytes(response.KnownIds)
					_ = txn.Set([]byte(response.NormalizedKey), v)
				}
				return nil
			}
			var val []string
			if err == nil {
				err := item.Value(func(v []byte) error {
					val = d.bytesToStringArray(v)
					if len(val) == 0 {
						val = nil
					}
					return nil
				})
				if err != nil {
					return nil
				}

				if id != "" && !d.stringArrayContains(val, id) {
					val = append(val, id)
					if update {
						v := d.stringArrayToBytes(val)
						_ = txn.Set([]byte(response.NormalizedKey), v)
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
	_ = d.db.Close()
}

func (d *ImportDb) stringArrayToBytes(v []string) []byte {
	buf := &bytes.Buffer{}
	if err := gob.NewEncoder(buf).Encode(v); err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

func (d *ImportDb) bytesToStringArray(v []byte) []string {
	buf := bytes.NewBuffer(v)
	strs := []string{}
	if err := gob.NewDecoder(buf).Decode(&strs); err != nil && err != io.EOF {
		log.Fatal(err)
	}
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
