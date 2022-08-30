package zmatch

import (
	"log"
	"time"
)

func GetAllKeys() []string {
	return []string{"test"}
}

func CycleCheck(d time.Duration) {
	keys := GetAllKeys()
	for _, key := range keys {
		go func(key string) {
			defer func() {
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()
			for {
				t := time.NewTimer(d)
				select {
				case <-t.C:
					service := GetPoolService()
					rooms, err := service.LRange(key)
					if err == ErrNotFound {
						break
					}
					if err != nil {
						panic(err)
					}

					for _, room := range rooms {
						if room.CanStart() {
							MatchStart(room)
						} else {
							_ = service.LPush("test", room)
						}
					}
				}
			}
		}(key)
	}
}
