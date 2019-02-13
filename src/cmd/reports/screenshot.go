// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
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

package reports

import (
	"bytes"
	"context"
	api "github.com/nlnwa/veidemann-api-go/veidemann_api"
	"github.com/nlnwa/veidemannctl/src/connection"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"time"
)

var (
	uri            string
	img            bool
	display        bool
	displaySeconds int32
)

// screenshotCmd represents the screenshot command
var screenshotCmd = &cobra.Command{
	Use:   "screenshot",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn := connection.NewReportClient()
		defer conn.Close()

		request := api.ScreenshotListRequest{}
		if len(args) > 0 {
			request.Id = args
		}
		if flags.executionId != "" {
			request.ExecutionId = flags.executionId
		}
		if img || display {
			request.ImgData = true
		}
		request.Filter = applyFilter(flags.filter)
		request.Page = flags.page
		request.PageSize = flags.pageSize

		r, err := client.ListScreenshots(context.Background(), &request)
		if err != nil {
			log.Fatalf("could not get page log: %v", err)
		}

		if img {
			printScreenshot(r.Value[0])
		} else if display {
			showScreenshot(r.Value)
		} else {
			ApplyTemplate(r, "screenshot.template")
		}
	},
}

func init() {
	ReportCmd.AddCommand(screenshotCmd)

	screenshotCmd.Flags().StringVarP(&uri, "uri", "", "", "All screenshots by URI")
	screenshotCmd.Flags().BoolVarP(&img, "img", "", false, "Image binary")
	screenshotCmd.Flags().BoolVarP(&display, "display", "", false, "Show image using ImageMagic's display")
	screenshotCmd.Flags().Int32VarP(&displaySeconds, "displayTime", "", 4, "The time in seconds to show an image. 0 for no timeout")
}

func printScreenshot(screenshot *api.Screenshot) {
	os.Stdout.Write(screenshot.Img)
}

func showScreenshot(screenshot []*api.Screenshot) {
	prog := "feh"
	if _, err := exec.LookPath(prog); err != nil {
		prog = "display"
		if _, err := exec.LookPath(prog); err != nil {
			log.Fatal("Neither 'feh' or 'display' is available")
		}
	}

	for _, im := range screenshot {
		cmd := exec.Command(prog, "-")
		cmd.Stdin = bytes.NewReader(im.Img)

		if displaySeconds == 0 {
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := cmd.Start()
			if err != nil {
				log.Fatal(err)
			}

			time.Sleep(time.Duration(displaySeconds) * time.Second)
			cmd.Process.Kill()
		}
	}
}
