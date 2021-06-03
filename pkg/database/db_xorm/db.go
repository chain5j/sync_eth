package db_xorm

import (
	log "github.com/chain5j/log15"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"xorm.io/core"
	"xorm.io/xorm"
)

var (
	masterEngine *xorm.Engine
	slaveEngine  *xorm.Engine
)

func MasterEngine(master MysqlConfig) *xorm.Engine {
	if masterEngine != nil {
		return masterEngine
	}

	lock.Lock()
	defer lock.Unlock()

	if masterEngine != nil {
		return masterEngine
	}

	engine, err := xorm.NewEngine(master.DriverName, GetConnURL(&master))
	if err != nil {
		log.Error("Instance Master DB error!!", "err", err)
		return nil
	}
	settings(engine, &master)
	engine.SetMapper(core.GonicMapper{})
	masterEngine = engine
	return masterEngine
}

var SysTimeLocation, _ = time.LoadLocation("Asia/Shanghai")

func settings(engine *xorm.Engine, info *MysqlConfig) {
	engine.ShowSQL(info.ShowSql)
	engine.SetTZLocation(SysTimeLocation)
	if info.MaxIdleConns > 0 {
		engine.SetMaxIdleConns(info.MaxIdleConns)
	}
	if info.MaxOpenConns > 0 {
		engine.SetMaxOpenConns(info.MaxOpenConns)
	}
	if info.ConnMaxLifetime > 0 {
		engine.SetConnMaxLifetime(time.Duration(info.ConnMaxLifetime) * time.Second)
	}

	// cache
	//cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	//engine.SetDefaultCacher(cacher)
}
