package handler

import (
	"errors"
	"fmt"
	"sync"
)

// AppBinding defines the behaviour the handler expects from the Wails App backend.
type AppBinding interface {
	Greet(name string) string
	GetImgSrc() string
}

// Handles exposes API entrypoints that delegate to the bound App implementation.
type Handles struct {
	mu  sync.RWMutex
	app AppBinding
}

var (
	handlesOnce      sync.Once
	handlesSingleton *Handles
)

// Instance returns the singleton Handles instance, creating it on first use.
func Instance() *Handles {
	handlesOnce.Do(func() {
		handlesSingleton = &Handles{}
	})

	return handlesSingleton
}

// Init binds the provided App to the singleton Handles instance.
func Init(app AppBinding) *Handles {
	h := Instance()
	fmt.Println("初始化绑定")
	if app != nil {
		h.BindApp(app)
	}

	return h
}

// BindApp replaces the current App binding on the handler singleton.
func (h *Handles) BindApp(app AppBinding) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.app = app
}

// Greet proxies the call to the bound App implementation.
func (h *Handles) Greet(name string) (string, error) {
	fmt.Println("handles.go Greet")
	app, err := h.currentApp()
	if err != nil {
		return "", err
	}

	return app.Greet(name), nil
}

// GetImgSrc proxies the call to the bound App implementation.
func (h *Handles) GetImgSrc() (string, error) {
	app, err := h.currentApp()
	if err != nil {
		return "", err
	}

	return app.GetImgSrc(), nil
}

func (h *Handles) currentApp() (AppBinding, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.app == nil {
		return nil, errors.New("handler: app binding is not configured")
	}

	return h.app, nil
}
