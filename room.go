package zmatch

import (
	"fmt"
	"time"
)

type Player struct {
	ID string `json:"id"`
	// ...
}

type Room struct {
	ID            string `json:"id"`
	Mode          string `json:"mode"`
	MaxPlayers    int    `json:"maxPlayers"`
	MinPlayers    int    `json:"minPlayers"`
	LastStartTime int64  `json:"last_start_time"`
	// ...
}

func NewRoom(id string) *Room {
	return &Room{ID: id}
}

func (r *Room) GetPlayerCount() int {
	service := GetPoolService()
	players, _ := service.GetRoomPlayers(r.ID)
	return len(players)
}

func (r *Room) CanStart() bool {
	if r.GetPlayerCount() >= r.MaxPlayers {
		return true
	}
	if r.GetPlayerCount() >= r.MinPlayers && time.Now().UnixNano() >= r.LastStartTime {
		return true
	}
	return false
}

func (r *Room) ToString() string {
	var playerStr string
	service := GetPoolService()
	players, _ := service.GetRoomPlayers(r.ID)
	for _, player := range players {
		playerStr += fmt.Sprintf(" %+v ", player)
	}
	return fmt.Sprintf("id:%s, players:%s, mode:%s, minPlayers:%d, maxPlayers:%d",
		r.ID, playerStr, r.Mode, r.MinPlayers, r.MaxPlayers)
}
