/*
 * Copyright 2020 National Library of Norway.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package importcmd

import (
	"testing"
)

func Test_importer_normalizeUri(t *testing.T) {
	tests := []struct {
		name           string
		seedDescriptor *seedDesc
		wantUri        string
		wantDedupKey   string
		wantErr        bool
	}{
		{"1", &seedDesc{Uri: "http://www.example.com"}, "http://www.example.com/", "http://www.example.com/", false},
		{"2", &seedDesc{Uri: "http://www.example.com#hash"}, "http://www.example.com/", "http://www.example.com/", false},
		{"3", &seedDesc{Uri: "http://www.example.com?query"}, "http://www.example.com/?query", "http://www.example.com/?query", false},
		{"4", &seedDesc{Uri: "http://www.example.com?query#hash"}, "http://www.example.com/?query", "http://www.example.com/?query", false},
		{"5", &seedDesc{Uri: "http://www.example.com/"}, "http://www.example.com/", "http://www.example.com/", false},
		{"6", &seedDesc{Uri: "http://www.example.com/#hash"}, "http://www.example.com/", "http://www.example.com/", false},
		{"7", &seedDesc{Uri: "http://www.example.com/?query"}, "http://www.example.com/?query", "http://www.example.com/?query", false},
		{"8", &seedDesc{Uri: "http://www.example.com/?query#hash"}, "http://www.example.com/?query", "http://www.example.com/?query", false},
		{"9", &seedDesc{Uri: "http://www.example.com/foo/bar"}, "http://www.example.com/foo/bar", "http://www.example.com/foo/bar", false},
		{"10", &seedDesc{Uri: "http://www.example.com/foo/bar#hash"}, "http://www.example.com/foo/bar", "http://www.example.com/foo/bar", false},
		{"11", &seedDesc{Uri: "http://www.example.com/foo/bar?query"}, "http://www.example.com/foo/bar?query", "http://www.example.com/foo/bar?query", false},
		{"12", &seedDesc{Uri: "http://www.example.com/foo/bar?query#hash"}, "http://www.example.com/foo/bar?query", "http://www.example.com/foo/bar?query", false},
		{"13", &seedDesc{Uri: "https://www.example.com/foo/bar#hash"}, "https://www.example.com/foo/bar", "https://www.example.com/foo/bar", false},
		{"14", &seedDesc{Uri: "HTTPS://www.example.com/foo/bar#hash"}, "https://www.example.com/foo/bar", "https://www.example.com/foo/bar", false},
		{"15", &seedDesc{Uri: "https://www.Example.Com/foo/bar#hash"}, "https://www.example.com/foo/bar", "https://www.example.com/foo/bar", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &importer{}
			importFlags.toplevel = false
			if err := i.normalizeUri(tt.seedDescriptor); (err != nil) != tt.wantErr {
				t.Errorf("normalizeUri() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.seedDescriptor.Uri != tt.wantUri {
				t.Errorf("normalizeUri() uri = %v, want %v", tt.seedDescriptor.Uri, tt.wantUri)
			}
		})
	}
}

func Test_UriKeyNormalizer_toplevelFlag(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		wantKey string
		wantErr bool
	}{
		{"1", "http://www.example.com", "http://www.example.com/", false},
		{"2", "http://www.example.com#hash", "http://www.example.com/", false},
		{"3", "http://www.example.com?query", "http://www.example.com/", false},
		{"4", "http://www.example.com?query#hash", "http://www.example.com/", false},
		{"5", "http://www.example.com/", "http://www.example.com/", false},
		{"6", "http://www.example.com/#hash", "http://www.example.com/", false},
		{"7", "http://www.example.com/?query", "http://www.example.com/", false},
		{"8", "http://www.example.com/?query#hash", "http://www.example.com/", false},
		{"9", "http://www.example.com/foo/bar", "http://www.example.com/", false},
		{"10", "http://www.example.com/foo/bar#hash", "http://www.example.com/", false},
		{"11", "http://www.example.com/foo/bar?query", "http://www.example.com/", false},
		{"12", "http://www.example.com/foo/bar?query#hash", "http://www.example.com/", false},
		{"13", "https://www.example.com/foo/bar#hash", "https://www.example.com/", false},
		{"14", "HTTPS://www.example.com/foo/bar#hash", "https://www.example.com/", false},
		{"15", "https://www.Example.Com/foo/bar#hash", "https://www.example.com/", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &UriKeyNormalizer{toplevel: true, ignoreScheme: false}
			key, err := i.Normalize(tt.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("Normalize() error = %v, wantErr %v", err, tt.wantErr)
			}
			if key != tt.wantKey {
				t.Errorf("Normalize() key = %v, want %v", key, tt.wantKey)
			}
		})
	}
}

func Test_UriKeyNormalizer_ignoreSchemeFlag(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		wantKey string
		wantErr bool
	}{
		{"1", "http://www.example.com", "//www.example.com/", false},
		{"2", "http://www.example.com#hash", "//www.example.com/", false},
		{"3", "http://www.example.com?query", "//www.example.com/?query", false},
		{"4", "http://www.example.com?query#hash", "//www.example.com/?query", false},
		{"5", "http://www.example.com/", "//www.example.com/", false},
		{"6", "http://www.example.com/#hash", "//www.example.com/", false},
		{"7", "http://www.example.com/?query", "//www.example.com/?query", false},
		{"8", "http://www.example.com/?query#hash", "//www.example.com/?query", false},
		{"9", "http://www.example.com/foo/bar", "//www.example.com/foo/bar", false},
		{"10", "http://www.example.com/foo/bar#hash", "//www.example.com/foo/bar", false},
		{"11", "http://www.example.com/foo/bar?query", "//www.example.com/foo/bar?query", false},
		{"12", "http://www.example.com/foo/bar?query#hash", "//www.example.com/foo/bar?query", false},
		{"13", "https://www.example.com/foo/bar#hash", "//www.example.com/foo/bar", false},
		{"14", "HTTPS://www.example.com/foo/bar#hash", "//www.example.com/foo/bar", false},
		{"15", "https://www.Example.Com/foo/bar#hash", "//www.example.com/foo/bar", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &UriKeyNormalizer{toplevel: false, ignoreScheme: true}
			key, err := i.Normalize(tt.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("Normalize() error = %v, wantErr %v", err, tt.wantErr)
			}
			if key != tt.wantKey {
				t.Errorf("Normalize() key = %v, want %v", key, tt.wantKey)
			}
		})
	}
}
