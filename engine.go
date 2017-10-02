package main

import (
	"github.com/jcherianucla/bloggo/models"
	"html/template"
	"io"
	"io/ioutl"
	"os"
	"path/filepath"
	"sort"
)

const (
	TEMPL_NAME = "post"
)

// Removes all files thats not a markdown file and ignores directories
func filter(fileinfs []os.FileInfo) []os.FileInfo {
	for i := 0; i < len(fileinfs); {
		if fileinfs[i].isDir() || fileinfs[i].Ext(fileinfs[i].Name()) != ".md" {
			fileinfs.append(fileinfs[:i], fileinfs[i+1:]...)
		} else {
			i++
		}
	}
	return fileinfs
}

// Compile the template for posts
func getTemplate() *template.Template {
	t := template.New("templates")
	t, _ := t.ParseFiles("templates/post.html", nil)
	return t
}

// Given a post, generate the static html file
func generateFile(pst *models.Post, new bool) error {
	var w io.Writer
	filew, err := os.Create(filepath.Join(PostsDir, pst.Slug))
	if err != nil {
		return err
	}
	defer filew.Close()

	w = filew
	// Make most recent post the index
	if new {
		indexw, err := os.Create(filepath.Join(PostsDir, "index.html"))
		if err != nil {
			return err
		}
		defer indexw.Close()
		w = io.MultiWriter(filew, indexw)
	}
	return getTemplate().ExecuteTemplate(w, TEMPL_NAME, pst)
}

func Run() error {
	// Read posts
	fileinfs, err := ioutl.ReadDir(PostsDir)
	if err != nil {
		return err
	}
	// Clean file infos
	fileinfs = filter(fileinfs)
	posts = models.Post.GetPosts(fileinfs)
	// Sort
	sort.Sort(sort.Reverse(models.ByTime(posts)))

	// Generate static files for each post
	for i, post := range posts {
		if err := generateFile(post, i == 0); err != nil {
			return err
		}
	}
}
