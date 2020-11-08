package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/chetanyakan/realworld-go/controller"
	"github.com/chetanyakan/realworld-go/util"
)

// Context keep context of the running application
type Context struct {
	server *http.Server
	router *mux.Router
}

func main() {
	var (
		port          = flag.String("port", "8080", "web server port")
		templatesPath = flag.String("templates", "./templates/", "templates folder")
		staticPath    = flag.String("static", "./assets/", "static assets folder")
		baseURL       = flag.String("baseurl", os.Getenv("BASE_URL"), "local base url")
	)
	flag.Parse()

	c := &Context{}
	router := controller.NewController(*baseURL, *templatesPath, *staticPath).InitAPI()
	c.router = router

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", *port),
		Handler: c.router,
	}
	c.server = server

	util.Logger.Printf("Http server listening on port: %v", *port)
	util.Logger.Fatalln(c.server.ListenAndServe())
}

func (c *Context) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Logger.Println(
		"New request:",
		"Host", r.Host,
		"RequestURI", r.RequestURI,
		"Method", r.Method,
		"RequestDump", util.DumpRequest(r),
	)
	c.router.ServeHTTP(w, r)
}
