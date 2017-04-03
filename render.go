package cvblog

import (
	"html/template"
	"os"
)

var (
	indexTmpl *template.Template
	postTmpl  *template.Template
	archTmpl  *template.Template
	tagsTmpl  *template.Template
	cateTmpl  *template.Template
)

type CategoryCount struct {
	Category string
	Count    int
}

type TagCount struct {
	Tag   string
	Count int
}

type Render struct {
	posts         []*Article
	categoryCount []*CategoryCount
	tagCount      []*TagCount
}

func init() {
	indexTmpl = template.Must(template.New("index.html").ParseFiles("./templates/index.html"))
	postTmpl = template.Must(template.New("post.html").ParseFiles("./templates/post.html"))
	archTmpl = template.Must(template.New("archive.html").ParseFiles("./templates/archive.html"))
	tagsTmpl = template.Must(template.New("tags.html").ParseFiles("./templates/tags.html"))
	cateTmpl = template.Must(template.New("category.html").ParseFiles("./templates/category.html"))
}

func NewRender(posts []*Article) *Render {
	catCount := make(map[string]int)
	tagCount := make(map[string]int)

	for _, v := range posts {
		if count, exist := catCount[v.Category]; exist {
			catCount[v.Category] = count + 1
		} else {
			catCount[v.Category] = 1
		}

		for _, tag := range v.Tags {
			if count, exist := tagCount[tag]; exist {
				tagCount[tag] = count + 1
			} else {
				tagCount[tag] = 1
			}
		}
	}

	catResult := []*CategoryCount{}
	for k, v := range catCount {
		cat := &CategoryCount{
			Category: k,
			Count:    v,
		}
		catResult = append(catResult, cat)
	}

	tagResult := []*TagCount{}
	for k, v := range tagCount {
		tag := &TagCount{
			Tag:   k,
			Count: v,
		}
		tagResult = append(tagResult, tag)
	}

	return &Render{
		posts:         posts,
		categoryCount: catResult,
		tagCount:      tagResult,
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

func (r *Render) ToTags() error {
	f, err := os.Create("html/tags.html")
	if err != nil {
		return err
	}

	return tagsTmpl.Execute(f, r.tagCount)
}

func (r *Render) ToCategory() error {
	f, err := os.Create("html/category.html")
	if err != nil {
		return err
	}

	return cateTmpl.Execute(f, r.categoryCount)
}

func (r *Render) ToIndex() error {
	f, err := os.Create("html/index.html")
	if err != nil {
		return err
	}
	length := len(r.posts)
	if length > 5 {
		length = 5
	}

	return indexTmpl.Execute(f, r.posts[:length])
}
