package web

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"text/template"

	"github.com/go-chi/chi"
)

type Server struct {
	templates    map[string]*template.Template
	lock         sync.RWMutex
}

func StartHttpServer(httpHost, httpPort string) {
	server := &Server{
		templates:    map[string]*template.Template{},
	}

	// load templates
	server.loadTemplates()

	router := chi.NewRouter()
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "web/public")
	FileServer(router, "/static", http.Dir(filesDir))
	server.registerHandlers(router)

	address := net.JoinHostPort(httpHost, httpPort)

	fmt.Printf("starting http server on %s\n", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		fmt.Println("Error starting web server")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (s *Server) loadTemplates() {
	layout := "web/views/layout.html"
	tpls := map[string]string{
		"balance.html": "web/views/balance.html",
	}

	for i, v := range tpls {
		tpl, err := template.New(i).ParseFiles(v, layout)
		if err != nil {
			log.Fatalf("error loading templates: %s", err.Error())
		}

		s.lock.Lock()
		s.templates[i] = tpl
		s.lock.Unlock()
	}
}

func (s *Server) render(tplName string, data map[string]interface{}, res http.ResponseWriter) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if tpl, ok := s.templates[tplName]; ok {
		err := tpl.Execute(res, data)
		if err != nil {
			log.Fatalf("error executing template: %s", err.Error())
		}
		return
	}

	log.Fatalf("template %s is not registered", tplName)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func (s *Server) registerHandlers(r *chi.Mux) {
	r.Get("/", s.GetBalance)
	r.Get("/send", s.GetSend)
	r.Post("/send", s.PostSend)
	r.Get("/receive", s.GetReceive)
}
