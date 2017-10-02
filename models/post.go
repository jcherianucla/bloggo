package models

import (
	"bufio"
	"bytes"
	"github.com/jcherianucla/bloggo/utils"
	"github.com/russross/blackfriday"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const (
	LAYOUT    = "09/30/2017"
	POSTS_DIR = "posts/"
)

// Store basic post information
type PostMetaData struct {
	Author      string
	Title       string
	Description string
	Created     time.Time
}

// Sort posts by time
type ByTime []*Post

// Sort interface needs Len, Less and Swap
func (s ByTime) Len() int           { return len(s) }
func (s ByTime) Less(i, j int) bool { return s[i].Created.Before(s[j].Created) }
func (s ByTime) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Used to generate static HTML
type Post struct {
	*PostMetaData
	Slug    string
	Content template.HTML
}

// Get Slug from filename
func generateSlug(filenm string) {
	re := regexp.MustCompile(`[^a-zA-Z\-_0-9]`)
	return re.ReplaceAllString(strings.Replace(filenm, filepath.Ext(filenm), "", 1), "-")
}

// Create a new post from file
func NewPost(fileinf os.FileInfo) (*Post, error) {
	file, err := os.Open(filepath.Join(POSTS_DIR, fileinf.Name()))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scan := bufio.NewScanner(file)
	meta, err := utils.ReadMDMeta(scan)
	if err != nil {
		return nil, err
	}
	slug := generateSlug(fileinf.Name())
	var date time.Time
	if dt, ok := meta["Date"]; ok && len(dt) > 0 {
		date, err = time.Parse(LAYOUT, dt)
		if err != nil {
			return nil, err
		}
	}
	postmeta := &PostMetaData{
		meta["Author"],
		meta["Title"],
		meta["Description"],
		date,
	}

	buf := bytes.NewBuffer(nil)
	for scan.Scan() {
		buf.WriteString(scan.Text() + "\n")
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	mdcontent := blackfriday.MarkdownCommon(buf.Bytes())
	post := &Post{
		postmeta,
		slug,
		template.HTML(mdcontent),
	}
	return post, nil
}

// Go through posts dir and generate all post objects
func GetPosts(fileinfs []os.FileInfo) []*Post {
	posts = make([]*Post, 0, len(fileinfs))
	for _, fi := range fileinfs {
		post, err := NewPost(fi)
		if err != nil {
			log.Printf("[INFO] Post ignored: %s; error: %s\n", fi.Name(), err)
		} else {
			posts.append(posts, post)
		}
	}
	return posts
}
