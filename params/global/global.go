//
// @author: xwc1125
// @date: 2020/10/05
package global

import (
	"github.com/chain5j/chain5j-pkg/database/leveldb"
	"github.com/chain5j/sync_eth/apis/sync/chain"
	"github.com/chain5j/sync_eth/params"
	"github.com/chain5j/sync_eth/pkg/cli"
	"github.com/chain5j/sync_eth/pkg/database/db_xorm"
	"github.com/chain5j/sync_eth/pkg/database/mq"
	"strings"
)

var (
	RootCli          *cli.Cli
	Config           *YamlConfig
	LevelDb          *leveldb.Database
	RpcClient        *chain.Eth
	RabbitMqTxRaw    *mq.Wrapper
	RabbitMqTxNotify *mq.Wrapper
)

type YamlConfig struct {
	Server      *ServerConfig   `json:"server" mapstructure:"server"`
	App         *AppConfig      `json:"app" mapstructure:"app"`
	Log         *LogConfig      `json:"log" mapstructure:"log"`
	Database    *DatabaseConfig `json:"database" mapstructure:"database"`
	ChainConfig *EthChainConfig `json:"chainConfig" mapstructure:"chainConfig"`
}

func (c *YamlConfig) GetApp() *AppConfig {
	if c.App == nil {
		return &AppConfig{
			Name:        params.App,
			Version:     params.Version,
			Description: params.Welcome,
		}
	}
	return c.App
}

type ServerConfig struct {
	BaseUrl string `json:"baseUrl" mapstructure:"baseUrl"`
	Addr    string `json:"addr" mapstructure:"addr"`
	Port    int    `json:"port" mapstructure:"port"`
}

type AppConfig struct {
	Name        string `json:"name" mapstructure:"name"`
	Version     string `json:"version" mapstructure:"version"`
	Description string `json:"description" mapstructure:"description"`
}

type DatabaseConfig struct {
	Master  db_xorm.MysqlConfig `json:"master" mapstructure:"master"`
	Es      *EsConfig           `json:"es" mapstructure:"es"`
	LevelDb *LevelDbConfig      `json:"levelDb" mapstructure:"levelDb"`
	Mq      *mq.Config          `json:"mq" mapstructure:"mq"`
}

type LogConfig struct {
	Modules  string `json:"modules" mapstructure:"modules"`
	FilePath string `json:"file_path" mapstructure:"file_path"`
	FileName string `json:"file_name" mapstructure:"file_name"`
	Format   string `json:"format" mapstructure:"format"`
	Console  bool   `json:"console" mapstructure:"console"`
	Level    int    `json:"level" mapstructure:"level"`
	UseColor bool   `json:"use_color" mapstructure:"use_color"`
}

//LevelDb config
type LevelDbConfig struct {
	Path string `json:"path" mapstructure:"path"`
}

func (c *LogConfig) GetModules() []string {
	if c.Modules == "" {
		return []string{"*"}
	}
	return strings.Split(c.Modules, ",")
}
