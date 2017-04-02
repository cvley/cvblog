package cvblog

import (
	"testing"
)

func TestParsePost(t *testing.T) {
	input := []string{
		"Date: 2012-10-25 12:22\nTitle: 我是文章的标题2\nCategory: Test\nTags: 标签1, 标签2\nStatus: draft\nURL: this-is-my-first-post\n\n然后开始写正文...",
	}

	for _, v := range input {
		if err := ParsePost([]byte(v)); err != nil {
			t.Fatal(err)
		}
	}
}
