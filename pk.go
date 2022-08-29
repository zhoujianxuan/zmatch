package zmatch

import (
	"log"
	"time"
)

type PKMatchRequest struct {
	Mode       string    `json:"mode"`
	Players    []*Player `json:"players"`
	MaxPlayers int       `json:"maxPlayers"`
	MinPlayers int       `json:"minPlayers"`
	// ...
}

type PKMatchResponse struct {
	RoomID string `json:"room_id"`
}

func PKMatch(request *PKMatchRequest) (*PKMatchResponse, error) {
	id, _ := time.Now().MarshalText()

	room := NewRoom(string(id))
	room.Mode = request.Mode
	room.Players = request.Players
	room.MaxPlayers = request.MaxPlayers
	room.MinPlayers = request.MinPlayers
	room.LastStartTime = time.Now().Add(time.Minute).UnixNano()

	service := GetPoolService()
	room, err := GetSuitableRoom(room, service)
	if err != nil {
		return nil, err
	}

	if room.CanStart() {
		// start
		log.Println("start")
	} else {
		for _, player := range room.Players {
			err = service.PlayerSaveRoom(player, room)
			if err != nil {
				return nil, err
			}
		}

		err = service.LPush(room.Mode, room)
		if err != nil {
			return nil, err
		}
	}
	return &PKMatchResponse{RoomID: room.ID}, nil
}

func JudgeSuitableRoom(room, pRoom *Room) bool {
	if room.GetPlayerCount() >= room.MaxPlayers {
		return false
	}
	if pRoom.GetPlayerCount() >= pRoom.MaxPlayers {
		return false
	}

	total := pRoom.GetPlayerCount() + room.GetPlayerCount()
	if total <= pRoom.MaxPlayers && total <= room.MaxPlayers {
		return true
	}
	return false
}

func GetSuitableRoom(room *Room, service PoolService) (*Room, error) {
	for i := 0; i < 3; i++ {
		pRoom, err := service.RPop(room.Mode)
		if err == ErrNotFound {
			break
		} else if err != nil {
			return nil, err
		}
		if JudgeSuitableRoom(room, pRoom) {
			pRoom.Players = append(pRoom.Players, room.Players...)
			for _, player := range room.Players {
				err = service.PlayerSaveRoom(player, pRoom)
				if err != nil {
					return nil, err
				}
			}
			return pRoom, nil
		}

		_ = service.LPush(room.Mode, room)
	}
	return room, nil
}
