package zmatch

import (
	"fmt"
	"math/rand"
	"time"
)

type PKMatchRequest struct {
	Mode           string        `json:"mode"`
	Players        []*Player     `json:"players"`
	MaxPlayers     int           `json:"max_players"`
	MinPlayers     int           `json:"min_players"`
	AfterStartTime time.Duration `json:"after_start_time"`
	// ...
}

type PKMatchResponse struct {
	Room *Room `json:"room"`
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func randStr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return string(b)
}

//PKMatch 匹配入口
func PKMatch(request *PKMatchRequest) (*PKMatchResponse, error) {
	id := fmt.Sprintf("%d%s", time.Now().UnixNano(), randStr(8))

	room := NewRoom(id)
	room.Mode = request.Mode
	room.Players = request.Players
	room.MaxPlayers = request.MaxPlayers
	room.MinPlayers = request.MinPlayers
	room.LastStartTime = time.Now().Add(request.AfterStartTime).UnixNano()

	service := GetPoolService()
	room, err := GetSuitableRoom(room, service)
	if err != nil {
		return nil, err
	}

	if room.CanStart() {
		// start
		MatchStart(room)
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
	return &PKMatchResponse{Room: room}, nil
}

func judgeSuitableRoom(room, pRoom *Room) bool {
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
		if judgeSuitableRoom(room, pRoom) {
			pRoom.Players = append(pRoom.Players, room.Players...)
			for _, player := range room.Players {
				err = service.PlayerSaveRoom(player, pRoom)
				if err != nil {
					return nil, err
				}
			}
			return pRoom, nil
		}
		_ = service.LPush(pRoom.Mode, pRoom)
	}
	return room, nil
}
