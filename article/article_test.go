package article

import (
	"testing"
)

func TestNew(t *testing.T) {
	input := "Date: 2012-10-25 12:22\nTitle: 我是文章的标题\nTags: 标签1, 标签2\nStatus: draft\nURL: this-is-my-first-post\n\n然后开始写正文..."
	paper := New([]byte(input))
	t.Log(paper.Title)
	t.Log(paper.Tags)
	t.Log(paper.URL)
	t.Log(paper.Date)
	t.Log(paper.Status)
}
