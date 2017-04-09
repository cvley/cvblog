package markdown

import (
	"bytes"
	"fmt"
	"regexp"
)

var (
	reHeader *regexp.Regexp
	reImage  *regexp.Regexp
	reCode   *regexp.Regexp
	reQuote  *regexp.Regexp
	reList   *regexp.Regexp
)

var (
	inlineReItalics  *regexp.Regexp
	inlineReEmphasis *regexp.Regexp
	inlineReStrike   *regexp.Regexp
	inlineReCode     *regexp.Regexp
	inlineReLink     *regexp.Regexp
)

var (
	lineTrail  = []byte("\n")
	blockTrail = []byte("\n\n")
	codeSign   = []byte("```")
)

type BlockType int

const (
	BlockTypeParagraph BlockType = iota
	BlockTypeHeader
	BlockTypeImage
	BlockTypeList
	BlockTypeCode
	BlockTypeQuote
)

func (tp BlockType) String() string {
	switch tp {
	case BlockTypeHeader:
		return "Header Block"

	case BlockTypeImage:
		return "Image Block"

	case BlockTypeList:
		return "List Block"

	case BlockTypeCode:
		return "Code Block"

	case BlockTypeQuote:
		return "Quote Block"

	default:
		return "Paragraph Block"
	}
}

type Block struct {
	data []byte
	tp   BlockType
}

func init() {
	reHeader = regexp.MustCompile(`^(#{1,6})\s*(\p{Han}+|[[:ascii:]]+)\s*#*$`)
	reImage = regexp.MustCompile(`^!\[(.*)\]\((.+)\)$`)
	reQuote = regexp.MustCompile(`^>\s(.*)$`)
	reList = regexp.MustCompile(`^[*|-]\s(.*)$`)
	reCode = regexp.MustCompile("^`{3}(\\w+)$")

	inlineReEmphasis = regexp.MustCompile(`\*{2}|\_{2}`)
	inlineReItalics = regexp.MustCompile(`\*|\_`)
	inlineReStrike = regexp.MustCompile(`\~{2}`)
	inlineReCode = regexp.MustCompile("`")
	inlineReLink = regexp.MustCompile(`\[([^\[]+)\]\(([^\]]+)\)`)
}

func Render(input []byte) []byte {
	blocks := bytes.Split(input, blockTrail)

	buffer := bytes.Buffer{}
	for _, data := range blocks {
		block := NewBlock(data)
		buffer.Write(block.Render())
	}

	return buffer.Bytes()
}

func NewBlock(input []byte) *Block {
	tp := getBlockType(input)
	return &Block{
		data: input,
		tp:   tp,
	}
}

func (block *Block) Render() []byte {
	switch block.tp {
	case BlockTypeCode:
		return parseCode(block.data)

	case BlockTypeHeader:
		return parseHeader(block.data)

	case BlockTypeImage:
		return parseImage(block.data)

	case BlockTypeList:
		data := block.renderInline()
		return parseList(data)

	case BlockTypeQuote:
		data := block.renderInline()
		return parseQuote(data)
	}

	buffer := bytes.Buffer{}
	buffer.WriteString("\n<p>")
	buffer.Write(block.data)
	buffer.WriteString("</p>\n")

	return buffer.Bytes()
}

func (block *Block) renderInline() []byte {
	result := parseInlineCode(block.data)
	result = parseInlineEmphasis(result)
	result = parseInlineItalics(result)
	result = parseInlineLink(result)
	return parseInlineStrike(result)
}

func getBlockType(input []byte) BlockType {
	if bytes.HasPrefix(input, codeSign) && bytes.HasSuffix(input, codeSign) {
		return BlockTypeCode
	}

	if reImage.Match(input) {
		return BlockTypeImage
	}

	if reHeader.Match(input) {
		return BlockTypeHeader
	}

	if isList(input) {
		return BlockTypeList
	}

	if isQuote(input) {
		return BlockTypeQuote
	}

	return BlockTypeParagraph
}

func isList(input []byte) bool {
	lists := bytes.Split(input, lineTrail)
	for _, list := range lists {
		if !reList.Match(list) {
			return false
		}
	}

	return true
}

func isQuote(input []byte) bool {
	quotes := bytes.Split(input, lineTrail)
	for _, quote := range quotes {
		if !reQuote.Match(quote) {
			return false
		}
	}

	return true
}

func parseCode(input []byte) []byte {
	contents := bytes.Split(input, lineTrail)

	pre := "\n<pre>\n<code>\n"
	result := reCode.FindSubmatch(contents[0])
	if len(result) == 2 {
		pre = fmt.Sprintf("\n<pre lang=\"%s\">\n<code>\n", result[1])
	}

	var buffer bytes.Buffer
	buffer.WriteString(pre)
	codes := contents[1 : len(contents)-1]
	buffer.Write(bytes.Join(codes, lineTrail))
	buffer.WriteString("\n</code>\n</pre>\n")

	return buffer.Bytes()
}

