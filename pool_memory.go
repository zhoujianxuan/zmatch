package zmatch

type MemoryPoolService struct {
	PoolMap       map[string]*Queue
	PlayerRoomMap map[string][]*Room
}

var memoryPoolService = &MemoryPoolService{PoolMap: map[string]*Queue{}}

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
	queue, ok := r.PoolMap[key]
	if !ok || queue.IsEmpty() {
		return nil, ErrNotFound
	}

	return queue.Pop(), nil
}

func (r *MemoryPoolService) LPush(key string, room *Room) error {
	_, ok := r.PoolMap[key]
	if !ok {
		r.PoolMap[key] = &Queue{}
	}

	r.PoolMap[key].Push(room)
	return nil
}

func (r *MemoryPoolService) LRange(key string) ([]*Room, error) {
	return nil, nil
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
