// description: sync_eth 
// 
// @author: xwc1125
// @date: 2020/10/05
package mq

import "fmt"

type Config struct {
	Host     string `json:"host" mapstructure:"host"`
	Port     int64  `json:"port" mapstructure:"port"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
	IsUse    bool   `json:"isUse" mapstructure:"isUse"`
}

func (c *Config) GetConnUtl() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", c.Username, c.Password, c.Host, c.Port)
}
