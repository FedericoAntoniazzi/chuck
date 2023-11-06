/*
Copyright © 2023 Federico Antoniazzi

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
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// The prefix for all environment variables
	envPrefix = "CHUCK"
)

// NewRootCommand creates a new chuck command instance
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "chuck",
		Short: "chuck - Container Image Update Checker",
		Long:  `chuck fetches the images from running containers and shows eventual updates`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Bind cobra and viper
			return initializeConfig(cmd)
		},
	}

	rootCmd.AddCommand(NewRunCommand())

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd := NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// initializeConfig binds cobra flags and viper config
func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	// Load configuration from env vars
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	bindFlags(cmd, v)

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name

		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
