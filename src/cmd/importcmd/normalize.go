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
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type UriKeyNormalizer struct {
	toplevel     bool
	ignoreScheme bool
}

func (u *UriKeyNormalizer) Normalize(s string) (key string, err error) {
	uri, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("unparseable URL '%v', cause: %v", s, err)
	}

	if uri.Host == "" {
		return "", errors.New("unparseable URL")
	}

	uri.Fragment = ""
	uri.Host = strings.ToLower(uri.Host)

	if u.toplevel {
		uri.Path = "/"
		uri.RawQuery = ""
	}

	if uri.Path == "" {
		uri.Path = "/"
	}

	if u.ignoreScheme {
		uri.Scheme = ""
	}

	key = uri.String()

	return
}
