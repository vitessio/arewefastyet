/*
Copyright 2021 The Vitess Authors.

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
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/vitessio/arewefastyet/go/cmd/admin"
	"github.com/vitessio/arewefastyet/go/cmd/api"
	"github.com/vitessio/arewefastyet/go/cmd/exec"
	"github.com/vitessio/arewefastyet/go/cmd/gen"
	"github.com/vitessio/arewefastyet/go/cmd/macrobench"
	"github.com/vitessio/arewefastyet/go/cmd/microbench"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	secretsFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "arewefastyet",
	Short: "Nightly Benchmarks Project",
	Long:  `Vitess has to ensure it's delivering flawless performance to its users. In order to meet this need, we created AreWeFastYet, a benchmarking monitoring tool for Vitess.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/arewefastyet/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&secretsFile, "secrets", "", "secrets file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.AddCommand(admin.AdminCmd())
	rootCmd.AddCommand(api.ApiCmd())
	rootCmd.AddCommand(microbench.MicroBenchCmd())
	rootCmd.AddCommand(macrobench.MacroBenchCmd())
	rootCmd.AddCommand(exec.ExecCmd())
	rootCmd.AddCommand(gen.GenCmd())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Get home directory.
		home, err := homedir.Expand("~/.config/arewefastyet")
		cobra.CheckErr(err)

		// Search config in home directory with name "~/.config/arewefastyet/config.yaml".
		viper.AddConfigPath(home)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("arewefastyet")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Println(err)
			os.Exit(1)
		}
	}
	if secretsFile != "" {
		viper.SetConfigFile(secretsFile)
		err := viper.MergeInConfig()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	postInitCommands(rootCmd.Commands())
}

// https://github.com/spf13/viper/issues/397
func postInitCommands(commands []*cobra.Command) {
	for _, cmd := range commands {
		presetRequiredFlags(cmd)
		if cmd.HasSubCommands() {
			postInitCommands(cmd.Commands())
		}
	}
}

func presetRequiredFlags(cmd *cobra.Command) {
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		log.Fatal(err)
	}
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			err = cmd.Flags().Set(f.Name, viper.GetString(f.Name))
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}
