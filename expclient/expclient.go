package expclient

import (
	"bytes"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	exec *http.Client
}

// for gomock
type Poster interface {
	GetPosts(id string) (*Resp, error)
}

func NewPoster(c *http.Client) Poster {
	return &Client{
		exec: c,
	}
}


const (
	urlMain = "https://jsonplaceholder.typicode.com/"
	postPath = urlMain + "post/"
	)

type Resp struct {
	User    int      `json:"user"`
	Id      int      `json:"id"`
	Title   string   `json:"title"`
	Body    string   `json:"body"`
}


func (r *Resp) ToUpperBody() {
	r.Body = strings.ToUpper(r.Body)
}


func New(exec *http.Client) *Client {
	if exec == nil {
		exec = http.DefaultClient
	}
	return &Client{
		exec:exec,
	}
}

func (c Client) GetPosts(id string) (*Resp, error) {

	buff := bytes.NewReader(make([]byte, 1024))

	reqUrl := "https://jsonplaceholder.typicode.com/posts/" + id

	log.Debug().Str("req_path", reqUrl).Msg("")

	req, err := http.NewRequest(http.MethodGet, reqUrl, buff)
	if err != nil {
		return nil, err
	}

	log.Debug().Str("url", req.URL.String()).Msg("requesting_external_api")

	resp, err := c.exec.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	out := new(Resp)

	err = json.Unmarshal(respBody, out)
	if err != nil {
		return nil, err
	}

	out.ToUpperBody()

	return out, nil
}