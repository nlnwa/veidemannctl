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
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nlnwa/veidemannctl/util"
	api "github.com/nlnwa/veidemannctl/veidemann_api"
	"golang.org/x/net/context"
)

var (
	label  string
	name   string
	file   string
	format string
)

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
		if len(args) > 0 && len(args) <= 2 {
			client, conn := util.NewControllerClient()
			defer conn.Close()

			var selector []string
			var ids []string

			if len(args) == 2 {
				ids = args[1:]
				fmt.Println("ID: ", ids)
			}

			switch args[0] {
			case "entity":
				request := util.CreateListRequest(ids, name, label, pageSize, page)

				r, err := client.ListCrawlEntities(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get entity: %v", err)
				}

				if util.Marshal(file, format, r) != nil {
					os.Exit(1)
				}
			case "seed":
				request := api.SeedListRequest{}

				selector = util.CreateSelector(label)

				request.Id = ids
				request.Name = name
				request.LabelSelector = selector
				request.Page = page
				request.PageSize = pageSize

				r, err := client.ListSeeds(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get seed: %v", err)
				}

				if util.Marshal(file, format, r) != nil {
					os.Exit(1)
				}
			case "job":
				request := util.CreateListRequest(ids, name, label, pageSize, page)

				r, err := client.ListCrawlJobs(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get job: %v", err)
				}

				if util.Marshal(file, format, r) != nil {
					os.Exit(1)
				}
			case "crawlconfig":
				request := util.CreateListRequest(ids, name, label, pageSize, page)

				r, err := client.ListCrawlConfigs(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get crawl config: %v", err)
				}

				if util.Marshal(file, format, r) != nil {
					os.Exit(1)
				}
			case "schedule":
				request := util.CreateListRequest(ids, name, label, pageSize, page)

				r, err := client.ListCrawlScheduleConfigs(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get schedule config: %v", err)
				}

				if util.Marshal(file, format, r) != nil {
					os.Exit(1)
				}
			case "browser":
				request := util.CreateListRequest(ids, name, label, pageSize, page)

				r, err := client.ListBrowserConfigs(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get browser config: %v", err)
				}

				if util.Marshal(file, format, r) != nil {
					os.Exit(1)
				}
			case "politeness":
				request := util.CreateListRequest(ids, name, label, pageSize, page)

				r, err := client.ListPolitenessConfigs(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get politeness config: %v", err)
				}

				if util.Marshal(file, format, r) != nil {
					os.Exit(1)
				}
			case "script":
				request := util.CreateListRequest(ids, name, label, pageSize, page)

				r, err := client.ListBrowserScripts(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get browser script: %v", err)
				}

				if util.Marshal(file, format, r) != nil {
					os.Exit(1)
				}
			case "group":
				request := util.CreateListRequest(ids, name, label, pageSize, page)

				r, err := client.ListCrawlHostGroupConfigs(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get crawl host group config: %v", err)
				}

				if util.Marshal(file, format, r) != nil {
					os.Exit(1)
				}
			case "loglevel":
				r, err := client.GetLogConfig(context.Background(), &empty.Empty{})
				if err != nil {
					log.Fatalf("could not get log config: %v", err)
				}

				if util.Marshal(file, format, r) != nil {
					os.Exit(1)
				}
			case "activerole":
				r, err := client.GetRolesForActiveUser(context.Background(), &empty.Empty{})
				if err != nil {
					log.Fatalf("could not get active role: %v", err)
				}

				if util.Marshal(file, format, r) != nil {
					os.Exit(1)
				}
			case "role":
				request := api.RoleMappingsListRequest{}
				request.Page = page
				request.PageSize = pageSize

				r, err := client.ListRoleMappings(context.Background(), &request)
				if err != nil {
					log.Fatalf("could not get active role: %v", err)
				}

				if util.Marshal(file, format, r) != nil {
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
	ValidArgs: util.GetObjectNames(),
}

func printValidObjectTypes() string {
	var names string
	for _, v := range util.GetObjectNames() {
		names += fmt.Sprintf("  * %s\n", v)
	}
	return fmt.Sprintf("Valid object types include:\n%s\n", names)
}

func init() {
	RootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	getCmd.PersistentFlags().StringVarP(&label, "label", "l", "", "List objects by label (<type>:<value> | <value>)")
	getCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "List objects by name (accepts regular expressions)")
	getCmd.PersistentFlags().StringVarP(&format, "format", "f", "table", "Output format (table|json|yaml)")
	getCmd.PersistentFlags().StringVarP(&file, "output", "o", "", "File name to write to")
	getCmd.PersistentFlags().Int32VarP(&pageSize, "pagesize", "s", 10, "Number of objects to get")
	getCmd.PersistentFlags().Int32VarP(&page, "page", "p", 0, "The page number")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
