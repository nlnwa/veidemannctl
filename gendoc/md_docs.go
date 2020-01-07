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

package main

import (
	"fmt"
	"github.com/nlnwa/veidemannctl/src/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra/doc"
	"path/filepath"
	"runtime"
)

func main() {
	// Find 'this' directory relative to this file to allow callers to be in any package
	var _, b, _, _ = runtime.Caller(0)
	var dir = filepath.Join(filepath.Dir(b), "../docs")
	fmt.Println("generating documentation")
	err := doc.GenMarkdownTree(cmd.RootCmd, dir)
	if err != nil {
		log.Fatal(err)
	}
}
