package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jordan-wright/unindexed"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	router := mux.NewRouter()
	//Маршруты
	router.HandleFunc(`/hello_world/{name}`, HandleHelloWorld)
	//Маршруты для .html
	router.HandleFunc(`/{filename}.html`, HandleHTML)
	//Статика
	staticPath, _ := filepath.Abs("www/static")
	fs := http.FileServer(unindexed.Dir(staticPath))
	router.PathPrefix(`/`).Handler(http.StripPrefix(`/`, fs))

	router.Use(loggingHandlingPageMiddleware)
	StartServer(router)
}

func StartServer(router *mux.Router) {
	server := &http.Server{
		Addr:         ":7020",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	log.Info().Msg(`Starting server...`)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Err(err).Msg(`Start server is fault!`)
			return
		}
	}()
	log.Info().Msg(`Server is starting!`)
	// Ожидание сигнала
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	// Завершение
	log.Info().Msg(`Server stopping...`)
	err := server.Close()
	if err != nil {
		log.Err(err).Msg(`Server stopped fault`)
	} else {
		log.Info().Msg(`Server stopped!`)
	}
}

func loggingHandlingPageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Trace().Msgf("[%s] %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func HandleHelloWorld(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "Hello, %s", vars[`name`])
	if err != nil {
		log.Err(err).Msg(`Fault handle route`)
	}
}
func HandleHTML(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tplPath, _ := filepath.Abs(fmt.Sprintf("www/tpl/%s.html", vars[`filename`]))

	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		log.Err(err).Msg(`File not found`)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = tpl.Execute(w, `uy`)
	if err != nil {
		log.Err(err).Msg(`Fault handle route`)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
