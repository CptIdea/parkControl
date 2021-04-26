package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"
)

var MapList struct {
	M map[string]Map
	sync.Mutex
}

var UserList = struct {
	U Users
	sync.Mutex
}{}

func main() {
	MapList.M = GetMap()
	fmt.Println(MapList.M)
	UserList.U = GetUsers()
	fmt.Println(UserList.U)

	go HTTPServer()
	for {
		MapList.Lock()
		TXT, _ := json.Marshal(MapList.M)
		MapList.Unlock()
		ioutil.WriteFile("Map.json", []byte(TXT), 700)

		UserList.Lock()
		TXT, _ = json.Marshal(UserList.U)
		UserList.Unlock()
		ioutil.WriteFile("Users.json", []byte(TXT), 700)
		time.Sleep(30 * time.Second)
	}
}

func tokenGenerator() string {
	b := make([]byte, 13)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
