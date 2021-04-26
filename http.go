package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func HTTPServer() {
	//http.HandleFunc("/pk/get_map", GetMapHandler)         	  //Не используется
	http.HandleFunc("/pk/choose_slot", ChooseSlotHandler)  //Выбор слота token,map_id,slot_id
	http.HandleFunc("/pk/get_token", GetTokenHandler)      //Получения токена user_id
	http.HandleFunc("/pk/update_map", UpdMapHandler)       //Обновление карты map_id
	http.HandleFunc("/pk/parking_enter", ParEntHandler)    //Въезд на парковку token, map_id
	http.HandleFunc("/pk/parking_leave", ParkLeaveHandler) //Выезд из парковки token
	fmt.Println("Server ready...")
	log.Fatal(http.ListenAndServe(":2229", nil))
}
func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	fmt.Println(r.RemoteAddr, r.RequestURI)

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Blocked", "application/json")
	if len(r.Form["user_id"]) > 0 && r.Form["user_id"][0] != "" {
		UserList.Lock()
		_, ok := UserList.U[r.Form["user_id"][0]]
		if !ok {
			UserList.U[r.Form["user_id"][0]] = User{
				ID:     r.Form["user_id"][0],
				Token:  tokenGenerator(),
				Status: 0,
			}
			User := UserList.U[r.Form["user_id"][0]]
			UserList.Unlock()
			Response := Response{
				Data:   []map[string]interface{}{{"User": User}},
				Result: "ok",
			}
			toSend, err := json.Marshal(Response)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprintf(w, string(toSend))
		} else {
			User := UserList.U[r.Form["user_id"][0]]
			UserList.Unlock()
			Response := Response{
				Data:   []map[string]interface{}{{"User": User}},
				Result: "ok",
			}
			toSend, err := json.Marshal(Response)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprintf(w, string(toSend))
		}

	} else {
		fmt.Fprintf(w, ERRFailAuth)
	}
}

func UpdMapHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	fmt.Println(r.RemoteAddr, r.RequestURI)

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Blocked", "application/json")

	var Response Response
	if len(r.Form["map_id"]) > 0 {
		if !MapExist(r.Form["map_id"][0]) {
			fmt.Fprintf(w, ERRMapNotFound)
			return
		}
		MapList.Lock()
		Response.Data = append(Response.Data, map[string]interface{}{"map": MapList.M[r.Form["map_id"][0]], "free": MapList.M[r.Form["map_id"][0]].FreeSlot(), "all": MapList.M[r.Form["map_id"][0]].AllSlot(), "hvz": MapList.M[r.Form["map_id"][0]].HVZSlot()})
		Response.Result = "ok"
		toSend, err := json.Marshal(Response)
		if err != nil {
			fmt.Println(err)
		}
		MapList.Unlock()
		fmt.Fprintf(w, string(toSend))
	} else {
		fmt.Fprintf(w, ERRFailAuth)
	}

}

func GetMapHandler(w http.ResponseWriter, r *http.Request) {

}

func ParEntHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	fmt.Println(r.RemoteAddr, r.RequestURI)

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Blocked", "application/json")

	var Response Response
	if len(r.Form["token"]) > 0 {
		UserList.Lock()
		var UserID string
		if UserList.U.GetIDByToken(r.Form["token"][0]) != "" {
			UserID = UserList.U.GetIDByToken(r.Form["token"][0])
		} else {
			fmt.Fprintf(w, ERRFailAuth)
			UserList.Unlock()
			return
		}
		UserList.Unlock()

		if len(r.Form["map_id"]) > 0 {
			if !MapExist(r.Form["map_id"][0]) {
				fmt.Fprintf(w, ERRMapNotFound)
				return
			}
			MapID := r.Form["map_id"][0]
			var Guest = true
			MapList.Lock()
			UserList.Lock()
			if MapList.M[MapID].FreeSlot() == 0 {
				MapList.Unlock()
				UserList.Unlock()
				fmt.Fprintf(w, ERRParkFull)
				return
			}
			for _, cur := range UserList.U[UserID].Parks {
				if cur == MapID {
					Guest = false
					break
				}
			}
			if UserList.U[UserID].Status != 0 {
				MapList.Unlock()
				UserList.Unlock()
				fmt.Fprintf(w, ERRUserCanT)
				return
			}
			User := UserList.U[UserID]
			User.Status = 1
			User.Parking = MapID
			UserList.U[UserID] = User // Обновление статуса

			MapList.Unlock()
			UserList.Unlock()
			Response.Data = append(Response.Data, map[string]interface{}{"guest": Guest})
			Response.Result = "ok"
			toSend, err := json.Marshal(Response)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Fprintf(w, string(toSend))
		} else {
			fmt.Fprintf(w, ERRFailAuth)
		}
	} else {
		fmt.Fprintf(w, ERRFailAuth)
	}

}

func ChooseSlotHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	fmt.Println(r.RemoteAddr, r.RequestURI)

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Blocked", "application/json")

	var Response Response

	if len(r.Form["token"]) > 0 {
		UserList.Lock()
		var UserID string
		if UserList.U.GetIDByToken(r.Form["token"][0]) != "" {
			UserID = UserList.U.GetIDByToken(r.Form["token"][0])
		} else {
			fmt.Fprintf(w, ERRFailAuth)
			return
		}
		UserList.Unlock()

		if len(r.Form["map_id"]) > 0 {
			if !MapExist(r.Form["map_id"][0]) {
				fmt.Fprintf(w, ERRMapNotFound)
				return
			}
			MapID := r.Form["map_id"][0]
			var slot string
			if len(r.Form["slot_id"]) < 1 {
				fmt.Fprintf(w, ERRSlotNotFound)
				return
			} else {
				slot = r.Form["slot_id"][0]
			}

			MapList.Lock()
			UserList.Lock()

			if UserList.U[UserID].Status != 1 {
				MapList.Unlock()
				UserList.Unlock()
				fmt.Fprintf(w, ERRUserCanT)
				return
			}
			Map := MapList.M[MapID]
			if Map.GetSlotById(slot).Blocked {
				MapList.Unlock()
				UserList.Unlock()
				fmt.Fprintf(w, ERRSlotBlocked)
				return
			}
			Map.GetSlotById(slot).Block(UserID)
			MapList.M[MapID] = Map

			User := UserList.U[UserID]
			User.Status = 2
			UserList.U[UserID] = User // Обновление статуса

			MapList.Unlock()
			UserList.Unlock()
			Response.Result = "ok"
			toSend, err := json.Marshal(Response)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Fprintf(w, string(toSend))
		} else {
			fmt.Fprintf(w, ERRFailAuth)
		}
	} else {
		fmt.Fprintf(w, ERRFailAuth)
	}

}

func ParkLeaveHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	fmt.Println(r.RemoteAddr, r.RequestURI)

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Blocked", "application/json")

	var Response Response

	if len(r.Form["token"]) > 0 {
		UserList.Lock()
		var UserID string
		if UserList.U.GetIDByToken(r.Form["token"][0]) != "" {
			UserID = UserList.U.GetIDByToken(r.Form["token"][0])
		} else {
			fmt.Fprintf(w, ERRFailAuth)
			return
		}
		MapList.Lock()
		MapID := UserList.U[UserID].Parking
		var slot string

		if UserList.U[UserID].Status != 2 {
			MapList.Unlock()
			UserList.Unlock()
			fmt.Fprintf(w, ERRUserCanT)
			return
		}
		Map := MapList.M[MapID]
		fmt.Println(Map)
		for _, cur := range Map.Map {
			if cur.userID == UserID {
				slot = cur.ID
				break
			}
		}
		if slot == "" {
			User := UserList.U[UserID]
			User.Status = 0
			UserList.U[UserID] = User // Обновление статуса

			MapList.Unlock()
			UserList.Unlock()
			fmt.Fprintf(w, SendOK)
			return
		}
		fmt.Println(Map)
		Map.GetSlotById(slot).unBlock()
		fmt.Println(Map)
		MapList.M[MapID] = Map

		User := UserList.U[UserID]
		User.Status = 0
		UserList.U[UserID] = User // Обновление статуса

		MapList.Unlock()
		UserList.Unlock()
		Response.Result = "ok"
		toSend, err := json.Marshal(Response)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Fprintf(w, string(toSend))
	} else {
		fmt.Fprintf(w, ERRFailAuth)
	}

}

func JSFix(w http.ResponseWriter, r *http.Request) {
	if len(r.Form) < 1 {
		fmt.Fprintf(w, "<!DOCTYPE html><script type=\"text/javascript\">\n\nif (location.href!=location.href.replace(\"#\",\"?\"))location.replace(location.href.replace(\"#\",\"?\"));\n\n</script>")
	}
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
}
