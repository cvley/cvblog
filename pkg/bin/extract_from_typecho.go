package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Post struct {
	Cid      int
	Title    string
	Modified int64
	Text     string
}

func (post *Post) Byte() []byte {
	var buffer bytes.Buffer
	date := time.Unix(post.Modified, 0).Format("2006-01-02 15:04")
	buffer.WriteString(fmt.Sprintf("Date: %s\n", date))
	buffer.WriteString(fmt.Sprintf("Title: %s\n", post.Title))
	buffer.WriteString(fmt.Sprintf("Tags: \n"))
	buffer.WriteString(fmt.Sprintf("Status: \n"))
	buffer.WriteString(fmt.Sprintf("URL: \n\n"))

	buffer.WriteString(post.Text[15:])
	buffer.WriteString("\n")

	return buffer.Bytes()
}

func (post *Post) ToFile(w io.Writer) error {
	_, err := w.Write(post.Byte())
	return err
}

func main() {
	username := flag.String("user", "", "mysql user name")
	passwd := flag.String("passwd", "", "mysql password")
	host := flag.String("host", "127.0.0.1", "mysql host")
	port := flag.Int("port", 3306, "mysql port")
	database := flag.String("db", "blog", "db name")
	flag.Parse()

	conf := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		*username,
		*passwd,
		*host,
		*port,
		*database,
	)

	db, err := sql.Open("mysql", conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	query := "select cid, title, modified, text from typecho_contents"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Cid, &post.Title, &post.Modified, &post.Text); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		f, err := os.Create("posts/" + post.Title)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := post.ToFile(f); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	return
}
