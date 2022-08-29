package zmatch

import (
	"testing"
)

func TestPKMatch(t *testing.T) {
	for i := 0; i < 10; i++ {
		_, _ = PKMatch(&PKMatchRequest{
			Mode: "test",
			Players: []*Player{
				{
					ID: "155211",
				},
			},
			MaxPlayers: 5,
			MinPlayers: 5,
		})
	}
}
