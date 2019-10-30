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

package main

import (
	"github.com/nlnwa/veidemannctl/src/cmd"
	v "github.com/nlnwa/veidemannctl/src/version"
	log "github.com/sirupsen/logrus"
	stdlogger "log"
)

var version = "master"

//go:generate go run -tags=dev bindata/assets_generate.go
func main() {
	v.Version.SetGitVersion(version)
	cmd.Execute()
}

func init() {
	stdlogger.SetOutput(log.StandardLogger().Writer())

	log.SetLevel(log.WarnLevel)
}
