package pithy

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// SwappableHandler is a http.Handler thats allows you to swap mux during working
type SwappableHandler struct {
	lock   sync.RWMutex
	router *mux.Router
}

// Swap swaps the router of the handler
func (h *SwappableHandler) Swap(r *mux.Router) {
	h.lock.Lock()
	h.router = r
	h.lock.Unlock()
}

// ServeHTTP implement the http.Handler interface
func (h *SwappableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.lock.RLock()
	router := h.router
	h.lock.RUnlock()
	router.ServeHTTP(w, r)
}

// APIRouter defines an interface to determine which handler is called by path and method
type APIRouter interface {
	Path() string
	Method() string
	Handler() APIFunc
}