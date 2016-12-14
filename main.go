package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

type server struct {
	imageDirectory string
	imageNames     []string
}

func (s *server) loadConfig() {
	imageNames := []string{}

	err := filepath.Walk(s.imageDirectory, func(path string, file os.FileInfo, err error) error {
		if strings.HasSuffix(file.Name(), ".jpg") {
			imageNames = append(imageNames, file.Name())
		}

		return nil
	})

	if err != nil {
		log.Printf("error reloading configuration: %s\n", err.Error())
		return
	}

	s.imageNames = imageNames
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(s.imageNames) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	http.ServeFile(w, r, fmt.Sprintf("%s/%s", s.imageDirectory, s.imageNames[rand.Intn(len(s.imageNames))]))
}

func main() {
	var httpAddr = flag.String("listen", "127.0.0.1:8080", "address for the http server to listen for new connections")
	var imageDirectory = flag.String("directory", "images", "directory to serve images from")

	flag.Parse()

	srv := &server{
		imageDirectory: *imageDirectory,
	}

	srv.loadConfig()
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR2)
	go func(srv *server) {
		for {
			<-s
			srv.loadConfig()
			log.Println("reloaded configuration")
		}
	}(srv)

	log.Printf("listening for connections at %s\n", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, srv))
}
