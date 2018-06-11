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
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/nlnwa/veidemannctl/src/cmd/reports"
	"github.com/nlnwa/veidemannctl/src/configutil"
	"github.com/nlnwa/veidemannctl/src/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)

var (
	cfgFile            string
	controllerAddress  string
	rootCAs            string
	serverNameOverride string
	debug              bool
)

const (
	bash_completion_func = `__veidemannctl_parse_get() {
    local template
    template="{{println .id}}"
    local veidemannctl_out
    if mapfile -t veidemannctl_out < <( veidemannctl get "$1" -o template -t "${template}" 2>/dev/null ); then
        mapfile -t COMPREPLY < <( compgen -W "$( printf '%q ' "${veidemannctl_out[@]}" )" -- "$cur" | awk '/ / { print "\""$0"\"" } /^[^ ]+$/ { print $0 }' )
    fi
}

__veidemannctl_get_resource() {
    if [[ ${#nouns[@]} -eq 0 ]]; then
        return 1
    fi
    __veidemannctl_parse_get ${nouns[${#nouns[@]} -1]}
    if [[ $? -eq 0 ]]; then
        return 0
    fi
}

__custom_func() {
    case ${last_command} in
        veidemannctl_get)
            __veidemannctl_get_resource
            return
            ;;
        *)
            ;;
    esac
}

__veidemannctl_get_name() {
    if [[ ${#nouns[@]} -eq 0 ]]; then
        return 1
    fi
	local noun
	noun=${nouns[${#nouns[@]} -1]}
    local template
    template="{{println .meta.name}}"
    local veidemannctl_out
    if mapfile -t veidemannctl_out < <( veidemannctl get "$noun" -n "^${cur}" -s20 -o template -t "${template}" 2>/dev/null ); then
	    mapfile -t COMPREPLY < <( printf '%q\n' "${veidemannctl_out[@]}" | awk -v IGNORECASE=1 -v p="$cur" '{p==substr($0,0,length(p))} / / {print "\""$0"\"" } /^[^ ]+$/ { print $0 }' )
    fi
}`
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "veidemannctl",
	Short: "Veidemann command line client",
	Long:  `A command line client for Veidemann which can manipulate configs and request status of the crawler.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	BashCompletionFunction: bash_completion_func,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	RootCmd.Version = version.Version.GetVersionString()

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.veidemannctl.yaml)")

	RootCmd.PersistentFlags().StringVarP(&controllerAddress, "controllerAddress", "c", "localhost:50051", "Address to the Controller service")
	viper.BindPFlag("controllerAddress", RootCmd.PersistentFlags().Lookup("controllerAddress"))

	RootCmd.PersistentFlags().StringVar(&rootCAs, "trusted-ca", "", "File with trusted certificate chains for the idp and controller."+
		" These are in addition to the default certs configured for the OS.")

	RootCmd.PersistentFlags().StringVar(&serverNameOverride, "serverNameOverride", "",
		"If set, it will override the virtual host name of authority (e.g. :authority header field) in requests.")
	viper.BindPFlag("serverNameOverride", RootCmd.PersistentFlags().Lookup("serverNameOverride"))

	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Turn on debugging")

	RootCmd.SetVersionTemplate("{{.Version}}")

	RootCmd.AddCommand(reports.ReportCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	var home string
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		var err error
		home, err = homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		// Search config in home directory with name ".veidemannctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".veidemannctl")
	}

	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	} else {
		viper.SetConfigFile(home + "/.veidemannctl.yaml")
	}

	if rootCAs != "" {
		rootCABytes, err := ioutil.ReadFile(rootCAs)
		if err != nil {
			log.Fatalf("failed to read root-ca: %v", err)
		}
		viper.Set("rootCAs", string(rootCABytes))
	}

	if RootCmd.PersistentFlags().Changed("controllerAddress") ||
		RootCmd.PersistentFlags().Changed("idp") ||
		RootCmd.PersistentFlags().Changed("trusted-ca") ||
		RootCmd.PersistentFlags().Changed("serverNameOverride") {

		configutil.WriteConfig()
	}
}
