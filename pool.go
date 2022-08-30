package zmatch

import "errors"

var (
	ErrNotFound = errors.New("not found")
)

type PoolService interface {
	// RPop Dequeue (isolation based on key), to ensure that there will be no room in the matching pool after dequeuing
	RPop(key string) (*Room, error)
	// LPush join the team
	LPush(key string, room *Room) error
	// LRange Dequeue in batches
	LRange(key string) ([]*Room, error)
	// Clean All information about cleaning the room
	Clean(key string, room *Room) error

	SaveRoomPlayers(roomId string, players []*Player) error

	GetRoomPlayers(roomId string) ([]*Player, error)

	DelRoomPlayers(roomId string, userIds []string) error

	// PlayerSaveRoom Players save their room information
	PlayerSaveRoom(player *Player, room *Room) error
	// PlayerDelRoom Player removes their room information
	PlayerDelRoom(player *Player, roomId string) error
	// PlayerGetRoom Players get their own room information
	PlayerGetRoom(player *Player) ([]*Room, error)
}

func GetPoolService() PoolService {
	return memoryPoolService
}
