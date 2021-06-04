// description: sync_eth
//
// @author: xwc1125
// @date: 2020/10/05
package cmd

import (
	"errors"
	log "github.com/chain5j/log15"
	"github.com/chain5j/sync_eth/apis/sync/chain"
	"github.com/chain5j/sync_eth/apis/sync/engine"
	"github.com/chain5j/sync_eth/params/global"
	"os"
	"os/signal"
)

func server() error {
	initRpc()

	sync, err := engine.NewSync(global.Config.Database.Es.Host)
	if err != nil {
		log.Error("engine.NewSync err", "err", err)
		return err
	}
	err = sync.Start()
	if err != nil {
		log.Error("sync.Start err", "err", err)
		return err
	}

	log.Info("Enter Control + C Shutdown Server")

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutdown Server ...")

	log.Info("Server exiting")
	return nil
}

func initRpc() error {
	host := global.Config.ChainConfig.Host
	if host == "" {
		return errors.New("rpc host is empty")
	}
	clientIdentifier := global.Config.ChainConfig.ClientIdentifier
	if clientIdentifier == "" {
		clientIdentifier = "eth"
	}

	eth, err := chain.NewETH(host, clientIdentifier, global.Config.ChainConfig.ChainId, global.Config.ChainConfig.IsEip155)
	if err != nil {
		return err
	}
	global.RpcClient = eth
	return nil
}
