package dropbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	Host        = "https://api.dropboxapi.com"
	DownloadURI = "/2/files/list_folder"
)

// Entry is one item return from list_folder call
//  Tag is the type of the file, such as `folder` or `file`
//  Name is the name of the file
//  Id is the unique id of the file from the dropbox, used to download.
type Entry struct {
	Tag  string
	Name string
	Id   string
}

type Client struct {
	token string
}

func New(token string) *Client {
	authToken := fmt.Sprintf("Bearer %s", token)
	return &Client{token: authToken}
}

func (c *Client) ListFolder(folder string) ([]*Entry, error) {
	requrl := "https://api.dropboxapi.com/2/files/list_folder"

	msg := make(map[string]string)
	msg["path"] = folder
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", requrl, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Close = true

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Printf("%+v\n\n", resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status %d not 200 OK", resp.StatusCode)
	}

	return parseEntries(body)
}

func (c *Client) Download(fileid string) ([]byte, error) {
	requrl := "https://content.dropboxapi.com/2/files/download"

	msg := make(map[string]string)
	msg["path"] = fileid
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", requrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.token)
	req.Header.Set("Dropbox-API-Arg", string(b))
	req.Close = true

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status %d not 200 OK", resp.StatusCode)
	}

	return body, nil
}

func (c *Client) Upload(path string, file io.Reader) error {
	requrl := "https://content.dropboxapi.com/2/files/upload"

	msg := make(map[string]string)
	msg["path"] = path
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", requrl, file)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.token)
	req.Header.Set("Dropbox-API-Arg", string(b))
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Close = true

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status %d not 200 OK", resp.StatusCode)
	}

	result := make(map[string]interface{})
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}
	if _, exist := result["id"]; !exist {
		return fmt.Errorf("invalid response %s", string(body))
	}

	return nil

}
