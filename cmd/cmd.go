// Copyright © 2017 National Library of Norway.
// Licensed under the Apache License, GitVersion 2.0 (the "License");
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
	"os"

	"github.com/nlnwa/veidemannctl/cmd/abort"
	"github.com/nlnwa/veidemannctl/cmd/abortjobexecution"
	"github.com/nlnwa/veidemannctl/cmd/activeroles"
	configcmd "github.com/nlnwa/veidemannctl/cmd/config"
	"github.com/nlnwa/veidemannctl/cmd/create"
	deletecmd "github.com/nlnwa/veidemannctl/cmd/delete"
	"github.com/nlnwa/veidemannctl/cmd/get"
	importcmd "github.com/nlnwa/veidemannctl/cmd/import"
	"github.com/nlnwa/veidemannctl/cmd/logconfig"
	"github.com/nlnwa/veidemannctl/cmd/login"
	"github.com/nlnwa/veidemannctl/cmd/logout"
	"github.com/nlnwa/veidemannctl/cmd/pause"
	"github.com/nlnwa/veidemannctl/cmd/report"
	"github.com/nlnwa/veidemannctl/cmd/run"
	"github.com/nlnwa/veidemannctl/cmd/script_parameters"
	"github.com/nlnwa/veidemannctl/cmd/status"
	"github.com/nlnwa/veidemannctl/cmd/unpause"
	"github.com/nlnwa/veidemannctl/cmd/update"
	"github.com/nlnwa/veidemannctl/cmd/version"
	"github.com/nlnwa/veidemannctl/config"

	"github.com/spf13/cobra"
)

// NewRootCmd returns the root command.
func NewRootCmd() *cobra.Command {
	cobra.EnableCommandSorting = false

	cmd := &cobra.Command{
		Use:               "veidemannctl",
		Short:             "veidemannctl controls the Veidemann web crawler",
		Long:              "veidemannctl controls the Veidemann web crawler",
		DisableAutoGenTag: true,
	}

	// Add global flags
	cmd.PersistentFlags().String("config", "", "Path to the config file to use (By default configuration file is stored under $HOME/.veidemann/contexts/")
	cmd.PersistentFlags().String("context", "", "The name of the context to use")
	cmd.PersistentFlags().String("server", "", "The address of the Veidemann server to use")
	cmd.PersistentFlags().String("server-name-override", "",
		"If set, it will override the virtual host name of authority (e.g. :authority header field) in requests")
	cmd.PersistentFlags().String("api-key", "",
		"If set, it will be used as the bearer token for authentication")
	cmd.PersistentFlags().String("log-level", "info", `set log level, available levels are "panic", "fatal", "error", "warn", "info", "debug" and "trace"`)
	cmd.PersistentFlags().String("log-format", "pretty", `set log format, available formats are: "pretty" or "json"`)
	cmd.PersistentFlags().Bool("log-caller", false, "include information about caller in log output")

	// Add subcommands
	cmd.AddCommand(configcmd.NewConfigCmd()) // config

	cmd.AddGroup(&cobra.Group{
		ID:    "basic",
		Title: "Basic Commands:",
	})
	cmd.AddCommand(get.NewGetCmd())          // get
	cmd.AddCommand(create.NewCreateCmd())    // create
	cmd.AddCommand(update.NewUpdateCmd())    // update
	cmd.AddCommand(deletecmd.NewDeleteCmd()) // delete

	cmd.AddGroup(&cobra.Group{
		ID:    "advanced",
		Title: "Advanced Commands:",
	})
	cmd.AddCommand(report.NewReportCmd())    // report
	cmd.AddCommand(importcmd.NewImportCmd()) // import

	cmd.AddGroup(&cobra.Group{
		ID:    "run",
		Title: "Crawl Commands:",
	})
	cmd.AddCommand(run.NewRunCmd())                             // run
	cmd.AddCommand(abort.NewAbortCmd())                         // abort
	cmd.AddCommand(abortjobexecution.NewAbortJobExecutionCmd()) // abortjobexecution

	cmd.AddGroup(&cobra.Group{
		ID:    "status",
		Title: "Management Commands:",
	})
	cmd.AddCommand(status.NewStatusCmd())   // status
	cmd.AddCommand(pause.NewPauseCmd())     // pause
	cmd.AddCommand(unpause.NewUnpauseCmd()) // unpause

	cmd.AddGroup(&cobra.Group{
		ID:    "login",
		Title: "Authentication Commands:",
	})
	cmd.AddCommand(login.NewLoginCmd())   // login
	cmd.AddCommand(logout.NewLogoutCmd()) // logout

	cmd.AddCommand(version.NewVersionCmd()) // version

	cmd.AddGroup(&cobra.Group{
		ID:    "debug",
		Title: "Troubleshooting and Debug Commands:",
	})
	cmd.AddCommand(scriptparameters.NewScriptParametersCmd()) // script-parameters
	cmd.AddCommand(logconfig.NewLogConfigCmd())               // logconfig
	cmd.AddCommand(activeroles.NewActiveRolesCmd())           // activeroles

	return cmd
}

// Execute initializes the root command and executes it.
func Execute() {
	// Initialize root command
	cmd := NewRootCmd()

	// Register function to run after command is initialized
	cobra.OnInitialize(func() {
		// Initialize config from flags
		err := config.Init(cmd.PersistentFlags())
		if err != nil {
			fmt.Printf("Initialization failed: %v\n", err)
			os.Exit(1)
		}
	})

	// Execute root command
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
