package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cvley/cvblog"
)

var (
	dir string
)

func init() {
	flag.StringVar(&dir, "dir", "", "markdown file directory")
}

func main() {
	flag.Parse()
	if dir == "" {
		flag.PrintDefaults()
		return
	}

	files, err := filepath.Glob("*.md")
	if err != nil {
		fmt.Println(err)
		return
	}

	posts := []*cvblog.Article{}
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Println(err)
			continue
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			continue
		}
		post := cvblog.NewArticle(b)
		posts = append(posts, post)
	}

	render := cvblog.NewRender(posts, "just about")
	render.SetOutputDir("html")

	render.ToIndex()
	render.ToPosts()
	render.ToAbout()
	render.ToArchive()
	render.ToCategory()
	render.ToTags()
}
