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

package importutil

import (
	"errors"
	"net/url"
	"strings"
)

type UriKeyNormalizer struct {
	// Toplevel will normalize the uri to the top level domain
	Toplevel bool
	// IgnoreScheme will ignore the scheme when normalizing
	IgnoreScheme bool
}

func (u *UriKeyNormalizer) Normalize(s string) (string, error) {
	uri, err := url.Parse(s)
	if err != nil {
		return "", err
	}

	uri.Fragment = ""
	uri.Host = strings.ToLower(uri.Host)

	if uri.Hostname() == "" {
		return "", errors.New("missing hostname")
	}

	if u.Toplevel {
		uri.Path = "/"
		uri.RawQuery = ""
	}

	if uri.Path == "" {
		uri.Path = "/"
	}

	if u.IgnoreScheme {
		uri.Scheme = ""
	}

	return uri.String(), nil
}
