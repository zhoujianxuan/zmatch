package zmatch

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestPKMatch(t *testing.T) {
	CycleCheck()
	RegisterMatchStartNotice(func(room *Room) {
		fmt.Println(room.ToString())
	})

	for i := 0; i < 100; i++ {
		request := &PKMatchRequest{
			Mode:       "test",
			MaxPlayers: 4,
			MinPlayers: 3,
		}
		r := rand.Int()%4 + 1
		for j := 0; j < r; j++ {
			request.Players = append(request.Players, &Player{
				ID: randStr(6),
			})
		}
		time.Sleep(500 * time.Millisecond)
		_, _ = PKMatch(request)
	}
}
