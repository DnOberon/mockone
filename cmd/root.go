/*
Copyright © 2020 John Darrington john@thirdhousesoftware.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/spf13/cobra"
	"log"
	"net/http/httputil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mockone [file to serve]",
	Short: "Single endpoint server for multiple file and response types",
	Long: `Mockone allows users to create a single endpoint http server
which serves a variety of files and types.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// process the file
		dir, err := filepath.Abs(args[0])
		if err != nil {
			log.Fatal("unable to get directory of file")
		}

		fmt.Println(dir)

		extension := filepath.Ext(dir)
		if extension == "" {
			log.Fatalf("unable to verify extension of file, valid extensions are .json, .html, and .*")
		}

		// create and run the server
		e := echo.New()

		e.GET("/", func(c echo.Context) error {
			dump, err := httputil.DumpRequest(c.Request(), true)
			if err != nil {
				log.Fatalf("cannot dump request body %s", err.Error())
			}

			log.Println(string(dump))

			return c.File(dir)
		})

		e.Logger.Fatal(e.Start(fmt.Sprintf(":%s",cmd.Flag("port").Value)))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().Int("port", 8091, "port to serve application on - default 8091")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".mockone" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mockone")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
