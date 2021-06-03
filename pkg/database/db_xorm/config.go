// description: Leetcode
//
// @author: xwc1125
// @date: 2020/10/05
package db_xorm

type MysqlConfig struct {
	DriverName   string `json:"driverName" mapstructure:"driverName"`
	User         string `json:"user" mapstructure:"user"`
	Password     string `json:"password" mapstructure:"password"`
	Host         string `json:"host" mapstructure:"host"`
	Port         int    `json:"port" mapstructure:"port"`
	Database     string `json:"database" mapstructure:"database"`
	Charset      string `json:"charset" mapstructure:"charset"`
	ShowSql      bool   `json:"showSql" mapstructure:"showSql"`
	LogLevel     string `json:"logLevel" mapstructure:"logLevel"`
	MaxIdleConns int    `json:"maxIdleConns" mapstructure:"maxIdleConns"`
	MaxOpenConns int    `json:"maxOpenConns" mapstructure:"maxOpenConns"`
	//ParseTime       bool   `json:"parseTime" mapstructure:"parseTime"`
	ConnMaxLifetime int64 `json:"connMaxLifetime" mapstructure:"connMaxLifetime: 10"` // s
	//Sslmode         string `json:"sslmode" mapstructure:"sslmode"`
}
