package blog

import (
	"html/template"
	"os"
	"log"

	"github.com/cvley/cvblog/article"
)

var (
	postTmpl *template.Template
)

func init() {
	postTmpl = template.Must(template.New("post.html").ParseFiles("./templates/post.html"))
}

func ParsePost(input []byte) error {
	post := article.New(input)
	log.Printf("%+v", post)

	f, err := os.Create(post.URL)
	if err != nil {
		return err
	}

	if err := postTmpl.Execute(f, post); err != nil {
		return err
	}

	return nil
}
