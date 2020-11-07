package controller

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"

	"github.com/chetanyakan/realworld-go/util"
)

type Controller struct {
	baseUrl       string
	templatesPath string
	staticPath    string
}

func NewController(baseUrl, templatesPath, staticPath string) *Controller {
	return &Controller{
		baseUrl,
		templatesPath,
		staticPath,
	}
}

// InitAPI initializes the REST API
func (c *Controller) InitAPI() *mux.Router {
	r := mux.NewRouter()
	r.Use(c.withRecovery, c.withLogging)

	// Health-check endpoint
	r.HandleFunc("/ping", c.healthCheck).Methods(http.MethodGet)

	c.handleStaticFiles(r)
	s := r.PathPrefix("/api/v1").Subrouter()

	// Add the custom plugin routes here
	s.HandleFunc("/", c.hello).Methods(http.MethodGet)

	// 404 handler
	r.Handle("{anything:.*}", http.NotFoundHandler())
	return r
}

func (c *Controller) healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]string{"OK"})
}

func (c *Controller) withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.Logger.Println(
			"New request:",
			"Host", r.Host,
			"RequestURI", r.RequestURI,
			"Method", r.Method,
			"RequestDump", util.DumpRequest(r),
		)

		next.ServeHTTP(w, r)
	})
}

func (c *Controller) withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				util.Logger.Println("Recovered from a panic",
					"url", r.URL.String(),
					"error", x,
					"stack", string(debug.Stack()))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// handleStaticFiles handles the static files under the assets directory.
func (c *Controller) handleStaticFiles(r *mux.Router) {
	// This will serve static files from the 'assets' directory under '/static/<filename>'
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(c.staticPath))))
}
