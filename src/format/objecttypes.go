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
	"fmt"
	configV1 "github.com/nlnwa/veidemann-api-go/config/v1"
	"sort"
	"strings"
)

//	{"crawlEntity", []string{"Id", "Meta.Name", "Meta.Description"}},
//	{"seed", []string{"Id", "Meta.Name", "Spec.Seed.EntityRef", "Spec.Seed.Scope.SurtPrefix", "Spec.Seed.JobRef", "Spec.Seed.Disabled"}},
//	{"crawlConfig", []string{"Id", "Meta.Name", "Meta.Description", "Spec.CrawlConfig.CollectionRef", "Spec.CrawlConfig.BrowserConfigRef", "Spec.CrawlConfig.PolitenessRef", "Spec.CrawlConfig.Extra"}},
//	{"crawlJob", []string{"Id", "Meta.Name", "Meta.Description", "Spec.CrawlJob.ScheduleRef", "Spec.CrawlJob.Limits", "Spec.CrawlJob.CrawlConfigRef", "Spec.CrawlJob.Disabled"}},
//	{"crawlScheduleConfig", []string{"Id", "Meta.Name", "Meta.Description", "Spec.CrawlScheduleConfig.CronExpression", "Spec.CrawlScheduleConfig.ValidFrom", "Spec.CrawlScheduleConfig.ValidTo"}},
//	{"browserConfig", []string{"Id", "Meta.Name", "Meta.Description", "Spec.BrowserConfig.UserAgent", "Spec.BrowserConfig.WindowWidth", "Spec.BrowserConfig.WindowHeight", "Spec.BrowserConfig.PageLoadTimeoutMs", "Spec.BrowserConfig.MaxInactivityTimeMs"}},
//	{"politenessConfig", []string{"Id", "Meta.Name", "Meta.Description", "Spec.PolitenessConfig.RobotsPolicy", "Spec.PolitenessConfig.MinTimeBetweenPageLoadMs", "Spec.PolitenessConfig.MaxTimeBetweenPageLoadMs", "Spec.PolitenessConfig.DelayFactor", "Spec.PolitenessConfig.MaxRetries", "Spec.PolitenessConfig.RetryDelaySeconds", "Spec.PolitenessConfig.CrawlHostGroupSelector"}},
//	{"browserScript", []string{"Id", "Meta.Name", "Meta.Description", "Spec.BrowserScript.Script", "Spec.BrowserScript.UrlRegexp"}},
//	{"crawlHostGroupConfig", []string{"Id", "Meta.Name", "Meta.Description", "Spec.CrawlHostGroup.IpRange"}},
//	{"roleMapping", []string{"Id", "Spec.RoleMapping.EmailOrGroup.Email", "Spec.RoleMapping.EmailOrGroup.Group", "Spec.RoleMapping.Role"}},
//	{"collection", []string{"Id", "Meta.Name", "Spec.Collection.CollectionDedupPolicy", "Spec.Collection.FileRotationPolicy", "Spec.Collection.SubCollections"}},

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

func GetObjectNames() []string {
	result := make([]string, len(configV1.Kind_name)-1)
	idx := 0
	for _, n := range configV1.Kind_name {
		if n != "undefined" {
			result[idx] = n
			idx++
		}
	}
	sort.Strings(result)
	return result
}

func GetTemplateBaseName(obj interface{}) string {
	return fmt.Sprintf("%T", obj)
}
