package uniquemutex

import "sync"

type UqMutex struct {
	m map[string]*sync.RWMutex
	sync.Mutex
}

func NewUqMutex() *UqMutex {
	return &UqMutex{
		m: make(map[string]*sync.RWMutex),
	}
}

func (u *UqMutex) GetMutex(key string) *sync.RWMutex {
	u.Lock()
	defer u.Unlock()
	if _, ok := u.m[key]; !ok {
		u.m[key] = &sync.RWMutex{}
	}

	return u.m[key]
}
