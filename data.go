package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func GetUsers() Users {
	JSON, err := ioutil.ReadFile("Users.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var result Users
	err = json.Unmarshal(JSON, &result)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return result
}

func GetMap() map[string]Map {
	JSON, err := ioutil.ReadFile("Map.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var result map[string]Map
	err = json.Unmarshal(JSON, &result)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return result
}

func MapExist(ID string) bool {
	MapList.Lock()
	_, ans := MapList.M[ID]
	MapList.Unlock()
	return ans
}

func (m Map) SaveToDB() {

}
