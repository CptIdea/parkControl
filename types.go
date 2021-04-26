package main

import "sync"

//0-Пустота
//1-Свободное место
//2-Занятое место
type Map struct {
	Map     []Slot `json:"map"`
	Address string `json:"address"`
	mu      sync.Mutex
}

type Slot struct {
	ID      string `json:"id"`
	Blocked bool   `json:"blocked"`
	HVZ     bool   `json:"hvz"`
	userID  string
}

type Response struct {
	Data   []map[string]interface{} `json:"data"`
	Result string                   `json:"result"`
}

type User struct {
	ID      string   `json:"id"`
	Token   string   `json:"token"`
	Status  int      `json:"status"`
	Parking string   `json:"parking"`
	Parks   []string `json:"parks"`
}
type Users map[string]User
