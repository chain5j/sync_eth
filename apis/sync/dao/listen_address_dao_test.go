// description: sync_eth 
// 
// @author: xwc1125
// @date: 2020/10/05
package dao

import (
	"fmt"
	"github.com/chain5j/sync_eth/apis/sync/models"
	"github.com/chain5j/sync_eth/pkg/database/es"
	"testing"
)

func TestListenAddressDao_BatchAddListenAddress(t *testing.T) {
	es, _ := es.NewES([]string{"http://127.0.0.1:9200"})
	err := ListenAddressDao.BatchAddListenAddress(es, []*models.ListenAddress{
		{
			Address: "0xAeff996F0Efb374fCf95Eb6b38fd4aA5E4bbC1b1",
			Remark:  "ganache",
		},
		{
			Address: "0x70ac2DBCee2c9f2cA6ab536b09e36205E93d2B24",
			Remark:  "ganache2",
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}
