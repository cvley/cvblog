package markdown

import (
	"testing"

	"bytes"
)

func TestParseHeader(t *testing.T) {
	input := [][]byte{
		[]byte("#### test!"),
		[]byte("### 中文"),
		[]byte("invalid"),
	}

	output := [][]byte{
		[]byte("<h4> test! </h4>"),
		[]byte("<h3> 中文 </h3>"),
		[]byte("invalid"),
	}

	for i, v := range input {
		result := ParseHeader(v)
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
		[]byte("<img src=\"http://www.hackcv.com/test.jpg\" alt=\"xxx\">"),
		[]byte("xxx.jpg"),
	}

	for i, v := range input {
		result := ParseImage(v)
		if !bytes.Equal(result, output[i]) {
			t.Fatalf("ParseImage fail, [%s] vs [%s]", string(result), string(output[i]))
		}
	}
}

func TestParseQuate(t *testing.T) {
	input := [][]byte{
		[]byte("> block quote"),
		[]byte("invalid"),
	}

	output := [][]byte{
		[]byte("<blockquote>block quote</blockquote>"),
		[]byte("invalid"),
	}
	for i, v := range input {
		result := ParseQuate(v)
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
		r := ParseInlineEmphasis(v)
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
		r := ParseInlineItalics(v)
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
		r := ParseInlineStrike(v)
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
		r := ParseInlineCode(v)
		if !bytes.Equal(r, output[i]) {
			t.Fatalf("ParseInlineCode fail, [%s] vs [%s]", string(r), string(output[i]))
		}
	}
}
