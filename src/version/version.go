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

package version

import (
	"fmt"
	"runtime"
)

type VeidemannVersion struct {
	gitVersion string
}

func (v *VeidemannVersion) SetGitVersion(version string) {
	v.gitVersion = version
}

func (v *VeidemannVersion) GetVersionString() string {
	return fmt.Sprintf("Client version: %s, Go version: %s, Platform: %s/%s\n",
		v.gitVersion,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH)
}

var Version = &VeidemannVersion{}
