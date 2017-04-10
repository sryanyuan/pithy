package pithy

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPContext struct {
	r *http.Request
	w http.ResponseWriter
}

type APIFunc func(*HTTPContext) *APIResult

func wrapHTTPHandler(fn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &HTTPContext{
			w: w,
			r: r,
		}
		result := fn(ctx)
		if nil == result {
			result = NewAPIResult(0, 0, "")
		}
		ctx.WriteAPIResult(result)
	}
}

func (c *HTTPContext) GetRequest() *http.Request {
	return c.r
}

func (c *HTTPContext) GetResponseWriter() http.ResponseWriter {
	return c.w
}

func (c *HTTPContext) WriteAPIResult(result *APIResult) {
	if jsonBytes, err := json.Marshal(result); nil != err {
		panic(err)
	} else {
		if 0 != result.StatusCode {
			c.w.WriteHeader(result.StatusCode)
		}
		c.w.Write(jsonBytes)
	}
}

func (c *HTTPContext) GetPathVar(key string) string {
	vars := mux.Vars(c.r)
	v, ok := vars[key]
	if !ok {
		return ""
	}
	return v
}

func (c *HTTPContext) GetPathVersion() string {
	return c.GetPathVar("version")
}
