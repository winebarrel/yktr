// cf. https://docs.esa.io/posts/102
package esa

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/bxcodec/httpcache"
	"github.com/gin-gonic/gin"
)

const (
	Endpoint       = "api.esa.io"
	Stylesheet     = "https://assets.esa.io/assets/application-860deb72f57963abb3cecce7b8070ab4e106b68cee8e3205d457110507b494f4.css"
	DefaultPerPage = "5"
)

type Client struct {
	*Config
}

func NewClient(cfg *Config) *Client {
	cli := &Client{
		Config: cfg,
	}

	return cli
}

type Author struct {
	Myself     bool
	Name       string
	ScreenName string `json:"screen_name"`
	Icon       string
}

type Post struct {
	Number         int
	Name           string
	FullName       string `json:"full_name"`
	Wip            bool
	BodyMd         string        `json:"body_md"`
	BodyHtml       template.HTML `json:"body_html"`
	CreatedAt      string        `json:"created_at"`
	Message        string
	Url            string
	UpdatedAt      string `json:"updated_at"`
	Tags           []string
	Category       string
	RevisionNumber int    `json:"revision_number"`
	CreatedBy      Author `json:"created_by"`
	UpdatedBy      Author `json:"updated_by"`
}

type Posts struct {
	Posts      []Post
	PrevPage   int `json:"prev_page"`
	NextPage   int `json:"next_page"`
	TotalCount int `json:"total_count"`
	Page       int
	PerPage    int `json:"per_page"`
	MaxPerPage int `json:"max_per_page"`
}

func (cli *Client) Posts(c *gin.Context, q string) (*Posts, error) {
	url := fmt.Sprintf("https://%s/v1/teams/%s/posts", Endpoint, cli.Team)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("sort", "number")
	query.Add("order", "desc")

	if q != "" {
		query.Add("q", q)
	}

	req.URL.RawQuery = query.Encode()

	body, err := cli.request(c, req)

	if err != nil {
		return nil, err
	}

	posts := &Posts{}
	err = json.Unmarshal(body, &posts)

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (cli *Client) request(c *gin.Context, req *http.Request) ([]byte, error) {
	req.Header.Add("Authorization", "Bearer "+cli.Token)

	page, err := strconv.Atoi(c.Query("page"))
	query := req.URL.Query()
	query.Add("per_page", DefaultPerPage)

	if err == nil && page > 0 {
		query.Add("page", strconv.Itoa(page))
	}

	req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	_, err = httpcache.NewWithInmemoryCache(client, false, time.Second*60)

	if err != nil {
		return nil, err
	}

	fmt.Println(req)
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", res.Status, body)
	}

	return body, err
}
