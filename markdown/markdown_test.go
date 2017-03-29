package markdown

import (
	"testing"

	"bytes"
	"io/ioutil"
	"os"
)

func TestParseHeader(t *testing.T) {
	input := [][]byte{
		[]byte("#### test!"),
		[]byte("### 中文"),
		[]byte("invalid"),
	}

	output := [][]byte{
		[]byte("\n<h4> test! </h4>\n"),
		[]byte("\n<h3> 中文 </h3>\n"),
		[]byte("invalid"),
	}

	for i, v := range input {
		result := parseHeader(v)
		if !bytes.Equal(result, output[i]) {
			t.Fatalf("ParseHeader fail, [%s] vs [%s]", string(result), string(output[i]))
		}
	}
}

func TestParseImage(t *testing.T) {
	input := [][]byte{
		[]byte("![xxx](http://www.hackcv.com/test.jpg)"),
		[]byte("xxx.jpg"),
	}

	output := [][]byte{
		[]byte("\n<img src=\"http://www.hackcv.com/test.jpg\" alt=\"xxx\">\n"),
		[]byte("xxx.jpg"),
	}

	for i, v := range input {
		result := parseImage(v)
		if !bytes.Equal(result, output[i]) {
			t.Fatalf("ParseImage fail, [%s] vs [%s]", string(result), string(output[i]))
		}
	}
}

func TestParseList(t *testing.T) {
	input := [][]byte{
		[]byte("- block quote"),
		[]byte("* valid"),
		[]byte("invalid"),
	}

	output := [][]byte{
		[]byte("\n<ul>\n<li>block quote</li>\n</ul>\n"),
		[]byte("\n<ul>\n<li>valid</li>\n</ul>\n"),
		[]byte("invalid"),
	}
	for i, v := range input {
		result := parseList(v)
		if !bytes.Equal(result, output[i]) {
			t.Fatalf("ParseQuote fail, [%s] vs [%s]", string(result), string(output[i]))
		}
	}
}

func TestParseQuote(t *testing.T) {
	input := [][]byte{
		[]byte("> block quote"),
		[]byte("invalid"),
	}

	output := [][]byte{
		[]byte("\n<blockquote>block quote</blockquote>\n"),
		[]byte("invalid"),
	}
	for i, v := range input {
		result := parseQuote(v)
		if !bytes.Equal(result, output[i]) {
			t.Fatalf("ParseQuote fail, [%s] vs [%s]", string(result), string(output[i]))
		}
	}
}

func TestParseInlineEmphasis(t *testing.T) {
	input := [][]byte{
		[]byte("just **test** **test** test"),
		[]byte("just __test__ __test__ test"),
		[]byte("****normal test"),
		[]byte("normal test"),
	}

	output := [][]byte{
		[]byte("just <strong>test</strong> <strong>test</strong> test"),
		[]byte("just <strong>test</strong> <strong>test</strong> test"),
		[]byte("normal test"),
		[]byte("normal test"),
	}

	for i, v := range input {
		r := parseInlineEmphasis(v)
		if !bytes.Equal(r, output[i]) {
			t.Fatalf("ParseInlineEmphasis fail, [%s] vs [%s]", string(r), string(output[i]))
		}
	}
}

func TestParseInlineItalics(t *testing.T) {
	input := [][]byte{
		[]byte("just *test* *test* test"),
		[]byte("just _test_ _test_ test"),
		[]byte("**normal test"),
		[]byte("normal test"),
	}

	output := [][]byte{
		[]byte("just <em>test</em> <em>test</em> test"),
		[]byte("just <em>test</em> <em>test</em> test"),
		[]byte("normal test"),
		[]byte("normal test"),
	}

	for i, v := range input {
		r := parseInlineItalics(v)
		if !bytes.Equal(r, output[i]) {
			t.Fatalf("ParseInlineItalics fail, [%s] vs [%s]", string(r), string(output[i]))
		}
	}
}

func TestParseInlineStrike(t *testing.T) {
	input := [][]byte{
		[]byte("just ~~test~~ ~~test~~ test"),
		[]byte("~~~~normal test"),
		[]byte("normal test"),
	}

	output := [][]byte{
		[]byte("just <del>test</del> <del>test</del> test"),
		[]byte("normal test"),
		[]byte("normal test"),
	}

	for i, v := range input {
		r := parseInlineStrike(v)
		if !bytes.Equal(r, output[i]) {
			t.Fatalf("ParseInlineStrike fail, [%s] vs [%s]", string(r), string(output[i]))
		}
	}
}

func TestParseInlineCode(t *testing.T) {
	input := [][]byte{
		[]byte("just `test` `test` test"),
		[]byte("````normal test"),
		[]byte("normal test"),
	}

	output := [][]byte{
		[]byte("just <code>test</code> <code>test</code> test"),
		[]byte("normal test"),
		[]byte("normal test"),
	}

	for i, v := range input {
		r := parseInlineCode(v)
		if !bytes.Equal(r, output[i]) {
			t.Fatalf("ParseInlineCode fail, [%s] vs [%s]", string(r), string(output[i]))
		}
	}
}

func TestParseCode(t *testing.T) {
	input := [][]byte{
		[]byte("```test\nblock quote\n```"),
		[]byte("```\nvalid\n```"),
	}

	output := [][]byte{
		[]byte("\n<pre lang=\"test\">\n<code>\nblock quote\n</code>\n</pre>\n"),
		[]byte("\n<pre>\n<code>\nvalid\n</code>\n</pre>\n"),
	}
	for i, v := range input {
		result := parseCode(v)
		if !bytes.Equal(result, output[i]) {
			t.Fatalf("ParseCode fail, [%s] vs [%s]", string(result), string(output[i]))
		}
	}
}

func TestParseInlineLink(t *testing.T) {
	input := [][]byte{
		[]byte("test [link](http://hackcv.com)"),
		[]byte("[test](http://baidu.com) test [link](http://hackcv.com)"),
	}

	output := [][]byte{
		[]byte("test <a href=\"http://hackcv.com\">link</a>"),
		[]byte("<a href=\"http://baidu.com\">test</a> test <a href=\"http://hackcv.com\">link</a>"),
	}
	for i, v := range input {
		result := parseInlineLink(v)
		if !bytes.Equal(result, output[i]) {
			t.Fatalf("ParseCode fail, [%s] vs [%s]", string(result), string(output[i]))
		}
	}
}

func TestBlock(t *testing.T) {
	input := "```shell\n./configure\nmake\nmake install\n```"
	block := NewBlock([]byte(input))
	result := block.Render()
	t.Log(string(result))
}

func TestRender(t *testing.T) {
	input := "README.md"
	f, err := os.Open(input)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	result := Render(b)
	t.Log(string(result))
}
