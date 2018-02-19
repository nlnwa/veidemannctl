// Copyright © 2017 National Library of Norway
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

package util

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/coreos/go-oidc"
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"math/rand"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

// return an HTTP client which trusts the provided root CAs.
func httpClientForRootCAs() *http.Client {
	tlsConfig := tls.Config{RootCAs: x509.NewCertPool()}
	if viper.GetString("rootCAs") == "" {
		return nil
	}
	rootCABytes := []byte(viper.GetString("rootCAs"))
	if !tlsConfig.RootCAs.AppendCertsFromPEM(rootCABytes) {
		log.Warn("no certs found in root CA file")
		return nil
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tlsConfig,
			Proxy:           http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

type auth struct {
	clientID     string
	clientSecret string
	redirectURI  string
	rawIdToken   string
	idToken      *oidc.IDToken
	oauth2Token  *oauth2.Token

	idTokenVerifier *oidc.IDTokenVerifier
	provider        *oidc.Provider

	// Does the provider use "offline_access" scope to request a refresh token
	// or does it use "access_type=offline" (e.g. Google)?
	offlineAsScope bool

	client  *http.Client
	ctx     context.Context
	state   string
	enabled bool
}

func NewAuth() *auth {
	a := auth{}
	a.enabled = true
	a.offlineAsScope = true
	a.clientID = "veidemann-cli"
	a.clientSecret = "cli-app-secret"
	a.redirectURI = "urn:ietf:wg:oauth:2.0:oob"

	a.client = httpClientForRootCAs()

	if a.client == nil {
		a.client = http.DefaultClient
	}

	a.ctx = oidc.ClientContext(context.Background(), a.client)

	// Initialize a provider by specifying dex's issuer URL.
	idp := viper.GetString("idp")
	p, err := oidc.NewProvider(a.ctx, idp)
	if err != nil {
		log.Warn(fmt.Errorf("Could not connect to authentication server '%s': %v. Proceeding without authentication", idp, err))
		a.enabled = false
	} else {
		a.provider = p
		oc := oidc.Config{ClientID: a.clientID}
		a.idTokenVerifier = a.provider.Verifier(&oc)
	}
	return &a
}

func (a *auth) oauth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     a.clientID,
		ClientSecret: a.clientSecret,
		Endpoint:     a.provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "groups", "offline_access", "audience:server:client_id:veidemann-api"},
		RedirectURL:  a.redirectURI,
	}
}

func (a *auth) CreateAuthCodeURL() string {
	a.state = RandStringBytesMaskImprSrc(16)
	viper.Set("nonce", a.state)

	nonce := oidc.Nonce(a.state)
	return a.oauth2Config().AuthCodeURL(a.state, nonce)
}

func (a *auth) Openbrowser(authCodeURL string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", authCodeURL).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", authCodeURL).Start()
	case "darwin":
		err = exec.Command("open", authCodeURL).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func (a *auth) VerifyCode(code string) {
	oauth2Token, err := a.oauth2Config().Exchange(a.ctx, code)
	if err != nil {
		log.Fatal(err)
	}
	a.oauth2Token = oauth2Token

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Fatal("No token found")
	}
	a.rawIdToken = rawIDToken

	a.verifyIdToken()

	viper.Set("accessToken", marshlAccessToken(oauth2Token))
}

func (a *auth) verifyIdToken() {
	// Parse and verify ID Token payload.
	idToken, err := a.idTokenVerifier.Verify(a.ctx, a.rawIdToken)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("IdToken: %s", idToken)

	if idToken.Nonce != viper.GetString("nonce") {
		log.Fatal("Nonce did not match")
	}
	a.idToken = idToken
}

func (a *auth) CheckStoredAccessToken() {
	accessToken := viper.GetString("accessToken")
	if accessToken == "" {
		log.Debugf("No accessToken")
		return
	}

	a.oauth2Token = unmarshalAccessToken(accessToken)
	log.Debugf("AccessToken: %v", a.oauth2Token)

	// Extract the ID Token from OAuth2 token.
	rawIdToken, ok := a.oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Fatal("No token found")
	}
	a.rawIdToken = rawIdToken

	a.verifyIdToken()
}

// Extract custom claims.
type Claims struct {
	Email    string   `json:"email"`
	Verified bool     `json:"email_verified"`
	Groups   []string `json:"groups"`
	Name     string   `json:"name"`
}

func (a *auth) Claims() (claims Claims) {
	if err := a.idToken.Claims(&claims); err != nil {
		log.Fatal(err)
	}
	return claims
}

type yamlToken struct {
	// AccessToken is the token that authorizes and authenticates
	// the requests.
	AccessToken string

	// TokenType is the type of token.
	// The Type method returns either this or "Bearer", the default.
	TokenType string

	// RefreshToken is a token that's used by the application
	// (as opposed to the user) to refresh the access token
	// if it expires.
	RefreshToken string

	// Expiry is the optional expiration time of the access token.
	//
	// If zero, TokenSource implementations will reuse the same
	// token forever and RefreshToken or equivalent
	// mechanisms for that TokenSource will not be used.
	Expiry time.Time

	// Raw optionally contains extra metadata from the server
	// when updating a token.
	IdToken interface{}
}

func marshlAccessToken(token *oauth2.Token) string {
	tmp := yamlToken{
		AccessToken:  token.AccessToken,
		TokenType:    token.Type(),
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
		IdToken:      token.Extra("id_token"),
	}
	r, err := yaml.Marshal(tmp)
	if err != nil {
		log.Fatalf("Could not marshal Access Token: %v", err)
	}
	s := base64.StdEncoding.EncodeToString(r)
	return s
}

func unmarshalAccessToken(accessToken string) *oauth2.Token {
	s, err := base64.StdEncoding.DecodeString(accessToken)
	if err != nil {
		log.Fatalf("Could not unmarshal Access Token: %v", err)
	}
	var tmp yamlToken
	if err := yaml.Unmarshal(s, &tmp); err != nil {
		log.Fatalf("Could not unmarshal Access Token: %v", err)
	}

	at := &oauth2.Token{
		AccessToken:  tmp.AccessToken,
		TokenType:    tmp.TokenType,
		RefreshToken: tmp.RefreshToken,
		Expiry:       tmp.Expiry,
	}
	extra := map[string]interface{}{"id_token": tmp.IdToken}
	at = at.WithExtra(extra)

	return at
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
