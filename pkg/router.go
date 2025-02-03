//go:generate mockgen -source=router.go -package=mocks -destination=mocks/router.go
package pkg

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type Router interface {
	//  Starts the router to handle incoming requests, returning an error if it fails.
	Serve(ctx context.Context) error
	// Gracefully shuts down the router, returning an error if the operation fails.
	Shutdown(ctx context.Context) error
	// Registers a GET route with a path, action name, handler, and optional middlewares.
	GET(path string, actionName string, handler Handler, middlewares ...Middleware)
	// Registers a POST route with a path, action name, handler, and optional middlewares.
	POST(path string, actionName string, handler Handler, middlewares ...Middleware)
	// Registers a PUT route with a path, action name, handler, and optional middlewares.
	PUT(path string, actionName string, handler Handler, middlewares ...Middleware)
	// Registers a DELETE route with a path, action name, handler, and optional middlewares.
	DELETE(path string, actionName string, handler Handler, middlewares ...Middleware)
	// Converts an `http.HandlerFunc` into a `Handler` for compatibility.
	WrapHandler(handler http.HandlerFunc) Handler
}

// A function that processes an HTTP request using a `Context` object.
type Handler func(Context)

// A function that wraps a `Handler` to add pre/post-processing logic.
// Middlewares are executed in reverse order (last one first).
type Middleware func(actionName string, next Handler) Handler

type Context interface {
	// Parses the request body into a struct or variable.
	ParseBody(v any) error
	// Validates a struct and returns a list of validation errors.
	ValidateStruct(v any) []string
	// Validates a single variable against a tag and returns an error message.
	ValidateVar(v any, tag string) string
	// Returns the underlying `context.Context` of the request.
	GetRequestContext() context.Context
	// Retrieves a path parameter by its key.
	GetParam(key string) string
	// Retrieves a query parameter by its key.
	GetQueryValue(key string) string
	// Retrieves a header value by its key.
	GetHeader(key string) string
	// Stores a value in the context for later retrieval.
	SaveValue(key string, value any)
	// Retrieves a stored value from the context by its key.
	GetValue(key string) any
	// Retrieves a form value by its key.
	GetFormValue(key string) string
	// Opens an uploaded file from a form by its key, returning the file and its metadata.
	OpenFormFile(key string) (File, FileInfo, error)
	// Saves a cookie with a name, value, and expiration time.
	SaveCookie(name string, value string, expiresIn time.Duration)
	// Retrieves a cookie value by its name.
	GetCookie(name string) string
	// Starts a tracing span and returns the updated context and span object.
	StartSpan() (context.Context, Span)
	// Redirects the request to a specified URL with a status code.
	Redirect(url string, code int)
	// Serves a file to the client.
	ServeFile(filename string)
	// Sends a response to the client using a `Response` object.
	Respond(r Response)
}

type File interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

type FileInfo interface {
	FileName() string
	FileSize() int64
	ContentType() string
}

type Response interface {
	GetCode() int
}

type Validator interface {
	GetTag() string
	GetFunc() validator.FuncCtx
}
