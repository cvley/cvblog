package cvblog

import (
	"testing"
)

var (
	posts []*Article
)

func init() {
	input := []string{
		"Date: 2012-10-25 12:22\nTitle: 我是文章的标题2\nCategory: Test\nTags: 标签1, 标签2\nStatus: draft\nURL: this-is-my-first-post\n\n然后开始写正文...",
	}

	posts = make([]*Article, len(input))
	for i, v := range input {
		a := NewArticle([]byte(v))
		posts[i] = a
	}
}

func TestRender(t *testing.T) {
	render := NewRender(posts)
	if err := render.ToPosts(); err != nil {
		t.Fatal(err)
	}

	if err := render.ToArchive(); err != nil {
		t.Fatal(err)
	}
}
