package main

import (
	"errors"
	"net/http"
	"strings"
)

func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

var (
	errorNotAllowed = errors.New("not allowed")
)

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	filesystem, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := filesystem.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, errorNotAllowed
	}

	return filesystem, nil
}
