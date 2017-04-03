package cvblog

import (
	"bytes"
	"html/template"
	"regexp"
	"strings"
	"time"

	"github.com/cvley/cvblog/markdown"
)

const (
	defaultCategory = "心得体会"
)

var (
	reDate     *regexp.Regexp
	reTitle    *regexp.Regexp
	reCategory *regexp.Regexp
	reTag      *regexp.Regexp
	reStatus   *regexp.Regexp
	reURL      *regexp.Regexp
)

type Article struct {
	Time     time.Time
	Date     string
	Title    template.HTML
	Category string
	Tags     []string
	Status   string
	URL      string
	Body     template.HTML
}

type ArticleSortByTime []*Article

func init() {
	reDate = regexp.MustCompile(`^Date: (.+)$`)
	reTitle = regexp.MustCompile(`^Title: (.+)$`)
	reTag = regexp.MustCompile(`^Tags: (.+)$`)
	reStatus = regexp.MustCompile(`^Status: (.+)$`)
	reURL = regexp.MustCompile(`^URL: (.+)$`)
}

func NewArticle(input []byte) *Article {
	content := bytes.SplitN(input, []byte("\n\n"), 2)

	result := &Article{
		Body:     template.HTML(markdown.Render(content[1])),
		Category: defaultCategory,
	}
	prefixs := bytes.Split(content[0], []byte("\n"))
	for _, prefix := range prefixs {
		if reDate.Match(prefix) {
			date := reDate.FindSubmatch(prefix)
			t, err := time.Parse("2006-01-02 15:04", string(date[1]))
			if err != nil {
				t = time.Now()
			}
			result.Time = t
			result.Date = t.Format("2006-01-02 15:04")
		}
		if reTitle.Match(prefix) {
			title := reTitle.FindSubmatch(prefix)
			result.Title = template.HTML(title[1])
		}
		if reTag.Match(prefix) {
			tags := reTag.FindSubmatch(prefix)
			result.Tags = strings.Split(string(tags[1]), ",")
		}
		if reStatus.Match(prefix) {
			status := reStatus.FindSubmatch(prefix)
			result.Status = string(status[1])
		}
		if reURL.Match(prefix) {
			urls := reURL.FindSubmatch(prefix)
			if bytes.HasSuffix(urls[1], []byte(".html")) {
				result.URL = string(urls[1])
			} else {
				result.URL = string(urls[1]) + ".html"
			}
		}
	}

	return result
}

func (a *Article) SetCategory(c string) {
	a.Category = c
}

func (a *Article) Summary() string {
	length := len(a.Body)
	if length > 500 {
		length = 500
	}

	return string(a.Body)[:length]
}

func (byTime ArticleSortByTime) Len() int {
	return len(byTime)
}

func (byTime ArticleSortByTime) Swap(i, j int) {
	byTime[i], byTime[j] = byTime[j], byTime[i]
}

func (byTime ArticleSortByTime) Less(i, j int) bool {
	return byTime[i].Time.Unix() > byTime[j].Time.Unix()
}
