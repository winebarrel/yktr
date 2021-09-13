package yktr

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gliderlabs/sigil/builtin"
	"github.com/winebarrel/yktr/esa"
)

//go:embed templates
var content embed.FS

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

func NewServer(cfg *Config) (*Server, error) {
	r := gin.Default()
	t := template.New("").Funcs(sigilFuncMap)
	t, err := t.ParseFS(content, "templates/index.html")

	if err != nil {
		return nil, err
	}

	r.SetHTMLTemplate(t)

	esaCli := esa.NewClient(&esa.Config{
		Team:  cfg.Team,
		Token: cfg.Token,
	})

	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		handleReq(c, cfg, esaCli)
	})

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

	if len(posts.Posts) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"q":          c.Query("q"),
		"category":   category,
		"domain":     fmt.Sprintf("%s.esa.io", cfg.Team),
		"stylesheet": esa.Stylesheet,
		"posts":      posts.Posts,
		"prev_page":  posts.PrevPage,
		"next_page":  posts.NextPage,
	})
}

func emoji(str string) (interface{}, error) {
	r := regexp.MustCompile(":([[:alnum:]]+):")
	str = r.ReplaceAllString(str, `<img class="emoji" title=":$1:" alt=":$1:" src="https://assets.esa.io/images/emoji/$1.png">`)
	return template.HTML(str), nil
}
