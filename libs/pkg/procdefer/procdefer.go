package procdefer

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	handlers []func()
	mu       sync.RWMutex
)

func init() {
	go start()
}

// AddDeferFunc description
// @param f
func AddDeferFunc(f func()) {
	mu.Lock()
	handlers = append(handlers, f)
	mu.Unlock()
}

// Clear description
func Clear() {
	mu.Lock()
	handlers = []func(){}
	mu.Unlock()
}

func start() {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	for range sigch {
		mu.RLock()
		for _, handler := range handlers {
			handler()
		}
		mu.RUnlock()
		os.Exit(0)
	}
}
