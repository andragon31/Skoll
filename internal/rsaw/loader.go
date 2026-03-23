package rsaw

import (
	"github.com/charmbracelet/log"
	"github.com/yuin/goldmark"
)

type Loader struct {
	root   string
	gm     goldmark.Markdown
	logger *log.Logger
}

func NewLoader(logger *log.Logger) *Loader {
	return &Loader{
		root:   ".",
		gm:     goldmark.New(),
		logger: logger,
	}
}
