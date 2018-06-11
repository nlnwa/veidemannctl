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

package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nlnwa/veidemannctl/src/apiutil"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/nlnwa/veidemannctl/src/format"
	api "github.com/nlnwa/veidemannctl/veidemann_api"
	"golang.org/x/net/context"
	"github.com/golang/protobuf/proto"
)

var flags struct {
	label      string
	name       string
	file       string
	format     string
	goTemplate string
	pageSize   int32
	page       int32
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [object_type]",
	Short: "Get the value(s) for an object type",
	Long: `Display one or many objects.

` +
		printValidObjectTypes() +
		`Examples:
  #List all seeds.
  veidemannctl get seed

  #List all seeds in yaml output format.
  veidemannctl get seed -f yaml`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			client, conn := connection.NewControllerClient()
			defer conn.Close()

			var selector []string
			var ids []string

			if len(args) == 2 {
				ids = args[1:]
				fmt.Println("ID: ", ids)
			}

			s := &format.MarshalSpec{
				Filename: flags.file,
				Format:   flags.format,
				Template: flags.goTemplate,
			}
			defer s.Close()

			var msg proto.Message

			switch args[0] {
			case "entity":
				request := apiutil.CreateListRequest(ids, flags.name, flags.label, flags.pageSize, flags.page)

				r, err := client.ListCrawlEntities(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get entity: %v", err)
				}

				msg = r
				if format.Marshal(s, msg) != nil {
					os.Exit(1)
				}
			case "seed":
				request := api.SeedListRequest{}

				selector = apiutil.CreateSelector(flags.label)

				request.Id = ids
				request.Name = flags.name
				request.LabelSelector = selector
				request.Page = flags.page
				request.PageSize = flags.pageSize

				r, err := client.ListSeeds(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get seed: %v", err)
				}

				msg = r
				if format.Marshal(s, msg) != nil {
					os.Exit(1)
				}
			case "job":
				request := apiutil.CreateListRequest(ids, flags.name, flags.label, flags.pageSize, flags.page)

				r, err := client.ListCrawlJobs(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get job: %v", err)
				}

				msg = r
				if format.Marshal(s, msg) != nil {
					os.Exit(1)
				}
			case "crawlconfig":
				request := apiutil.CreateListRequest(ids, flags.name, flags.label, flags.pageSize, flags.page)

				r, err := client.ListCrawlConfigs(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get crawl config: %v", err)
				}

				msg = r
				if format.Marshal(s, msg) != nil {
					os.Exit(1)
				}
			case "schedule":
				request := apiutil.CreateListRequest(ids, flags.name, flags.label, flags.pageSize, flags.page)

				r, err := client.ListCrawlScheduleConfigs(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get schedule config: %v", err)
				}

				msg = r
				if format.Marshal(s, msg) != nil {
					os.Exit(1)
				}
			case "browser":
				request := apiutil.CreateListRequest(ids, flags.name, flags.label, flags.pageSize, flags.page)

				r, err := client.ListBrowserConfigs(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get browser config: %v", err)
				}

				msg = r
				if format.Marshal(s, msg) != nil {
					os.Exit(1)
				}
			case "politeness":
				request := apiutil.CreateListRequest(ids, flags.name, flags.label, flags.pageSize, flags.page)

				r, err := client.ListPolitenessConfigs(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get politeness config: %v", err)
				}

				msg = r
				if format.Marshal(s, msg) != nil {
					os.Exit(1)
				}
			case "script":
				request := apiutil.CreateListRequest(ids, flags.name, flags.label, flags.pageSize, flags.page)

				r, err := client.ListBrowserScripts(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get browser script: %v", err)
				}

				msg = r
				if format.Marshal(s, msg) != nil {
					os.Exit(1)
				}
			case "group":
				request := apiutil.CreateListRequest(ids, flags.name, flags.label, flags.pageSize, flags.page)

				r, err := client.ListCrawlHostGroupConfigs(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get crawl host group config: %v", err)
				}

				msg = r
				if format.Marshal(s, msg) != nil {
					os.Exit(1)
				}
			case "loglevel":
				r, err := client.GetLogConfig(context.Background(), &empty.Empty{})
				if err != nil {
					log.Fatalf("could not get log config: %v", err)
				}

				msg = r
				if format.Marshal(s, msg) != nil {
					os.Exit(1)
				}
			case "activerole":
				r, err := client.GetRolesForActiveUser(context.Background(), &empty.Empty{})
				if err != nil {
					log.Fatalf("could not get active role: %v", err)
				}

				msg = r
				if format.Marshal(s, msg) != nil {
					os.Exit(1)
				}
			case "role":
				request := api.RoleMappingsListRequest{}
				request.Page = flags.page
				request.PageSize = flags.pageSize

				r, err := client.ListRoleMappings(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get active role: %v", err)
				}

				msg = r
				if format.Marshal(s, msg) != nil {
					os.Exit(1)
				}
			default:
				fmt.Printf("Unknown object type\n")
				cmd.Usage()
			}
		} else {
			fmt.Print("You must specify the object type to get. ")
			fmt.Println(printValidObjectTypes())
			fmt.Println("See 'veidemannctl get -h' for help")
		}
	},
	ValidArgs: format.GetObjectNames(),
}

func printValidObjectTypes() string {
	var names string
	for _, v := range format.GetObjectNames() {
		names += fmt.Sprintf("  * %s\n", v)
	}
	return fmt.Sprintf("Valid object types include:\n%s\n", names)
}

func init() {
	RootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	getCmd.PersistentFlags().StringVarP(&flags.label, "label", "l", "", "List objects by label (<type>:<value> | <value>)")

	getCmd.PersistentFlags().StringVarP(&flags.name, "name", "n", "", "List objects by name (accepts regular expressions)")
	annotation := make(map[string][]string)
	annotation[cobra.BashCompCustom] = []string{"__veidemannctl_get_name"}

	getCmd.PersistentFlags().Lookup("name").Annotations = annotation
	getCmd.PersistentFlags().StringVarP(&flags.format, "output", "o", "table", "Output format (table|json|yaml|template|template-file)")
	getCmd.PersistentFlags().StringVarP(&flags.goTemplate, "template", "t", "", "A Go template used to format the output")
	getCmd.PersistentFlags().StringVarP(&flags.file, "filename", "f", "", "File name to write to")
	getCmd.PersistentFlags().Int32VarP(&flags.pageSize, "pagesize", "s", 10, "Number of objects to get")
	getCmd.PersistentFlags().Int32VarP(&flags.page, "page", "p", 0, "The page number")
}
