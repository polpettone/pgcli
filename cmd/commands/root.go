package commands

import (
	"fmt"
	"github.com/polpettone/pgcli/cmd/config"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pgcli",
	Short: "A brief description of your application",
	Long: "",

	Run: func(cmd *cobra.Command, args []string) {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "pgcli: try 'pgcli --help' for more information")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(InitConfig)
	config.InitLogging()
}

func InitConfig() {
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

		viper.AddConfigPath(home + "/.config/pgcli")
		viper.SetConfigType("yaml")
		viper.SetConfigName("conf")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {

		out := `No Config file found
Provide a config file named .pgcli in $HOME
Format yaml
Content:
	project_id: <project id>
	url: https://gitlab.com/api/v4/projects
	api_token: <api token>`

		fmt.Println(out)
		os.Exit(1)
	}
}

