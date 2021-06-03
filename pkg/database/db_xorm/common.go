// description: superwallet
//
// @author: xwc1125
// @date: 2020/10/05
package db_xorm

import (
	"fmt"
	"sync"
)

var (
	lock sync.Mutex
)

func GetConnURL(info *MysqlConfig) (url string) {
	url = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		info.User,
		info.Password,
		info.Host,
		info.Port,
		info.Database,
		info.Charset)
	return
}
