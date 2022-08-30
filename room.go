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
	ID            string    `json:"id"`
	Mode          string    `json:"mode"`
	Players       []*Player `json:"players"`
	MaxPlayers    int       `json:"maxPlayers"`
	MinPlayers    int       `json:"minPlayers"`
	LastStartTime int64     `json:"last_start_time"`
	// ...
}

func NewRoom(id string) *Room {
	return &Room{ID: id}
}

func (r *Room) GetPlayerCount() int {
	return len(r.Players)
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
	for _, player := range r.Players {
		playerStr += fmt.Sprintf(" %+v ", player)
	}
	return fmt.Sprintf("id:%s, players:%s, mode:%s, minPlayers:%d, maxPlayers:%d",
		r.ID, playerStr, r.Mode, r.MinPlayers, r.MaxPlayers)
}
