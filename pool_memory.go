package zmatch

import "sync"

type MemoryPoolService struct {
	PoolMap sync.Map
}

var memoryPoolService = &MemoryPoolService{}

//Queue A FIFO queue.
type Queue []*Room

//Push the element into the queue.
func (q *Queue) Push(v *Room) {
	*q = append(*q, v)
}

//Pop element from head.
func (q *Queue) Pop() *Room {
	head := (*q)[0]
	*q = (*q)[1:]
	return head
}

//IsEmpty Returns if the queue is empty or not.
func (q *Queue) IsEmpty() bool {
	return len(*q) == 0
}

func (r *MemoryPoolService) RPop(key string) (*Room, error) {
	data, ok := r.PoolMap.Load(key)
	if !ok {
		return nil, ErrNotFound
	}
	queue := data.(*Queue)
	if queue.IsEmpty() {
		return nil, ErrNotFound
	}

	return queue.Pop(), nil
}

func (r *MemoryPoolService) LPush(key string, room *Room) error {
	if data, ok := r.PoolMap.Load(key); !ok {
		queue := &Queue{}
		queue.Push(room)
		r.PoolMap.Store(key, queue)
	} else {
		queue := data.(*Queue)
		queue.Push(room)
	}
	return nil
}

func (r *MemoryPoolService) LRange(key string) ([]*Room, error) {
	data, ok := r.PoolMap.Load(key)
	if !ok {
		return nil, ErrNotFound
	}

	queue := data.(*Queue)
	if queue.IsEmpty() {
		return nil, ErrNotFound
	}

	var rooms []*Room
	var i int
	for !queue.IsEmpty() && i < 1000 {
		i++
		rooms = append(rooms, queue.Pop())
	}
	return rooms, nil
}

func (r *MemoryPoolService) Clean(key string, room *Room) error {
	return nil
}

func (r *MemoryPoolService) PlayerSaveRoom(player *Player, room *Room) error {
	return nil
}

func (r *MemoryPoolService) PlayerDelRoom(player *Player, roomId string) error {
	return nil
}

func (r *MemoryPoolService) PlayerGetRoom(player *Player) ([]*Room, error) {
	return nil, nil
}
