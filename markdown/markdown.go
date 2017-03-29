package markdown

import (
	"bytes"
	"fmt"
	"regexp"

	_ "log"
)

var (
	reHeader *regexp.Regexp
	reImage  *regexp.Regexp
	reTable  *regexp.Regexp
	reCode   *regexp.Regexp
	reQuote  *regexp.Regexp
	reList   *regexp.Regexp
)

var (
	inlineReItalics  *regexp.Regexp
	inlineReEmphasis *regexp.Regexp
	inlineReStrike   *regexp.Regexp
	inlineReCode     *regexp.Regexp
)

var (
	lineTrail = []byte("\r\n")
)

func init() {
	reHeader = regexp.MustCompile(`^(#{1,6})\s*(\p{Han}+|[[:ascii:]]+)\s*#*$`)
	reImage = regexp.MustCompile(`^!\[(\w+)\]\((.*)\)$`)
	reQuote = regexp.MustCompile(`^>\s(.*)`)

	inlineReEmphasis = regexp.MustCompile(`\*{2}|\_{2}`)
	inlineReItalics = regexp.MustCompile(`\*|\_`)
	inlineReStrike = regexp.MustCompile(`\~{2}`)
	inlineReCode = regexp.MustCompile("`")
}

type Parser interface {
	Parse([]byte) []byte
	Match([]byte) bool
}

// BlockData contains the block byte data and type
type BlockData struct {
	Data []byte
}

func ParseInlineCode(input []byte) []byte {
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

func ParseInlineStrike(input []byte) []byte {
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

func ParseInlineItalics(input []byte) []byte {
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

func ParseInlineEmphasis(input []byte) []byte {
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

func ParseQuate(input []byte) []byte {
	if !reQuote.Match(input) {
		return input
	}

	var buffer bytes.Buffer
	buffer.Write([]byte("<blockquote>"))

	lines := bytes.Split(input, lineTrail)
	for _, line := range lines {
		ret := reQuote.FindSubmatch(line)
		buffer.Write(ret[1])
	}
	buffer.Write([]byte("</blockquote>"))

	return buffer.Bytes()
}

func ParseHeader(input []byte) []byte {
	if !reHeader.Match(input) {
		return input
	}

	ret := reHeader.FindSubmatch(input)
	num := len(ret[1])
	header := fmt.Sprintf("<h%d> %s </h%d>", num, string(ret[2]), num)
	return []byte(header)
}

func ParseImage(input []byte) []byte {
	if !reImage.Match(input) {
		return input
	}

	ret := reImage.FindSubmatch(input)
	alt := ret[1]
	src := ret[2]
	result := fmt.Sprintf("<img src=\"%s\" alt=\"%s\">", src, alt)
	return []byte(result)
}
