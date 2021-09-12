package yktr

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/winebarrel/yktr/esa"
)

//go:embed templates
var content embed.FS

type Server struct {
	*Config
	engine *gin.Engine
}

func NewServer(cfg *Config) (*Server, error) {
	r := gin.New()
	t, err := template.ParseFS(content, "templates/index.html")

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

	c.HTML(http.StatusOK, "index.html", gin.H{
		"q":          q,
		"domain":     fmt.Sprintf("%s.esa.io", cfg.Team),
		"stylesheet": esa.Stylesheet,
		"posts":      posts.Posts,
		"prev_page":  posts.PrevPage,
		"next_page":  posts.NextPage,
	})
}
