package cvblog

import (
	"html/template"
	"log"
	"os"
)

var (
	indexTmpl *template.Template
	postTmpl  *template.Template
)

func init() {
	indexTmpl = template.Must(template.New("index.html").ParseFiles("./templates/index.html"))
	postTmpl = template.Must(template.New("post.html").ParseFiles("./templates/post.html"))
}

func ParsePost(input []byte) error {
	post := NewArticle(input)
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
