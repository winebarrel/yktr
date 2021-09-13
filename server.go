package yktr

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/gliderlabs/sigil/builtin"
	"github.com/winebarrel/yktr/esa"
)

//go:embed templates
var tpls embed.FS

//go:embed assets
var assets embed.FS

var sigilFuncMap = template.FuncMap{
	"append":       builtin.Append,
	"base64decode": builtin.Base64Decode,
	"base64encode": builtin.Base64Encode,
	"capitalize":   builtin.Capitalize,
	"default":      builtin.Default,
	"dir":          builtin.Dir,
	"dirs":         builtin.Dirs,
	"drop":         builtin.Drop,
	"exists":       builtin.Exists,
	"file":         builtin.File,
	"files":        builtin.Files,
	"httpget":      builtin.HttpGet,
	"include":      builtin.Include,
	"indent":       builtin.Indent,
	"jmespath":     builtin.JmesPath,
	"join":         builtin.Join,
	"joinkv":       builtin.JoinKv,
	"json":         builtin.Json,
	"lower":        builtin.Lower,
	"match":        builtin.Match,
	"pointer":      builtin.Pointer,
	"render":       builtin.Render,
	"replace":      builtin.Replace,
	"seq":          builtin.Seq,
	"shell":        builtin.Shell,
	"split":        builtin.Split,
	"splitkv":      builtin.SplitKv,
	"stdin":        builtin.Stdin,
	"substring":    builtin.Substring,
	"text":         builtin.Text,
	"tojson":       builtin.ToJson,
	"toyaml":       builtin.ToYaml,
	"trim":         builtin.Trim,
	"uniq":         builtin.Uniq,
	"upper":        builtin.Upper,
	"var":          builtin.Var,
	"yaml":         builtin.Yaml,
	"emoji":        emoji,
}

type Server struct {
	*Config
	engine *gin.Engine
}

var redirectPaths []string = []string{
	"posts",
	"members",
}

func NewServer(cfg *Config) (*Server, error) {
	r := gin.Default()
	store := persistence.NewInMemoryStore(cfg.CacheTTL * time.Second)
	err := letHTMLTemplates(r)

	if err != nil {
		return nil, err
	}

	esaCli := esa.NewClient(&esa.Config{
		Team:    cfg.Team,
		Token:   cfg.Token,
		PerPage: cfg.PerPage,
	})

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("assets/favicon.ico", http.FS(assets))
	})

	for _, path := range redirectPaths {
		path := path
		urlPath := fmt.Sprintf("/%s/:id", path)

		r.GET(urlPath, func(c *gin.Context) {
			url := fmt.Sprintf("http://%s.esa.io/%s/%s", cfg.Team, path, c.Param("id"))
			c.Redirect(http.StatusTemporaryRedirect, url)
		})
	}

	handler := func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		handleReq(c, cfg, esaCli)
	}

	if cfg.CacheTTL > 0 {
		handler = cache.CachePage(store, cfg.CacheTTL*time.Second, handler)
	}

	r.NoRoute(handler)

	return &Server{Config: cfg, engine: r}, nil
}

func (svr *Server) Run() error {
	listen := fmt.Sprintf("%s:%d", svr.Addr, svr.Port)
	return svr.engine.Run(listen)
}

func handleReq(c *gin.Context, cfg *Config, esaCli *esa.Client) {
	q := c.Query("q")
	category := strings.TrimLeft(c.Request.URL.Path, "/")

	if category != "" {
		q = fmt.Sprintf(`%s in:"%s"`, q, category)
	}

	posts, err := esaCli.Posts(c, strings.TrimSpace(q))

	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		return
	}

	data := gin.H{
		"q":          c.Query("q"),
		"category":   category,
		"domain":     fmt.Sprintf("%s.esa.io", cfg.Team),
		"stylesheet": esa.Stylesheet,
		"posts":      posts.Posts,
		"prev_page":  posts.PrevPage,
		"next_page":  posts.NextPage,
	}

	if len(posts.Posts) == 0 {
		c.HTML(http.StatusOK, "post_not_found.html", data)
		return
	}

	c.HTML(http.StatusOK, "index.html", data)
}

func letHTMLTemplates(r *gin.Engine) error {
	t := template.New("").Funcs(sigilFuncMap)

	entries, err := tpls.ReadDir("templates")

	if err != nil {
		return err
	}

	for _, e := range entries {
		t, err := t.ParseFS(tpls, path.Join("templates", e.Name()))

		if err != nil {
			return err
		}

		r.SetHTMLTemplate(t)
	}

	return nil
}

func emoji(str string) (interface{}, error) {
	r := regexp.MustCompile(":([[:alnum:]]+):")
	str = r.ReplaceAllString(str, `<img class="emoji" title=":$1:" alt=":$1:" src="https://assets.esa.io/images/emoji/$1.png">`)
	return template.HTML(str), nil
}
