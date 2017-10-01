package main

import (
	"github.com/jcherianucla/bloggo/controllers"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	PORT = ":8000"
)

func main() {
	http.HandleFunc("/", controllers.GetPosts)
	http.ListenAndServe(PORT, nil)
}
