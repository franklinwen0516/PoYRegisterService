package db

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

const reloadTime int = 300 // 5分钟切换一次buffer

type UserRegisterInfoSqlData struct {
	Public_key string `db:"public_key"`
	Images     []byte `db:"images"`
}

type UserRegisterInfoMysqlMgr struct {
	mysqlClient *sqlx.DB
	notifySync  chan int
	stopMdbSync bool

	userRegisterInfoBuffers     [][]UserRegisterInfoSqlData
	userRegisterInfoBufferIndex int32
}

var (
	UserRegisterInfoInstance UserRegisterInfoMysqlMgr
)

func (c *UserRegisterInfoMysqlMgr) Init() error {
	client, err := sqlx.Connect("mysql", "root:1234@tcp(127.0.0.1:3306)/register_info")
	if err != nil {
		log.Fatalf("connect database error")
		return err
	}
	c.mysqlClient = client
	c.initDataCache()
	c.notifySync = make(chan int, 1)
	c.stopMdbSync = false
	c.runDataWriteThread()
	return nil
}

func (c *UserRegisterInfoMysqlMgr) getUserRegisterInfosIndex() int32 {
	return c.userRegisterInfoBufferIndex
}

func (c *UserRegisterInfoMysqlMgr) getUserRegisterInfosBackIndex() int32 {
	return 1 - c.userRegisterInfoBufferIndex
}

func (c *UserRegisterInfoMysqlMgr) switchUserRegisterInfosIndex() {
	c.userRegisterInfoBufferIndex = 1 - c.userRegisterInfoBufferIndex
}

func (c *UserRegisterInfoMysqlMgr) initDataCache() {
	c.userRegisterInfoBufferIndex = 0
	c.userRegisterInfoBuffers = make([][]UserRegisterInfoSqlData, 2)
	c.userRegisterInfoBuffers[c.getUserRegisterInfosIndex()] = []UserRegisterInfoSqlData{}
	c.userRegisterInfoBuffers[c.getUserRegisterInfosBackIndex()] = []UserRegisterInfoSqlData{}
}

func (c *UserRegisterInfoMysqlMgr) runDataWriteThread() {
	go func() {
		syncTimer := time.NewTicker(time.Duration(reloadTime) * time.Second)
		defer syncTimer.Stop()
		for !c.stopMdbSync {
			c.writeUsersInfo()
			select {
			case <-syncTimer.C:
			case <-c.notifySync:
			}
		}
	}()
}

func (c *UserRegisterInfoMysqlMgr) writeUsersInfo() error {
	var userInfoItems []UserRegisterInfoSqlData
	err := c.mysqlClient.Select(&userInfoItems,
		"SELECT public_key, images FROM register_images")
	if err != nil {
		log.Fatal(err)
		return err
	}
	// tmpUserBalanceMap := c.processUserInfoSqlData(userInfoItems)
	// if len(tmpUserBalanceMap) > 0 {
	// 	c.userLoginInfos[c.getUserLoginInfosBackIndex()] = tmpUserBalanceMap
	// 	c.switchUserLoginInfosIndex()
	// }
	// debug
	// for k, v := range tmpUserBalanceMap {
	// 	log.Printf("key:%v|value:%v\n", k, v)
	// }
	return nil
}

// 提供对外注册接口
func (c *UserRegisterInfoMysqlMgr) userDataWrite() {

}
