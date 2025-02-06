/*
Copyright © 2025 Federico Antoniazzi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/FedericoAntoniazzi/chuck/pkg/outputs"
	"github.com/FedericoAntoniazzi/chuck/pkg/registry"
	"github.com/FedericoAntoniazzi/chuck/pkg/runtime"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Check image updates in running containers",
	Run: func(cmd *cobra.Command, args []string) {
		rt, err := runtime.NewDockerRuntime()
		if err != nil {
			slog.Error("Error creating docker client", "err", err)
		}

		ctx := context.Background()
		containers, err := rt.ListRunningContainers(ctx)
		if err != nil {
			slog.Error("Error listing containers", "err", err)
		}

		console := outputs.ConsoleOutput{}

		for _, cnt := range containers {
			tags, err := registry.ListNewerTags(cnt.Image)
			if err != nil {
				slog.Error("Error listing tags", "err", err, "container", cnt.Name, "image", cnt.Image)
			}

			if len(tags) > 1 {
				newestTag := tags[len(tags)-1]
				message := fmt.Sprintf("Container %s can be upgraded to version %s", cnt.Name, newestTag)
				console.Send(message)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
