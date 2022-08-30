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
	room.MaxPlayers = request.MaxPlayers
	room.MinPlayers = request.MinPlayers
	room.LastStartTime = time.Now().Add(request.AfterStartTime).UnixNano()

	service := GetPoolService()
	room, err := GetSuitableRoom(room, request.Players, service)
	if err != nil {
		return nil, err
	}

	if room.CanStart() {
		// start
		MatchStart(room)
	} else {
		_ = service.SaveRoomPlayers(room.ID, request.Players)
		for _, player := range request.Players {
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

func judgeSuitableRoom(userCount, maxPlayers int, pRoom *Room) bool {
	if userCount >= maxPlayers {
		return false
	}

	pRoomUserCount := pRoom.GetPlayerCount()
	if pRoomUserCount >= pRoom.MaxPlayers {
		return false
	}

	total := pRoomUserCount + userCount
	if total <= pRoom.MaxPlayers && total <= maxPlayers {
		return true
	}
	return false
}

func GetSuitableRoom(room *Room, players []*Player, service PoolService) (*Room, error) {
	for i := 0; i < 3; i++ {
		pRoom, err := service.RPop(room.Mode)
		if err == ErrNotFound {
			break
		} else if err != nil {
			return nil, err
		}
		if judgeSuitableRoom(len(players), room.MaxPlayers, pRoom) {
			_ = service.SaveRoomPlayers(pRoom.ID, players)
			return pRoom, nil
		}
		_ = service.LPush(pRoom.Mode, pRoom)
	}
	return room, nil
}
