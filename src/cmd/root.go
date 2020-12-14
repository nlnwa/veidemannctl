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
	"github.com/nlnwa/veidemannctl/bindata"
	"github.com/nlnwa/veidemannctl/src/cmd/config"
	"github.com/nlnwa/veidemannctl/src/cmd/importcmd"
	"github.com/nlnwa/veidemannctl/src/cmd/logconfig"
	"github.com/nlnwa/veidemannctl/src/cmd/reports"
	"github.com/nlnwa/veidemannctl/src/configutil"
	"github.com/nlnwa/veidemannctl/src/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	cfgFile string
	debug   bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "veidemannctl",
	Short: "Veidemann command line client",
	Long:  `A command line client for Veidemann which can manipulate configs and request status of the crawler.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
	DisableAutoGenTag: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	RootCmd.Version = version.Version.GetVersionString()

	data, err := bindata.Asset("completion.sh")
	if err != nil {
		panic(err)
	}
	RootCmd.BashCompletionFunction = string(data)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.veidemannctl.yaml)")

	RootCmd.PersistentFlags().StringVar(&configutil.GlobalFlags.Context, "context", "", "The name of the veidemannconfig context to use.")

	RootCmd.PersistentFlags().StringVarP(&configutil.GlobalFlags.ControllerAddress, "controllerAddress", "c", "localhost:50051", "Address to the Controller service")

	RootCmd.PersistentFlags().StringVar(&configutil.GlobalFlags.ServerNameOverride, "serverNameOverride", "",
		"If set, it will override the virtual host name of authority (e.g. :authority header field) in requests.")

	RootCmd.PersistentFlags().StringVar(&configutil.GlobalFlags.ApiKey, "apiKey", "",
		"Api-key used for authentication instead of interactive logon trough IDP.")

	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Turn on debugging")

	RootCmd.PersistentFlags().BoolVar(&configutil.GlobalFlags.IsShellCompletion, "comp", false, "Clean output used for shell completion")
	_ = RootCmd.PersistentFlags().MarkHidden("comp")

	RootCmd.SetVersionTemplate("{{.Version}}")

	RootCmd.AddCommand(reports.ReportCmd)
	RootCmd.AddCommand(logconfig.LogconfigCmd)
	RootCmd.AddCommand(config.ConfigCmd)
	RootCmd.AddCommand(importcmd.ImportCmd)
	RootCmd.AddCommand(scriptParametersCmd)
}

// initConfig reads in config file and ENV variables if set.
var contextDir string

func initConfig() {
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		var err error
		configutil.GlobalFlags.Context, err = configutil.GetCurrentContext()
		if err != nil {
			log.Fatalf("Could not get current context: %v", err)
		}
		contextDir = configutil.GetConfigDir("contexts")

		// Search config in home directory with name ".veidemannctl" (without extension).
		viper.AddConfigPath(contextDir)
		viper.SetConfigName(configutil.GlobalFlags.Context)

		log.Debug("Using context: ", configutil.GlobalFlags.Context)
	}

	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	} else {
		log.Debugf("Setting config file to the non-existing file: %s", filepath.Join(contextDir, configutil.GlobalFlags.Context+".yaml"))
		viper.SetConfigFile(filepath.Join(contextDir, configutil.GlobalFlags.Context+".yaml"))
	}

	if !RootCmd.PersistentFlags().Changed("controllerAddress") {
		configutil.GlobalFlags.ControllerAddress = viper.GetString("controllerAddress")
	}

	if !RootCmd.PersistentFlags().Changed("serverNameOverride") {
		configutil.GlobalFlags.ServerNameOverride = viper.GetString("serverNameOverride")
	}

	if !RootCmd.PersistentFlags().Changed("apiKey") {
		configutil.GlobalFlags.ApiKey = viper.GetString("apiKey")
	}
	log.Debug("Using controller address: ", configutil.GlobalFlags.ControllerAddress)
	log.Debug("Using server name override: ", configutil.GlobalFlags.ServerNameOverride)
	log.Debug("Using api-key: ", configutil.GlobalFlags.ApiKey)
}
