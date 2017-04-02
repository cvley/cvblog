package article

import (
	"sort"
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
	t.Log(paper.Summary())
}

func TestSortByDate(t *testing.T) {
	input := []string{
		"Date: 2012-10-25 12:22\nTitle: 我是文章的标题2\nCategory: Test\nTags: 标签1, 标签2\nStatus: draft\nURL: this-is-my-first-post\n\n然后开始写正文...",
		"Date: 2012-10-25 13:22\nTitle: 我是文章的标题1\nTags: 标签1, 标签2\nStatus: draft\nURL: this-is-my-first-post\n\n然后开始写正文...",
	}

	articles := []Article{}
	for _, v := range input {
		paper := New([]byte(v))
		articles = append(articles, paper)
	}

	sort.Sort(ArticleSortByTime(articles))
	t.Log(articles[0].Date)
	t.Log(articles[0].Title)
}
