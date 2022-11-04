package util

import "sync"

type SyncMap struct {
	sync.Map
}

func (s *SyncMap) Size() int {
	i := 0
	s.Range(func(key, value any) bool {
		i++
		return true
	})
	return i
}
