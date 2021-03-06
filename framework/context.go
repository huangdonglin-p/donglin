package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Context is a context manager that encapsulates the request context,
// timeout settings and other functions, similar to Gin's Context.
type Context struct {
	// Implement the function of WebService receiving and processing protocol text
	request        *http.Request
	responseWriter http.ResponseWriter
	//
	ctx     context.Context
	handler ControllerHandler

	//
	hasTimeout bool
	//
	writerMux *sync.Mutex
}

// NewContext returns a new Context
func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		writerMux:      &sync.Mutex{},
	}
}

/********** Context base function **************/

// WriterMux returns lock
func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.writerMux
}

// GetRequest returns http.Request
func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

// GetResponse returns http.ResponseWriter
func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.responseWriter
}

// SetHasTimeout sets the timeout flag
func (ctx *Context) SetHasTimeout() {
	ctx.hasTimeout = true
}

// HasTimeout returns ctx.hasTimeout
func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

/********** implement context.Context **************/

// BaseContext returns context.Context
func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

// Deadline implements context.Context.Deadline method
func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}

// Done implements context.Context.Done method
func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

// Err implements context.Context.Err method
func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

// Value implements context.Context.Value method
func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.BaseContext().Value(key)
}

/**************** query url, Get method parameters *****************/

// QueryAll returns Get parameter
// The parameters of the Gin framework are obtained,
// and finally the method initQueryCache() is called.
// It also uses c.Request.URL.Query() internally.
func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request != nil {
		return map[string][]string(ctx.request.URL.Query())
	}
	return map[string][]string{}
}

// QueryInt returns an int type value, or returns the default value
func (ctx *Context) QueryInt(key string, def int) int {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if l := len(vals); l > 0 {
			intval, err := strconv.Atoi(vals[l-1])
			if err != nil {
				return def
			}
			return intval
		}
	}
	return def
}

// QueryString returns a string type value, or returns the default value
func (ctx *Context) QueryString(key string, def string) string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if l := len(vals); l > 0 {
			return vals[l-1]
		}
	}
	return def
}

// QueryArray returns a []string type value, or returns the default value
func (ctx *Context) QueryArray(key string, def []string) []string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

/**************** POST method parameters *****************/

// FormAll returns Post parameter
// The gin framework calls the initFormCache method,
// and internally calls the http.Request.ParseMultipartForm method first
func (ctx *Context) FormAll() map[string][]string {
	if ctx.request != nil {
		return map[string][]string(ctx.request.PostForm)
	}
	return map[string][]string{}
}

// FormInt returns an int type value, or returns the default value
func (ctx *Context) FormInt(key string, def int) int {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if l := len(vals); l > 0 {
			intval, err := strconv.Atoi(vals[l-1])
			if err != nil {
				return def
			}
			return intval
		}
	}
	return def
}

// FormString returns a string type value, or returns the default value
func (ctx *Context) FormString(key string, def string) string {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if l := len(vals); l > 0 {
			return vals[l-1]
		}
	}
	return def
}

// FormArray returns a []string type value, or returns the default value
func (ctx *Context) FormArray(key string, def []string) []string {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

/******************* Data binding ********************/

//
func (ctx *Context) BindJson(obj interface{}) error {
	if ctx.request != nil {
		body, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ctx.request empty")
	}
	return nil
}

/**************** Response rendering *****************/

// Json serializes the given struct as JSON into the response body.
// The gin framework implements this function through the
// Render function under the render package. Finally set ContentType to
// "application/json; charset=utf-8", and then execute http.ResponseWriter.Write(data) to write data.
func (ctx *Context) Json(status int, obj interface{}) error {
	if ctx.HasTimeout() {
		return nil
	}
	ctx.responseWriter.Header().Set("Content-Type", "application/json")
	ctx.responseWriter.WriteHeader(status)
	byt, err := json.Marshal(obj)
	if err != nil {
		ctx.responseWriter.WriteHeader(500)
		return err
	}
	ctx.responseWriter.Write(byt)
	return nil
}

//
func (ctx *Context) HTML(status int, obj interface{}, template string) error {
	return nil
}

//
func (ctx *Context) Text(status int, obj string) error {
	return nil
}
