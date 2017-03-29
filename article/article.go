package article

import (
	"bytes"
	"regexp"
	"strings"
)

var (
	reDate   *regexp.Regexp
	reTitle  *regexp.Regexp
	reTag    *regexp.Regexp
	reStatus *regexp.Regexp
	reURL    *regexp.Regexp
)

type Article struct {
	Date   string
	Title  string
	Tags   []string
	Status string
	URL    string
	Body   []byte
}

func init() {
	reDate = regexp.MustCompile(`^Date: (.+)$`)
	reTitle = regexp.MustCompile(`^Title: (.+)$`)
	reTag = regexp.MustCompile(`^Tags: (.+)$`)
	reStatus = regexp.MustCompile(`^Status: (.+)$`)
	reURL = regexp.MustCompile(`^URL: (.+)$`)
}
func New(input []byte) *Article {
	content := bytes.SplitN(input, []byte("\n\n"), 2)

	result := &Article{Body: content[1]}
	prefixs := bytes.Split(content[0], []byte("\n"))
	for _, prefix := range prefixs {
		if reDate.Match(prefix) {
			date := reDate.FindSubmatch(prefix)
			result.Date = string(date[1])
		}
		if reTitle.Match(prefix) {
			title := reTitle.FindSubmatch(prefix)
			result.Title = string(title[1])
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
			result.URL = string(urls[1])
		}
	}

	return result
}
