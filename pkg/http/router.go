package http

import (
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/shahzodshafizod/gocloud/pkg"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

type server struct {
	addr      string
	router    *http.ServeMux
	validator *validator.Validate
	server    *http.Server
	tracer    pkg.Tracer
}

func New(addr string, tracer pkg.Tracer, validators ...pkg.Validator) pkg.Router {
	validator := validator.New()
	for _, v := range validators {
		if v == nil {
			continue
		}
		validator.RegisterValidationCtx(v.GetTag(), v.GetFunc())
	}

	return &server{
		addr:      addr,
		router:    http.NewServeMux(),
		validator: validator,
		tracer:    tracer,
	}
}

func (s *server) Serve(context.Context) error {
	s.server = &http.Server{
		Addr:              s.addr,
		Handler:           s.router,
		ReadHeaderTimeout: time.Second * 3,
		ReadTimeout:       time.Second * 15,
		WriteTimeout:      time.Second * 15,
		IdleTimeout:       time.Second * 30,
	}
	return s.server.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return errors.New("server is not set")
	}
	return s.server.Shutdown(ctx)
}

func (s *server) GET(path string, name string, next pkg.Handler, mws ...pkg.Middleware) {
	s.handle(http.MethodGet, path, name, next, mws...)
}

func (s *server) POST(path string, name string, next pkg.Handler, mws ...pkg.Middleware) {
	s.handle(http.MethodPost, path, name, next, mws...)
}

func (s *server) PUT(path string, name string, next pkg.Handler, mws ...pkg.Middleware) {
	s.handle(http.MethodPut, path, name, next, mws...)
}

func (s *server) DELETE(path string, name string, next pkg.Handler, mws ...pkg.Middleware) {
	s.handle(http.MethodDelete, path, name, next, mws...)
}

func (s *server) WrapHandler(handler http.HandlerFunc) pkg.Handler {
	return pkg.Handler(func(ctx pkg.Context) {
		customCtx := ctx.(*customContext)
		handler.ServeHTTP(customCtx.response, customCtx.request)
	})
}

func (s *server) handle(method string, path string, name string, next pkg.Handler, mws ...pkg.Middleware) {

	for _, mw := range mws {
		next = mw(name, next)
	}

	prefix, maskedPath := maskPath(path)

	ha := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cctx := &customContext{
			name:        name,
			request:     r,
			response:    w,
			headers:     r.Header,
			validator:   s.validator,
			ctx:         r.Context(),
			params:      makeParams(prefix, maskedPath, r.URL.Path),
			queryValues: r.URL.Query(),
			values:      make(map[string]any),
			tracer:      s.tracer,
		}
		r.ParseMultipartForm(defaultMaxMemory)
		cctx.form = r.MultipartForm

		if r.Method != method {
			statusCode := http.StatusMethodNotAllowed
			cctx.Respond(&response{Code: statusCode, Message: http.StatusText(statusCode)})
			return
		}

		next(cctx)
	})

	s.router.HandleFunc(prefix, ha)
}

func maskPath(path string) (string, string) {
	// create regex for named params
	var firstIndex = -1
	for strings.Contains(path, "/:") {
		indexFrom := strings.Index(path, "/:") + 1
		if firstIndex < 0 {
			firstIndex = indexFrom
		}
		indexTo := indexFrom + strings.Index(path[indexFrom:], "/")
		if indexTo < indexFrom {
			indexTo = len(path)
		}
		if indexFrom+1 == indexTo {
			break
		}
		path = strings.Replace(path, path[indexFrom:indexTo], `(?P<`+path[indexFrom:indexTo][1:]+`>[\w-\.]+)`, 1)
	}

	var prefix = path
	if firstIndex >= 0 {
		prefix = path[:firstIndex]
	}

	return prefix, path
}

func makeParams(prefix, maskedPath, path string) map[string]string {
	params := make(map[string]string)
	if prefix != maskedPath {
		// parse param values from path
		re, err := regexp.Compile(maskedPath)
		if err != nil {
			return params
		}
		if matches := re.FindStringSubmatch(path); len(matches) > 0 {
			for index, name := range re.SubexpNames() {
				if name == "" {
					continue
				}
				params[name] = matches[index]
			}
		}
	}
	return params
}

