package controller

import (
	"html/template"
	"net/http"
	"path"

	"github.com/chetanyakan/realworld-go/util"
)

func (c *Controller) hello(w http.ResponseWriter, _ *http.Request) {
	lp := path.Join(c.templatesPath, "hello.html")

	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		util.Logger.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = tmpl.ExecuteTemplate(w, "hello.html", map[string]string{})
}
