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

package configutil

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	yaml2 "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type config struct {
	ControllerAddress  string `json:"controllerAddress"`
	AccessToken        string `json:"accessToken"`
	Nonce              string `json:"nonce"`
	RootCAs            string `json:"rootCAs"`
	ServerNameOverride string `json:"serverNameOverride"`
	ApiKey             string `json:"apiKey"`
}

var GlobalFlags struct {
	Context            string
	ControllerAddress  string
	ServerNameOverride string
	ApiKey             string
	IsShellCompletion  bool
}

func WriteConfig() {
	log.Debug("Writing config")

	c := config{
		viper.GetString("controllerAddress"),
		viper.GetString("accessToken"),
		viper.GetString("nonce"),
		viper.GetString("rootCAs"),
		viper.GetString("serverNameOverride"),
		viper.GetString("apiKey"),
	}

	y, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}

	f, err := os.Create(viper.ConfigFileUsed())
	if err != nil {
		log.Fatalf("Could not create config file '%s': %v", viper.ConfigFileUsed(), err)
	}
	if err = f.Chmod(0600); err != nil {
		log.Fatalf("Could not set access mode on config file '%s': %v", viper.ConfigFileUsed(), err)
	}
	defer f.Close()

	if _, err = f.Write(y); err != nil {
		log.Fatalf("Could write to config file '%s': %v", viper.ConfigFileUsed(), err)
	}
}

func GetConfigDir(subdir string) string {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(home, ".veidemann", subdir)
}

type context struct {
	Context string
}

func GetCurrentContext() (string, error) {
	switch GlobalFlags.Context {
	case "":
		if GlobalFlags.Context == "kubectl" {
		}
		contextDir := GetConfigDir("contexts")
		log.Debugf("Creating context directory: %s", contextDir)
		if err := os.MkdirAll(contextDir, 0777); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		contextFile := GetConfigDir("context.yaml")
		f, err := os.Open(contextFile)
		if err != nil {
			err = SetCurrentContext("default")
			if err != nil {
				return "", err
			}

			f, err = os.Open(contextFile)
			if err != nil {
				log.Fatalf("Could not read file '%s': %v", contextFile, err)
				return "", err
			}
		}
		defer f.Close()

		var c context
		dec := yaml2.NewDecoder(f)
		err = dec.Decode(&c)
		if err != nil {
			log.Fatalf("Could not read file '%s': %v", contextFile, err)
			return "", err
		}

		return c.Context, err
	case "kubectl":
		output, err := exec.Command("kubectl", "config", "current-context").CombinedOutput()
		if err != nil {
			_, _ = os.Stderr.WriteString(err.Error())
			return "", err
		}
		return strings.TrimSpace(string(output)), nil
	default:
		return GlobalFlags.Context, nil
	}
}

func SetCurrentContext(ctx string) error {
	contextFile := GetConfigDir("context.yaml")
	w, err := os.Create(contextFile)
	if err != nil {
		log.Fatalf("Could not open or create file '%s': %v", contextFile, err)
		return err
	}
	enc := yaml2.NewEncoder(w)
	err = enc.Encode(context{ctx})
	if err != nil {
		log.Fatalf("Could not write file '%s': %v", contextFile, err)
		return err
	}
	_ = enc.Close()
	_ = w.Close()

	return nil
}

func ListContexts() ([]string, error) {
	var files []string
	contextDir := GetConfigDir("contexts")
	fileInfo, err := ioutil.ReadDir(contextDir)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if !file.IsDir() {
			sufIdx := strings.LastIndex(file.Name(), ".")
			if sufIdx > 0 {
				files = append(files, file.Name()[:sufIdx])
			}
		}
	}
	return files, nil
}

func ContextExists(ctx string) (bool, error) {
	cs, err := ListContexts()
	if err != nil {
		return false, err
	}

	for _, c := range cs {
		if ctx == c {
			return true, nil
		}
	}
	return false, nil
}
