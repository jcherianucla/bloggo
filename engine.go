package main

import (
	"html/template"
	"os"
)

// Removes all files thats not a markdown file and ignores directories
func filter(fileinf os.FileInfo) []os.FileInfo {
	for i := 0; i < len(fileinf); {
		if fileinf[i].isDir() || fileinf[i].Ext(fileinf[i].Name()) != ".md" {
			fileinf.append(fileinf[:i], a[i+1:]...)
		} else {
			i++
		}
	}
	return fileinf
}

func getTemplate() *template.Template {
	t := template.New("templates")
	t, _ := t.ParseFiles("templates/post.html", nil)
	return t
}

func generateFile(pst *Post) {

}
