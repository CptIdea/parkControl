package main

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

//Свободные слоты
func (m Map) FreeSlot() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	var ans int
	for _, cur := range m.Map {
		if !cur.Blocked {
			ans++
		}
	}
	return ans
}

//Все слоты
func (m Map) AllSlot() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.Map)
}

func (m Map) HVZSlot() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	var ans int
	for _, cur := range m.Map {
		if cur.HVZ {
			ans++
		}
	}
	return ans
}

func GenerateParking(name string, slotCount int, prefix string) {
	var newMap = make([]Slot, slotCount)
	for i := range newMap {
		newMap[i] = Slot{
			ID:      prefix + strconv.Itoa(i+1),
			Blocked: false,
			userID:  "",
		}
	}
	MapList.Lock()
	MapList.M[name] = Map{
		Map: newMap,
		mu:  sync.Mutex{},
	}
	MapList.Unlock()

}

//Случайное заполнение
/*func (m *Map) Randomize() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for x, cur := range m.Map {
		if cur.Blocked != 0 {
			rand.Seed(time.Now().UnixNano())
			m.Map[x].Blocked = rand.Intn(2) + 1
			time.Sleep(10)
		}
	}
}

*/

/*func (m *Map) GetBestSlot() *Slot {
	if m.FreeSlot() == 0 {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	var bX, bY, bZ int
	var distX, distY = 1000, 1000
	for x, cur := range m.Map {
		for y, cur := range cur {
			for z, cur := range cur {
				if cur.Blocked == 1 {
					for _, exits := range m.Exits {
						if dist(x, exits[0]) == 0 && dist(y, exits[1]) == 0 {
							return &m.Map[x][y][z]
						}
						if dist(x, exits[0]) <= distX && dist(y, exits[1]) <= distY {

							distX = dist(x, exits[0])
							distY = dist(y, exits[1])
							bX = x
							bY = y
							bZ = z
						}
					}
				}
			}
		}
	}
	return &m.Map[bX][bY][bZ]
}

*/
func (m *Map) GetBestSlot() *Slot {
	if m.FreeSlot() == 0 {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	rand.Seed(time.Now().UnixNano())
	return &m.Map[rand.Intn(len(m.Map))]
}

func (s *Slot) Block(ID string) {
	s.Blocked = true
	s.userID = ID
}

func (s *Slot) unBlock() {
	s.Blocked = false
	s.userID = ""
}

func (m *Map) GetSlotById(ID string) *Slot {
	for x, slot := range m.Map {
		if slot.ID == ID {
			return &m.Map[x]
		}
	}
	return nil
}

func (u Users) GetIDByToken(token string) string {

	for _, user := range u {
		if user.Token == token {
			return user.ID
		}
	}
	return ""
}
