// Copyright © 2017 National Library of Norway.
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

package connection

import (
	"crypto/x509"
	"github.com/golang/protobuf/ptypes/empty"
	configV1 "github.com/nlnwa/veidemann-api-go/config/v1"
	controllerV1 "github.com/nlnwa/veidemann-api-go/controller/v1"
	reportV1 "github.com/nlnwa/veidemann-api-go/report/v1"
	"github.com/nlnwa/veidemannctl/src/configutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"strings"
	"time"
)

func NewControllerClient() (controllerV1.ControllerClient, *grpc.ClientConn) {
	conn := newConnection()
	c := controllerV1.NewControllerClient(conn)
	return c, conn
}

func NewReportClient() (reportV1.ReportClient, *grpc.ClientConn) {
	conn := newConnection()
	c := reportV1.NewReportClient(conn)
	return c, conn
}

func NewConfigClient() (configV1.ConfigClient, *grpc.ClientConn) {
	conn := newConnection()
	c := configV1.NewConfigClient(conn)
	return c, conn
}

func newConnection() *grpc.ClientConn {
	idp, tls := GetIdp()

	// Set up a connection to the server.
	conn, _ := connect(idp, tls)
	return conn
}

func connect(idp string, tls bool) (*grpc.ClientConn, bool) {
	address := configutil.GlobalFlags.ControllerAddress
	apiKey := configutil.GlobalFlags.ApiKey

	dialOptions := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	}

	if idp != "" {
		dialOptions = addCredentials(idp, apiKey, dialOptions)
	}

	// Set up a connection to the server.
	creds := clientTransportCredentials(tls)
	log.Debugf("connecting to %v", address)
	conn, err := blockingDial(context.Background(), address, creds, dialOptions...)
	if err != nil {
		if strings.Contains(err.Error(), "first record does not look like a TLS handshake") {
			log.Debug("Could not connect with TLS, retrying without credentials")
			tls = false
			conn, err = blockingDial(context.Background(), address, nil, dialOptions...)
			if err != nil {
				log.Fatalf("Could not connect: %v", err)
			}
			log.Warn("Connected with insecure transport. Server does not support TLS")
		} else {
			log.Fatalf("Could not connect: %v", err)
		}
	}
	return conn, tls
}

func clientTransportCredentials(tls bool) credentials.TransportCredentials {
	var creds credentials.TransportCredentials

	if tls {
		certPool, err := x509.SystemCertPool()
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}

		// Add CAs from config
		if viper.GetString("rootCAs") != "" {
			rootCABytes := []byte(viper.GetString("rootCAs"))
			if !certPool.AppendCertsFromPEM(rootCABytes) {
				log.Warn("no certs found in root CA file")
			}
		}

		serverNameOverride := configutil.GlobalFlags.ServerNameOverride
		log.Debugf("Using server name override: %s", serverNameOverride)
		creds = credentials.NewClientTLSFromCert(certPool, serverNameOverride)
	}

	return creds
}

func GetIdp() (string, bool) {
	conn, tls := connect("", true)
	defer conn.Close()

	c := controllerV1.NewControllerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Debug("requesting OpenIdConnectIssuer")
	reply, err := c.GetOpenIdConnectIssuer(ctx, &empty.Empty{})
	if err != nil {
		log.Fatalf("Could not get idp address: %v", err)
	}

	idp := reply.GetOpenIdConnectIssuer()
	log.Debugf("using OpenIdConnectIssuer %v", idp)
	if idp == "" {
		log.Warn("Server was not configured with an idp. Proceeding without authentication")
	}

	return idp, tls
}

type bearerTokenCred struct {
	tokenType string
	token     string
}

func (b bearerTokenCred) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": b.tokenType + " " + b.token,
	}, nil
}

func (b bearerTokenCred) RequireTransportSecurity() bool {
	return false
}

func addCredentials(idp, apiKey string, opts []grpc.DialOption) []grpc.DialOption {
	if apiKey != "" {
		at := &bearerTokenCred{"ApiKey", apiKey}
		opts = append(opts, grpc.WithPerRPCCredentials(at))
		return opts
	}

	a := NewAuth(idp)
	if !a.enabled {
		return opts
	}

	a.CheckStoredAccessToken()
	if a.oauth2Token != nil {
		log.Debugf("Raw IdToken: %s", a.oauth2Token.TokenType)
	}
	if a.rawIdToken == "" {
		return opts
	}

	bt := &bearerTokenCred{a.oauth2Token.TokenType, a.rawIdToken}
	opts = append(opts, grpc.WithPerRPCCredentials(bt))
	return opts
}

// BlockingDial is a helper method to dial the given address, using optional TLS credentials,
// and blocking until the returned connection is ready. If the given credentials are nil, the
// connection will be insecure (plain-text).
// This function is borrowed from https://github.com/fullstorydev/grpcurl/blob/master/grpcurl.go
func blockingDial(ctx context.Context, address string, creds credentials.TransportCredentials, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	// grpc.Dial doesn't provide any information on permanent connection errors (like
	// TLS handshake failures). So in order to provide good error messages, we need a
	// custom dialer that can provide that info. That means we manage the TLS handshake.
	result := make(chan interface{}, 1)

	writeResult := func(res interface{}) {
		// non-blocking write: we only need the first result
		select {
		case result <- res:
		default:
		}
	}

	dialer := func(address string, timeout time.Duration) (net.Conn, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		conn, err := (&net.Dialer{Cancel: ctx.Done()}).Dial("tcp", address)
		if err != nil {
			writeResult(err)
			return nil, err
		}
		if creds != nil {
			conn, _, err = creds.ClientHandshake(ctx, address, conn)
			if err != nil {
				writeResult(err)
				return nil, err
			}
		}
		return conn, nil
	}

	// Even with grpc.FailOnNonTempDialError, this call will usually timeout in
	// the face of TLS handshake errors. So we can't rely on grpc.WithBlock() to
	// know when we're done. So we run it in a goroutine and then use result
	// channel to either get the channel or fail-fast.
	go func() {
		opts = append(opts,
			grpc.WithBlock(),
			grpc.FailOnNonTempDialError(true),
			grpc.WithDialer(dialer),
			grpc.WithInsecure(), // we are handling TLS, so tell grpc not to
		)
		conn, err := grpc.DialContext(ctx, address, opts...)
		var res interface{}
		if err != nil {
			res = err
		} else {
			res = conn
		}
		writeResult(res)
	}()

	select {
	case res := <-result:
		if conn, ok := res.(*grpc.ClientConn); ok {
			return conn, nil
		} else {
			return nil, res.(error)
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
