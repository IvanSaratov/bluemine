package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/gorilla/sessions"

	"github.com/IvanSaratov/bluemine/config"
	"github.com/bradfitz/gomemcache/memcache"
)

var Core struct {
	DB    *sql.DB
	Store *sessions.CookieStore
}

type (
	SessionInfo struct {
		Id   uint64
		Name string
	}
)

var mc *memcache.Client

func InitSession() {
	rand.Seed(time.Now().UnixNano())
	mc = memcache.New(config.Conf.Memcache)
}

func CreateSession(info *SessionInfo) (id string, err error) {
	var datas []byte
	datas, err = json.Marshal(info)
	if err != nil {
		return
	}

	id = fmt.Sprint(rand.Int63())
	err = mc.Set(&memcache.Item{Key: "session_" + id, Value: datas})
	if err != nil {
		id = ""
		return
	}

	return
}

func GetSessionInfo(id string) (result *SessionInfo, err error) {
	var item *memcache.Item
	item, err = mc.Get("session_" + id)
	if err != nil {
		return
	}

	datas := item.Value

	result = new(SessionInfo)
	err = json.Unmarshal(datas, &result)

	return
}
