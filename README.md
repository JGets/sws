#Simple Web Server

````
import "github.com/JGets/sws"
````

##About

SWS implements a simple HTTP server backend to use in making web servers in Go.

The package is built around the [net/http](http://golang.org/pkg/net/http/) package, but adds the ability to use regular expressons for URL pattern matching.

As well, you can specify functions for the server to call when a 404 not found, or a 500 internal error are encountered, allowing for easy implementation of customized error pages.

Additionally, SWS has built in static file serving. For any request that does not match to a pattern, the server attempts to serve the file from the directory specified in the server's `StaticDir` variable. In the default server, this variable defaults to `static/` in the programs current working directory. Static files are served with a header setting the cache amx age to 1 week.

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


###Types

####SimpleWebServer
````
type SimpleWebServer struct {
	Handler *RegexpHandler
	NotFoundHandler func (http.ResponseWriter, *http.Request)
	InternalErrorHandler func(http.ResponseWriter, *http.Request, error)
	StaticDir string
	StaticFileCacheParam string
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

