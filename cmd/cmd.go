// @author: xwc1125
// @date: 2020/10/05
package cmd

import (
	"github.com/chain5j/chain5j-pkg/database/leveldb"
	log "github.com/chain5j/log15"
	"github.com/chain5j/sync_eth/params"
	"github.com/chain5j/sync_eth/params/global"
	"github.com/chain5j/sync_eth/pkg/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

func initCli() *cli.Cli {
	rootCli := cli.NewCli(params.App, params.Version, params.Welcome)

	err := rootCli.Init(
		func(viper *viper.Viper, rootFlags *pflag.FlagSet) {
		},
		// --config "./conf/config.yaml"
		func(viper *viper.Viper) {
			err := viper.Unmarshal(&global.Config)
			if err != nil {
				panic(err)
			}
			initLogs()
			initDB()
		})
	if err != nil {
		log.Error("initCli err", "err", err)
	}

	rootCli.RunE(func(cmd *cobra.Command, args []string) error {
		return server()
	})

	return rootCli
}

// init log
func initLogs() {
	logConfig := global.Config.Log
	ostream := log.StreamHandler(os.Stderr, log.TerminalFormat(logConfig.UseColor))
	gLogger := log.NewGlogHandler(ostream)
	gLogger.Verbosity(log.Lvl(logConfig.Level))
	gLogger.VModules(logConfig.GetModules())

	log.PrintOrigins(true)
	log.Root().SetHandler(gLogger)
}

func initDB() {
	// leveldb
	db, err := leveldb.New(global.Config.Database.LevelDb.Path, 0, 0, "")
	if err != nil {
		log.Error("leveldb new err", "err", err)
		panic(err)
	}
	global.LevelDb = db
	log.Info("leveldb init success", "DbType", "levelDB")
}

//Execute : apply commands
func Execute() {
	rootCli := initCli()
	global.RootCli = rootCli
	rootCli.Execute()
}
