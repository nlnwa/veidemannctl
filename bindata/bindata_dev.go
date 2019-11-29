// +build dev

/*
 * Copyright 2019 National Library of Norway.
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

package bindata

import (
	"net/http"
	"path/filepath"
	"runtime"
)

// Find the 'res' directory relative to this file to allow callers to be in any package
var _, b, _, _ = runtime.Caller(0)
var dir = filepath.Join(filepath.Dir(b), "../res")

// Assets contains project assets.
var Assets http.FileSystem = http.Dir(dir)
