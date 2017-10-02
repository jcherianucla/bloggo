package main

import (
	"log"
	"os"
	"path/filepath"
)

var (
	PostsDir     string
	TemplatesDir string
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("[ERROR]", err)
	}
	PostsDir = filepath.Join(wd, "posts")
	TemplatesDir = filepath.Join(wd, "templates")
}
