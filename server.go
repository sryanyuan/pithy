package pithy

import (
	"errors"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

// apiVersionMatcher can get the api version from the path
const apiVersionMatcher = "/v{version:[0-9.]+}"

// constant errors
var (
	ErrInvalidServerRouter  = errors.New("Invalid server router")
	ErrInvalidListenAddress = errors.New("Invalid listen address")
)

type HTTPServer struct {
	listener net.Listener
	wg       sync.WaitGroup
	svr      *http.Server
	router   SwappableHandler
}

func NewHTTPServer(address string) *HTTPServer {
	svr := &HTTPServer{
		svr: &http.Server{
			Addr: address,
		},
	}
	svr.svr.Handler = &svr.router
	return svr
}

// Serve with the speicified address
func (s *HTTPServer) Serve() error {
	// Is router already set ?
	if nil == s.router.router {
		return ErrInvalidServerRouter
	}
	if "" == s.svr.Addr {
		return ErrInvalidListenAddress
	}
	ls, err := net.Listen("tcp", s.svr.Addr)
	if nil != err {
		return err
	}

	s.listener = ls
	s.wg.Add(1)

	go s.serve()
	return nil
}

// Stop stops the server
func (s *HTTPServer) Stop() {
	s.listener.Close()
	s.wg.Wait()
}

// AppendRouters create mux and update handler of server
func (s *HTTPServer) AppendRouters(rts ...APIRouter) {
	if nil == s.router.router {
		s.router.router = mux.NewRouter()
	}

	for _, v := range rts {
		fn := wrapHTTPHandler(v.Handler())
		Debugf("Register http router %s|%s", v.Path(), v.Method())
		s.router.router.Path(apiVersionMatcher + v.Path()).Methods(v.Method()).Handler(fn)
	}
}

func (s *HTTPServer) serve() {
	Infoln("HTTP serve @ ", s.listener.Addr().String())
	s.svr.Handler = &s.router
	if err := s.svr.Serve(s.listener); !strings.Contains(err.Error(), "use of closed network connection") {
		Errorln("HTTP server stop serve with error : ", err)
	}
}