func parseQuote(input []byte) []byte {
	var buffer bytes.Buffer
	buffer.Write([]byte("\n<blockquote>"))

	lines := bytes.Split(input, lineTrail)
	for _, line := range lines {
		ret := reQuote.FindSubmatch(line)
		if ret == nil {
			return input
		}
		buffer.Write(ret[1])
	}
	buffer.Write([]byte("</blockquote>\n"))

	return buffer.Bytes()
}

func parseHeader(input []byte) []byte {
	if !reHeader.Match(input) {
		return input
	}

	ret := reHeader.FindSubmatch(input)
	num := len(ret[1])
	header := fmt.Sprintf("\n<h%d> %s </h%d>\n", num, string(ret[2]), num)
	return []byte(header)
}

func parseList(input []byte) []byte {
	lists := bytes.Split(input, lineTrail)

	var buffer bytes.Buffer
	buffer.WriteString("\n<ul>\n")
	for _, list := range lists {
		result := reList.FindAllSubmatch(list, -1)
		if result == nil {
			return input
		}

		for _, v := range result {
			buffer.WriteString("<li>")
			buffer.Write(v[1])
			buffer.WriteString("</li>\n")
		}
	}

	buffer.WriteString("</ul>\n")

	return buffer.Bytes()
}

func parseImage(input []byte) []byte {
	ret := reImage.FindSubmatch(input)
	if ret == nil {
		return input
	}

	alt := ret[1]
	src := ret[2]
	result := fmt.Sprintf("\n<img src=\"%s\" alt=\"%s\">\n", src, alt)
	return []byte(result)
}

func parseInlineCode(input []byte) []byte {
	result := inlineReCode.FindAllSubmatchIndex(input, -1)
	if result == nil || len(result)%2 == 1 {
		return input
	}

	var buffer bytes.Buffer
	var index int
	for i := 0; i < len(result); i = i + 2 {
		first := result[i]
		second := result[i+1]
		if first[1] == second[0] {
			index = second[1]
			continue
		}

		buffer.Write(input[index:first[0]])
		buffer.Write([]byte("<code>"))
		buffer.Write(input[first[1]:second[0]])
		buffer.Write([]byte("</code>"))
		index = second[1]
	}
	buffer.Write(input[index:])

	return buffer.Bytes()
}

func parseInlineStrike(input []byte) []byte {
	result := inlineReStrike.FindAllSubmatchIndex(input, -1)
	if result == nil || len(result)%2 == 1 {
		return input
	}

	var buffer bytes.Buffer
	var index int
	for i := 0; i < len(result); i = i + 2 {
		first := result[i]
		second := result[i+1]
		if first[1] == second[0] {
			index = second[1]
			continue
		}

		buffer.Write(input[index:first[0]])
		buffer.Write([]byte("<del>"))
		buffer.Write(input[first[1]:second[0]])
		buffer.Write([]byte("</del>"))
		index = second[1]
	}
	buffer.Write(input[index:])

	return buffer.Bytes()
}

func parseInlineItalics(input []byte) []byte {
	result := inlineReItalics.FindAllSubmatchIndex(input, -1)
	if result == nil || len(result)%2 == 1 {
		return input
	}

	var buffer bytes.Buffer
	var index int
	for i := 0; i < len(result); i = i + 2 {
		first := result[i]
		second := result[i+1]
		if first[1] == second[0] {
			index = second[1]
			continue
		}

		buffer.Write(input[index:first[0]])
		buffer.Write([]byte("<em>"))
		buffer.Write(input[first[1]:second[0]])
		buffer.Write([]byte("</em>"))
		index = second[1]
	}
	buffer.Write(input[index:])

	return buffer.Bytes()
}

func parseInlineEmphasis(input []byte) []byte {
	result := inlineReEmphasis.FindAllSubmatchIndex(input, -1)
	if result == nil || len(result)%2 == 1 {
		return input
	}

	var buffer bytes.Buffer
	var index int
	for i := 0; i < len(result); i = i + 2 {
		first := result[i]
		second := result[i+1]
		if first[1] == second[0] {
			index = second[1]
			continue
		}

		buffer.Write(input[index:first[0]])
		buffer.Write([]byte("<strong>"))
		buffer.Write(input[first[1]:second[0]])
		buffer.Write([]byte("</strong>"))
		index = second[1]
	}
	buffer.Write(input[index:])

	return buffer.Bytes()
}

func parseInlineLink(input []byte) []byte {
	indexes := inlineReLink.FindAllSubmatchIndex(input, -1)
	if indexes == nil {
		return input
	}

	var buffer bytes.Buffer
	var start int
	for _, index := range indexes {
		if start < index[0] {
			buffer.Write(input[start:index[0]])
		}

		var b []byte
		b = inlineReLink.Expand(b, []byte("<a href=\"$2\">$1</a>"), input, index)
		buffer.Write(b)
		start = index[1]
	}

	return buffer.Bytes()
}
