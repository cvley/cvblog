package dropbox

import (
	"os"
	"testing"
	"strings"
)

func TestListFolder(t *testing.T) {
	token := os.Getenv("TOKEN")
	if token == "" {
		t.Errorf("export TOKEN=xxxx first..")
	}

	client := New(token)

	if _, err := client.ListFolder("/Apps/hackcv"); err != nil {
		t.Errorf("list folder fail: %s", err)
	}

	if _, err := client.Download("id:YkRiURCPNyAAAAAAAAAAAg"); err != nil {
		t.Errorf("download fail: %s", err)
	}

	reader := strings.NewReader("just a test")
	if err := client.Upload("/Apps/hackcv/test.txt", reader); err != nil {
		t.Errorf("upload fail: %s", err)
	}

}
