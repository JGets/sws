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


func Handle(pattern string, handler http.Handler) {
	DefaultServer.Handle(pattern, handler)
}

func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request) ) {
	DefaultServer.HandleFunc(pattern, handler)
}

func Run(addr string) error {
	return DefaultServer.Run(addr)
}

func SetNotFoundHandler( handler func(http.ResponseWriter, *http.Request) ) {
	DefaultServer.NotFoundHandler = handler
}

func SetInternalErrorHandler( handler func(http.ResponseWriter, *http.Request, error) ) {
	DefaultServer.InternalErrorHandler = handler
}


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





type SimpleWebServer struct {
	Handler *RegexpHandler
	NotFoundHandler func (http.ResponseWriter, *http.Request)
	InternalErrorHandler func(http.ResponseWriter, *http.Request, error)
	StaticDir string
	StaticFileCacheParam string
}


func NewServer() *SimpleWebServer {
	//get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	//initialize a new SimpleWebServer object
	s := &SimpleWebServer{NewRegexpHandler(), DefaultNotFound, DefaultInternalError, path.Join(wd, "static"), "max-age=604800"}	//default cache max age of 1 week
	//make sure the server's RegexpHandler object has a pointer to the server
	s.Handler.Server = s
	//return a pointer to the server
	return s
}


func (s *SimpleWebServer) Handle(pattern string, handler http.Handler) {
	s.Handler.Handle(pattern, handler)
}

func (s *SimpleWebServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request) ){
	s.Handler.HandleFunc(pattern, handler)
}

func (s *SimpleWebServer) Run(addr string) error {
	http.Handle("/", s.Handler)
	err := http.ListenAndServe(addr, nil)
	return err
}


func (s *SimpleWebServer) SetStaticFileCacheMaxAge(age int) {
	s.StaticFileCacheParam = fmt.Sprintf("max-age=%v", age)
}


func (s *SimpleWebServer) attemptToServeStaticFile(w http.ResponseWriter, r *http.Request) {

	filePath := path.Join(s.StaticDir, r.URL.Path)

	if s.fileExists(filePath){

		// fmt.Printf("Serving static file: %v\n", filePath);
		fmt.Printf("Attempting to serve static file: %v\n", r.URL.Path)


		w.Header().Set("Cache-Control", s.StaticFileCacheParam)

		http.ServeFile(w, r, filePath)

	} else {

		fmt.Printf("File Not Found: %v\n", r.URL.Path)
		//http.NotFound(w, r)
		s.NotFoundHandler(w, r)

	}
	
}


func (s *SimpleWebServer) fileExists(filePath string) bool {

	info, err := os.Stat(filePath)
	if err != nil {
		return false;
	}

	return !info.IsDir()
}


func DefaultNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("404 Not Found"))
}

func DefaultInternalError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(500)
	w.Write([]byte("500 Internal Server Error\n"))
	w.Write([]byte(err.Error()))
}




type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type RegexpHandler struct {
	routes []*route
	Server *SimpleWebServer
	
}

func NewRegexpHandler() *RegexpHandler {
	return &RegexpHandler{make([]*route, 0, 4), nil, 86400}	//default static file cache life of 1 day
}


func (h *RegexpHandler) Handle(pattern string, handler http.Handler){

	rpat, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}

	h.routes = append(h.routes, &route{rpat, handler})
}

func (h *RegexpHandler) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request) ){

	rpat, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}


	h.routes = append(h.routes, &route{rpat, http.HandlerFunc(handler)})
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    for _, route := range h.routes {
        if route.pattern.MatchString(r.URL.Path) {
            route.handler.ServeHTTP(w, r)
            return
        }
    }

    //attempt to serve a static file for the request
    if h.Server != nil {
	    h.Server.attemptToServeStaticFile(w, r);
    } else {
    	w.WriteHeader(500)
    	w.Write([]byte("RegexpHandler: Error: No server object set"))
    }

}





