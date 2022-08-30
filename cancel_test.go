package zmatch

import (
	"fmt"
	"testing"
	"time"
)

func TestPKMatchCancel(t *testing.T) {
	CycleCheck(time.Second)
	RegisterMatchStartNotice(func(room *Room) {
		fmt.Println(room.ToString())
	})
	request := &PKMatchRequest{
		Mode:       "test",
		MaxPlayers: 4,
		MinPlayers: 3,
	}
	r := 3
	for j := 0; j < r; j++ {
		request.Players = append(request.Players, &Player{
			ID: randStr(6),
		})
	}
	response, _ := PKMatch(request)
	_, _ = CancelMatch(&CancelMatchRequest{RoomID: response.Room.ID, CancelPlayerIds: []string{request.Players[0].ID}})
	request.Players = []*Player{
		{ID: randStr(8)},
	}
	_, _ = PKMatch(request)
	time.Sleep(time.Second)
}
