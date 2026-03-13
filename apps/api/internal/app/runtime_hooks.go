package app

import (
	"context"
	"sync"
)

type postWriteHookFunc func(context.Context)

var (
	hooksMu             sync.RWMutex
	registeredWriteHook postWriteHookFunc
)

func SetPostWriteHook(hook func(context.Context)) {
	hooksMu.Lock()
	defer hooksMu.Unlock()
	registeredWriteHook = hook
}

func currentPostWriteHook() postWriteHookFunc {
	hooksMu.RLock()
	defer hooksMu.RUnlock()
	return registeredWriteHook
}
