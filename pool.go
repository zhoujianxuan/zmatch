package zmatch

import "errors"

var (
	ErrNotFound = errors.New("not found")
)

type PoolService interface {
	// RPop 出队(根据key进行隔离)，要保证出队之后匹配池中不会存在房间
	RPop(key string) (*Room, error)
	// LPush 入队
	LPush(key string, room *Room) error
	// LRange 批量出队
	LRange(key string) ([]*Room, error)
	// Clean 清理房间所有相关信息
	Clean(key string, room *Room) error

	// PlayerSaveRoom 玩家保存自己的房间信息
	PlayerSaveRoom(player *Player, room *Room) error
	// PlayerDelRoom 玩家移除自己的房间信息
	PlayerDelRoom(player *Player, roomId string) error
	// PlayerGetRoom 玩家获取自己的房间信息
	PlayerGetRoom(player *Player) ([]*Room, error)
}

func GetPoolService() PoolService {
	return memoryPoolService
}
