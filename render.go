package cvblog

import (
	"html/template"
	"os"
)

var (
	indexTmpl *template.Template
	postTmpl  *template.Template
	archTmpl  *template.Template
)

type Render struct {
	posts []*Article
}

func init() {
	indexTmpl = template.Must(template.New("index.html").ParseFiles("./templates/index.html"))
	postTmpl = template.Must(template.New("post.html").ParseFiles("./templates/post.html"))
	archTmpl = template.Must(template.New("archive.html").ParseFiles("./templates/archive.html"))
}

func NewRender(posts []*Article) *Render {
	return &Render{
		posts: posts,
	}
}

func (r *Render) ToPosts() error {
	for _, v := range r.posts {
		f, err := os.Create("html/" + v.URL)
		if err != nil {
			return err
		}

		if err := postTmpl.Execute(f, v); err != nil {
			return err
		}
	}

	return nil
}

func (r *Render) ToArchive() error {
	f, err := os.Create("html/archives.html")
	if err != nil {
		return err
	}

	return archTmpl.Execute(f, r.posts)
}
