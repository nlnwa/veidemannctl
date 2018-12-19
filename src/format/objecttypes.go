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

package format

import (
	configV1 "github.com/nlnwa/veidemann-api-go/config/v1"
	"strings"
)

var objectTypes = []struct {
	vName  string
	tabDef []string
}{
	{"crawlEntity", []string{"Id", "Meta.Name", "Meta.Description"}},
	{"seed", []string{"Id", "Meta.Name", "Spec.Seed.EntityId", "Spec.Seed.Scope.SurtPrefix", "Spec.Seed.JobId", "Spec.Seed.Disabled"}},
	{"crawlConfig", []string{"Id", "Meta.Name", "Meta.Description", "Spec.CrawlConfig.BrowserConfigId", "Spec.CrawlConfig.PolitenessId", "Spec.CrawlConfig.Extra"}},
	{"crawlJob", []string{"Id", "Meta.Name", "Meta.Description", "Spec.CrawlJob.ScheduleId", "Spec.CrawlJob.Limits", "Spec.CrawlJob.CrawlConfigId", "Spec.CrawlJob.Disabled"}},
	{"crawlScheduleConfig", []string{"Id", "Meta.Name", "Meta.Description", "CronExpression", "ValidFrom", "ValidTo"}},
	{"browserConfig", []string{"Id", "Meta.Name", "Meta.Description", "UserAgent", "WindowWidth", "WindowHeight", "PageLoadTimeoutMs", "SleepAfterPageloadMs"}},
	{"politenessConfig", []string{"Id", "Meta.Name", "Meta.Description", "RobotsPolicy", "MinTimeBetweenPageLoadMs", "MaxTimeBetweenPageLoadMs", "DelayFactor", "MaxRetries", "RetryDelaySeconds", "CrawlHostGroupSelector"}},
	{"browserScript", []string{"Id", "Meta.Name", "Meta.Description", "Script", "UrlRegexp"}},
	{"crawlHostGroup", []string{"Id", "Meta.Name", "Meta.Description", "IpRange"}},
	{"roleMapping", []string{"Id", "Spec.RoleMapping.EmailOrGroup.Email", "Spec.RoleMapping.EmailOrGroup.Group", "Spec.RoleMapping.Role"}},
}

// Get kind for string
func GetKind(Name string) configV1.Kind {
	Name = strings.ToLower(Name)
	for _, k := range configV1.Kind_name {
		if strings.ToLower(k) == Name {
			return configV1.Kind(configV1.Kind_value[k])
		}
	}
	return configV1.Kind_undefined
}

func GetTableDefForKind(kind configV1.Kind) []string {
	for _, ot := range objectTypes {
		if ot.vName == kind.String() {
			return ot.tabDef
		}
	}
	return nil
}

func GetObjectNames() []string {
	result := make([]string, len(objectTypes))
	for idx, ot := range objectTypes {
		result[idx] = ot.vName
	}
	return result
}
