package db

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const reloadTime int = 300 // 5分钟切换一次buffer
// const reloadTime int = 30 // test半分钟切换一次buffer

type UserRegisterInfoSqlData struct {
	Public_key string   `db:"public_key"`
	Images     [][]byte `db:"images"`
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
		log.Fatalf("connect database error %s", err)
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
	var userInfoItems = c.userRegisterInfoBuffers[c.getUserRegisterInfosBackIndex()]
	for _, userInfoItem := range userInfoItems {
		imagesMap := make(map[string][]byte)
		for i, imgData := range userInfoItem.Images {
			imagesMap[fmt.Sprintf("image%d", i+1)] = imgData
		}

		// Convert map to JSON
		imagesJSON, err := json.Marshal(imagesMap)
		if err != nil {
			log.Fatal(err)
		}

		// // Encrypt JSON data
		// key := []byte("your-secret-key") // Replace with your secret key
		// encryptedJSON, err := EncryptJSON(key, imagesJSON)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		//imageBase64 := base64.StdEncoding.EncodeToString(imagesJSON)

		_, err = c.mysqlClient.Exec(
			"INSERT INTO register_info.register_images (public_key, images) VALUES (?, ?)",
			userInfoItem.Public_key, imagesJSON,
		)
		if err != nil {
			log.Print(err)
			return err
		}
	}
	c.userRegisterInfoBuffers[c.getUserRegisterInfosBackIndex()] = []UserRegisterInfoSqlData{}
	c.switchUserRegisterInfosIndex()
	return nil
}

// 提供对外注册接口
func (c *UserRegisterInfoMysqlMgr) UserDataWrite(public_key string, images [][]byte) {
	userInfoItem := UserRegisterInfoSqlData{
		Public_key: public_key,
		Images:     images,
	}
	c.userRegisterInfoBuffers[c.getUserRegisterInfosIndex()] = append(c.userRegisterInfoBuffers[c.getUserRegisterInfosIndex()], userInfoItem)
}
