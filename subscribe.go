package dfuse

import (
	gocache "github.com/patrickmn/go-cache"
)

// 订阅管理
type Subscribe struct {
	cache    *gocache.Cache
	reqCache map[string][]byte
}