type customContext struct {
	name        string
	request     *http.Request
	response    http.ResponseWriter
	headers     http.Header
	validator   *validator.Validate
	ctx         context.Context
	params      map[string]string
	queryValues url.Values
	values      map[string]any
	form        *multipart.Form
	tracer      pkg.Tracer
}

func (c *customContext) ParseBody(v any) error {
	return json.NewDecoder(c.request.Body).Decode(v)
}

func (c *customContext) ValidateStruct(v any) []string {
	err := c.validator.StructCtx(c.ctx, v)
	if err != nil {
		return strings.Split(err.Error(), "\n")
	}
	return []string{}
}

func (c *customContext) ValidateVar(v any, tag string) string {
	err := c.validator.VarCtx(c.ctx, v, tag)
	if err != nil {
		return err.Error()
	}
	return ""
}

func (c *customContext) GetRequestContext() context.Context {
	return c.ctx
}

func (c *customContext) GetParam(name string) string {
	return c.params[name]
}

func (c *customContext) GetQueryValue(key string) string {
	return c.queryValues.Get(key)
}

func (c *customContext) GetHeader(key string) string {
	return c.headers.Get(key)
}

func (c *customContext) SaveValue(key string, value any) {
	c.values[key] = value
}

func (c *customContext) GetValue(key string) any {
	return c.values[key]
}

func (c *customContext) GetFormValue(key string) string {
	if c.form != nil {
		if values := c.form.Value[key]; len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

func (c *customContext) OpenFormFile(key string) (pkg.File, pkg.FileInfo, error) {

	if c.form == nil {
		return nil, nil, errors.New("no files parsed")
	}

	fileHeaders, found := c.form.File[key]
	if !found {
		return nil, nil, errors.New("no files found with the key " + key)
	}
	if len(fileHeaders) <= 0 {
		return nil, nil, errors.New("not found")
	}

	header := fileHeaders[0]
	file, err := header.Open()
	if err != nil {
		return nil, nil, err
	}

	info := &fileInfo{
		name:        header.Filename,
		size:        header.Size,
		contentType: header.Header.Get("Content-Type"),
	}

	return file, info, nil
}

type fileInfo struct {
	name        string
	size        int64
	contentType string
}

func (f *fileInfo) FileName() string    { return f.name }
func (f *fileInfo) FileSize() int64     { return f.size }
func (f *fileInfo) ContentType() string { return f.contentType }

func (c *customContext) SaveCookie(name string, value string, expiresIn time.Duration) {
	http.SetCookie(c.response, &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(expiresIn),
		Secure:   true,
		HttpOnly: true,
	})
}

func (c *customContext) GetCookie(name string) string {
	cookie, _ := c.request.Cookie(name)
	if cookie != nil {
		return cookie.Value
	}
	return ""
}

func (c customContext) StartSpan() (context.Context, pkg.Span) {
	return c.tracer.StartFromContext(c.ctx, c.name)
}

func (c *customContext) Redirect(url string, code int) {
	if code == 0 {
		code = http.StatusFound
	}
	http.Redirect(c.response, c.request, url, code)
}

func (c *customContext) ServeFile(filename string) {
	http.ServeFile(c.response, c.request, filename)
}

func (c *customContext) Respond(r pkg.Response) {
	if r == nil {
		return
	}
	c.response.Header().Set("Content-Type", "application/json")
	code := r.GetCode()
	if code == 0 {
		code = http.StatusNotImplemented
	}
	c.response.WriteHeader(code)
	err := json.NewEncoder(c.response).Encode(r)
	if err != nil {
		c.response.WriteHeader(http.StatusInternalServerError)
	}
}

type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r *response) GetCode() int { return r.Code }

type mapCarrier http.Header

var _ pkg.TextMapCarrier = mapCarrier{}

func (m mapCarrier) Get(key string) string {
	if value, ok := m[key]; ok && len(value) > 0 {
		return value[0]
	}
	return ""
}

func (m mapCarrier) Set(key string, value string) {
	m[key] = []string{value}
}

func (m mapCarrier) Keys() []string {
	var keys = make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
