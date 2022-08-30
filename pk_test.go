package zmatch

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestPKMatch(t *testing.T) {
	CycleCheck(time.Second)
	RegisterMatchStartNotice(func(room *Room) {
		fmt.Println(room.ToString())
	})

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
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
			time.Sleep(time.Duration(rand.Int()%10) * time.Second)
			_, _ = PKMatch(request)
			wg.Done()
		}()
	}
	wg.Wait()
	time.Sleep(time.Second)
}
