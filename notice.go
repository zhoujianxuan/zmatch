package zmatch

type MatchStartNotice func(room *Room)

var needNotices []MatchStartNotice

func RegisterMatchStartNotice(f MatchStartNotice) {
	needNotices = append(needNotices, f)
}

func MatchStart(room *Room) {
	for _, f := range needNotices {
		f(room)
	}
}
