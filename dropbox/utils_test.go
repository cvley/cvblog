package dropbox

import (
	"testing"
)

func TestParseEntries(t *testing.T) {
	testmsg := `{"entries": [{".tag": "folder", "name": "test", "path_lower": "/apps/hackcv/test", "path_display": "/Apps/hackcv/test", "id": "id:XZWT2grs02AAAAAAAAAAAQ" }, { ".tag": "file", "name": "Linux生日.md", "path_lower": "/apps/hackcv/linux生日.md", "path_display": "/Apps/hackcv/Linux生日.md", "id": "id:YPUcpkGDHuAAAAAAAAAAAQ", "client_modified": "2013-09-17T11:33:48Z", "server_modified": "2016-01-04T07:41:44Z", "rev": "142b038dd", "size": 2671 } ], "cursor": "AAFnN2fMkjkw7Pwf4YuY-zeERp7KvsXc6y6eUjHOs2ARUenw2d_qJE5lWHZOzp5CWUvPTS1puOrRkOf6OC7KtfE5fJPjyfPz8izfNiMM2G4BKzr2H23z1MY9seHp0ctwj4XqX1YvkwGxSCfyHj1BNE0_QORcvJeIxnlIgd_Pcn3Hr42eFvBN_3ufG_mxljbNff0", "has_more": false}`

	entries, err := parseEntries([]byte(testmsg))
	if err != nil {
		t.Errorf("parse entries fail: %s", err)
	}

	for _, entry := range entries {
		if entry.Id != "id:XZWT2grs02AAAAAAAAAAAQ" && entry.Id != "id:YPUcpkGDHuAAAAAAAAAAAQ" {
			t.Errorf("parse entries fail: invalid id")
		}
		if entry.Name != "test" && entry.Name != "Linux生日.md" {
			t.Errorf("parse entries fail: invalid name")
		}
		if entry.Tag != "folder" && entry.Tag != "file" {
			t.Errorf("parse entries fail: invalid tag")
		}
	}
}
