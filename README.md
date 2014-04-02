#Simple Web Server

````
import "github.com/JGets/sws"
````

##About

`sws` is a simple HTTP server backend to use in making web servers in Go.

The package is built around the [net/http](http://golang.org/pkg/net/http/) package, but adds the ability to use regular expressons for URL pattern matching.

As well, you can specify functions for the server to call when a 404 not found, or a 500 internal error are encountered, allowing for easy implementation of customized error pages.

Additionally, `sws` has built in static file serving. For any request that does not match to a pattern, the server attempts to serve the file from the directory specified in the server's `StaticDir` variable. In the default server, this variable defaults to `static/` in the programs current working directory. By default, static files are served with a header setting the cache max age to 1 week, however this cache value can be changed by setting the server's `StaticFileCacheParam`, or using the `SetStaticFileCacheMaxAge()` method.

##Documentation

###Variables

````
var DefaultServer = NewServer()
````

###Functions

####func Handle
````
func Handle(pattern string, handler http.Handler)
````
Adds a route to handle the given `pattern` using the given `handler`, to the `DefaultServer`.

####func HandleFunc
````
func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
````
Adds a route to handle the given `pattern` using the given `handler` function, to the `DefaultServer`.

####func Run
````
func Run(addr string) error
````
Runs the `DefaultServer` on the given address string.

####func SetNotFoundHandler
````
func SetNotFoundHandler(handler func(http.ResponseWriter, *http.Request))
````
Sets the `DefaultServer`'s `NotFoundHandler` to the given function.

####func SetInternalErrorHandler
````
func SetInternalErrorHandler(handler func(http.ResponseWriter, *http.Request, error))
````
Sets the `DefaultServer`'s `InternalErrorHandler` to the given function.

####func ParseParams
````
func ParseParams(r *http.Request) map[string]string
````
Attempts to parse the form parameters of the given request object into an easy to use map. If a paremeter has multiple values, only the first is populated into the map.

####func DefaultNotFound
````
func DefaultNotFound(w http.ResponseWriter, r *http.Request)
````
A simple handler function for when resources are not found. Sends a plaintext "404 Not Found" response with appropriate headers.

####func DefaultInternalError
````
func DefaultInternalError(w http.ResponseWriter, r *http.Request, err error)
````
A simple handler function for when an error occurs in the server. Sends a plaintext "500 Internal Server Error" along with the result of `err.Error()` response, along with the appropriate headers.


###Types

####SimpleWebServer
````
type SimpleWebServer struct {
	NotFoundHandler func (http.ResponseWriter, *http.Request)
		//A function to handle when resources are not found
	InternalErrorHandler func(http.ResponseWriter, *http.Request, error)
		//A function to handle internal server errors
	StaticDir string
		//The path to the directory from which static files are to be served
	StaticFileCacheParam string
		//The value which is to be set for the header's 'Cache-Control:' parameter for static files
	//contains other unexported fields
}
````

####func NewServer
````
func NewServer() *SimpleWebServer
````
Instantiates and returns a pointer to a new `SimpleWebServer` object with all the default values set.

####func (*SimpleWebServer) Handle
````
func (s *SimpleWebServer) Handle(pattern string, handler http.Handler)
````
Adds a route to the server's `Router` to handle the given `patern` using the given `handler`.

####func (*SimpleWebServer) HandleFunc
````
func (s *SimpleWebServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
````
Adds a route to the server's `Router` to handle the given `patern` using the given `handler` function.

####func (*SimpleWebServer) Run
````
func (s *SimpleWebServer) Run(addr string) error
````
Runs the server on the TCP network address `addr`

####func (*SimpleWebServer) SetStaticFileCacheMaxAge
````
func (s *SimpleWebServer) SetStaticFileCacheMaxAge(age int)
````
A convenience method to set the `StaticFileCacheParam` of the server to a maximum age of `age`, in seconds



