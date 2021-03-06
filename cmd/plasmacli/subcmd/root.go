package subcmd

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/config"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/flags"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/subcmd/eth"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/subcmd/keys"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/subcmd/query"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/subcmd/tx"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// default home directory
var homeDir = os.ExpandEnv("$HOME/.plasmacli/")

// RootCmd returns the initialized root cmd for plasmacli
func RootCmd() *cobra.Command {
	cobra.EnableCommandSorting = false
	rootCmd.PersistentFlags().String(flags.Home, homeDir, "home directory for plasmacli")
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootCmd.AddCommand(
		tx.RootCmd(),
		eth.RootCmd(),
		query.RootCmd(),
		client.LineBreak,

		RestServerCmd(),
		client.LineBreak,

		keys.RootCmd(),
		client.LineBreak,

		VersionCmd(),
	)

	return rootCmd
}

var rootCmd = &cobra.Command{
	Use:           "plasmacli",
	Short:         "Plasma Client",
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		config.RegisterViperAndEnv()

		homeDir := viper.GetString(flags.Home)
		store.InitKeystore(homeDir)

		configFilepath := filepath.Join(homeDir, "config.toml")
		if _, err := os.Stat(configFilepath); os.IsNotExist(err) {
			if err := config.WriteConfigFile(configFilepath, config.DefaultConfig()); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		viper.AddConfigPath(homeDir)
		viper.SetConfigName("config")
		if err := viper.MergeInConfig(); err != nil {
			return err
		}

		return nil
	},
}
