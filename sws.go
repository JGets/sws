/*
The MIT License (MIT)

Copyright (c) 2014 John Gettings

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package sws

import(
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
)

// the default server object to use
var DefaultServer = NewServer()

/*
	Adds a route to handle the given pattern using the given handler, to the DefaultServer.
*/
func Handle(pattern string, handler http.Handler) {
	DefaultServer.Handle(pattern, handler)
}

/*
	Adds a route to handle the given pattern using the gieven handler function, to the Default Server.
*/
func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request) ) {
	DefaultServer.HandleFunc(pattern, handler)
}

/*
	Runs the Default Server on the given address addr.
*/
func Run(addr string) error {
	return DefaultServer.Run(addr)
}

/*
	Sets the DefaultServer's NotFoundHandler to the given function.
*/
func SetNotFoundHandler( handler func(http.ResponseWriter, *http.Request) ) {
	DefaultServer.NotFoundHandler = handler
}

/*
	Sets the DefaultServer's InternalErrorHandler to the given function
*/
func SetInternalErrorHandler( handler func(http.ResponseWriter, *http.Request, error) ) {
	DefaultServer.InternalErrorHandler = handler
}

/*
	Attempts to parse the form parameters of the given request object into an easy to use map. If a paremeter has multiple values, only the first is populated into the map.
*/
func ParseParams(r *http.Request) map[string]string {
	r.ParseForm()

	ret := make(map[string]string)

	for k, v := range r.Form {
		if v != nil && len(v) > 0 {
			ret[k] = v[0]
		} else {
			fmt.Printf("Warning: No parameter found for key '%v'\n", k)
		}
	}
	return ret
}

/*
	A simple handler function for when resources are not found. Sends a plaintext "404 Not Found" response with appropriate headers.
*/
func DefaultNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("404 Not Found"))
}

/*
	A simple handler function for when an error occurs in the server. Sends a plaintext "500 Internal Server Error" along with the result of err.Error() response, along with the appropriate headers.
*/
func DefaultInternalError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(500)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("500 Internal Server Error\n"))
	w.Write([]byte(err.Error()))
}




type SimpleWebServer struct {
	router *regexpRouter
	NotFoundHandler func (http.ResponseWriter, *http.Request)	//A Handler function to handle when resources are not found.
	InternalErrorHandler func(http.ResponseWriter, *http.Request, error)	//A function to handle internal server errors
	StaticDir string 	//The path to the directory from which static files are to be served
	StaticFileCacheParam string 	//The value which is to be set for the header's 'Cache-Control:' parameter for static files
}


/*
	Instantiates and returns a pointer to a new SimpleWebServer object with all the default values set.
*/
func NewServer() *SimpleWebServer {
	//get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	//initialize a new SimpleWebServer object
	s := &SimpleWebServer{newRegexpRouter(), DefaultNotFound, DefaultInternalError, path.Join(wd, "static"), "max-age=604800"}	//default cache max age of 1 week
	//make sure the server's RegexpRouter object has a pointer to the server
	s.router.Server = s
	//return a pointer to the server
	return s
}

/*
	Adds a route to the server's Router to handle the given patern using the given handler.
*/
func (s *SimpleWebServer) Handle(pattern string, handler http.Handler) {
	s.router.Handle(pattern, handler)
}

/*
	Adds a route to the server's Router to handle the given patern using the given handler function.
*/
func (s *SimpleWebServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request) ){
	s.router.HandleFunc(pattern, handler)
}

/*
	Runs the server on the TCP network address addr
*/
func (s *SimpleWebServer) Run(addr string) error {

	mux := http.NewServeMux()

	mux.Handle("/", s.router)

	// http.Handle("/", s.router)
	err := http.ListenAndServe(addr, mux)
	return err
}

/*
	A convenience method to set the StaticFileCacheParam of the server to a maximum age of age, in seconds
*/
func (s *SimpleWebServer) SetStaticFileCacheMaxAge(age int) {
	s.StaticFileCacheParam = fmt.Sprintf("max-age=%v", age)
}

/*
	Attempts to serve a static file that matches the request.
*/
func (s *SimpleWebServer) attemptToServeStaticFile(w http.ResponseWriter, r *http.Request) {
	//create a path string for where the file should be
	filePath := path.Join(s.StaticDir, r.URL.Path)

	// fmt.Printf("attempting to serve file from: %v\n", filePath)

	if s.fileExists(filePath){	//if the file exists:

		// fmt.Printf("Attempting to serve static file: %v\n", r.URL.Path)


		//set the cache header field
		w.Header().Set("Cache-Control", s.StaticFileCacheParam)
		//serve the file
		http.ServeFile(w, r, filePath)

	} else {	//if the file doesn't exist

		// fmt.Printf("File Not Found: %v\n", r.URL.Path)

		//call the server's not found handler
		s.NotFoundHandler(w, r)
	}
}

/*
	determines if the file at filePath exists
*/
func (s *SimpleWebServer) fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false;
	}
	return !info.IsDir()
}




type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type regexpRouter struct {
	routes []*route
	Server *SimpleWebServer
}

/*
	Instantiates a new regexpRouter (with a nil Server) and returns a pointer to it.
*/
func newRegexpRouter() *regexpRouter {
	return &regexpRouter{make([]*route, 0, 4), nil}	//default static file cache life of 1 day
}

/*
	Adds a new route to the regexpRouter for the given pattern and Handler
*/
func (h *regexpRouter) Handle(pattern string, handler http.Handler){
	rpat, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}
	h.routes = append(h.routes, &route{rpat, handler})
}

/*
	Adds a new route to the regexpRouter for the given pattern and handler function
*/
func (h *regexpRouter) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request) ){
	rpat, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}
	h.routes = append(h.routes, &route{rpat, http.HandlerFunc(handler)})
}

/*
	Attempts to find a matching route for the request. If no route matches, and the router's Server is non-nil, it attempts to serve a static file from the Server, otherwise it serves a basic plaintext 404.
*/
func (h *regexpRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//iterate through all the routes, looking for a match
    for _, route := range h.routes {
        if route.pattern.MatchString(r.URL.Path) {
            route.handler.ServeHTTP(w, r)
            return
        }
    }
    //if the request doesn't match one of the registerd routes,
    if h.Server != nil {
    	//and if the router has a pointer to a server, attempt to serve a static file for the request
	    h.Server.attemptToServeStaticFile(w, r);
    } else {
    	//otherwise, just serve a basic plaintext 404
    	DefaultNotFound(w, r)
    }
}
