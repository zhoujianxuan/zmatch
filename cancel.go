package zmatch

type CancelMatchRequest struct {
	RoomID          string   `json:"room_id"`
	CancelPlayerIds []string `json:"cancel_player_ids"`
}

type CancelMatchResponse struct{}

func CancelMatch(request *CancelMatchRequest) (*CancelMatchResponse, error) {
	service := GetPoolService()
	err := service.DelRoomPlayers(request.RoomID, request.CancelPlayerIds)
	if err != nil {
		return nil, err
	}
	return &CancelMatchResponse{}, nil
}
