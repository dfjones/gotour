package mutex

import (
	"sync"
)

type visitMap struct {
	visitedUrls map[string]bool
	mutex       *sync.RWMutex
}

func New() *visitMap {
	return &visitMap{make(map[string]bool), new(sync.RWMutex)}
}

func (vm *visitMap) Visit(url string) {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()
	vm.visitedUrls[url] = true
}

func (vm *visitMap) IsVisited(url string) bool {
	vm.mutex.RLock()
	defer vm.mutex.RUnlock()
	_, ok := vm.visitedUrls[url]
	return ok
}

func (vm *visitMap) Close() map[string]bool {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()
	cpy := make(map[string]bool)
	for k, v := range vm.visitedUrls {
		cpy[k] = v
	}
	return cpy
}
