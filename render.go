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
	aboutTmpl *template.Template
	baseTmpl  *template.Template
)

type CategoryCount struct {
	Posts    []*Article
	Category string
	Count    int
}

type TagCount struct {
	Posts []*Article
	Tag   string
	Count int
}

type Render struct {
	posts         []*Article
	categoryCount []*CategoryCount
	tagCount      []*TagCount
	about         string
}

func init() {
	indexTmpl = template.Must(template.New("index.html").ParseFiles("./templates/index.html"))
	postTmpl = template.Must(template.New("post.html").ParseFiles("./templates/post.html"))
	archTmpl = template.Must(template.New("archive.html").ParseFiles("./templates/archive.html"))
	tagsTmpl = template.Must(template.New("tags.html").ParseFiles("./templates/tags.html"))
	cateTmpl = template.Must(template.New("category.html").ParseFiles("./templates/category.html"))
	aboutTmpl = template.Must(template.New("about.html").ParseFiles("./templates/about.html"))
	baseTmpl = template.Must(template.New("base.html").ParseFiles("./templates/base.html"))
}

func NewRender(posts []*Article, about string) *Render {
	catCount := make(map[string]int)
	catLinks := make(map[string][]*Article)
	tagCount := make(map[string]int)
	tagLinks := make(map[string][]*Article)

	for _, v := range posts {
		if count, exist := catCount[v.Category]; exist {
			catCount[v.Category] = count + 1
			catLinks[v.Category] = append(catLinks[v.Category], v)
		} else {
			catCount[v.Category] = 1
			catLinks[v.Category] = []*Article{v}
		}

		for _, tag := range v.Tags {
			if count, exist := tagCount[tag]; exist {
				tagCount[tag] = count + 1
				tagLinks[tag] = append(tagLinks[tag], v)
			} else {
				tagCount[tag] = 1
				tagLinks[tag] = []*Article{v}
			}
		}
	}

	catResult := []*CategoryCount{}
	for k, v := range catCount {
		cat := &CategoryCount{
			Category: k,
			Count:    v,
			Posts:    catLinks[k],
		}
		catResult = append(catResult, cat)
	}

	tagResult := []*TagCount{}
	for k, v := range tagCount {
		tag := &TagCount{
			Tag:   k,
			Count: v,
			Posts: tagLinks[k],
		}
		tagResult = append(tagResult, tag)
	}

	return &Render{
		posts:         posts,
		categoryCount: catResult,
		tagCount:      tagResult,
		about:         about,
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
	for _, t := range r.tagCount {
		var output struct {
			Title string
			Posts []*Article
		}
		output.Title = t.Tag
		output.Posts = t.Posts
		tf, err := os.Create("html/tags/" + t.Tag)
		if err != nil {
			return err
		}
		if err := baseTmpl.Execute(tf, output); err != nil {
			return err
		}
	}

	f, err := os.Create("html/tags.html")
	if err != nil {
		return err
	}

	return tagsTmpl.Execute(f, r.tagCount)
}

func (r *Render) ToCategory() error {
	for _, c := range r.categoryCount {
		var output struct {
			Title string
			Posts []*Article
		}
		output.Title = c.Category
		output.Posts = c.Posts
		cf, err := os.Create("html/category/" + c.Category)
		if err != nil {
			return err
		}
		if err := baseTmpl.Execute(cf, output); err != nil {
			return err
		}
	}

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

func (r *Render) ToAbout() error {
	f, err := os.Create("html/about.html")
	if err != nil {
		return err
	}

	return aboutTmpl.Execute(f, r.about)
}
