package main

import (
	"github.com/alecthomas/kong"
	_ "github.com/joho/godotenv/autoload"

	"github.com/sjansen/go-saml-demo/internal/config"
	"github.com/sjansen/go-saml-demo/internal/server"
)

type runserverCmd struct{}

type context struct {
	cfg *config.Config
}

var cli struct {
	Runserver runserverCmd `cmd:"cmd"`
}

func main() {
	ctx := kong.Parse(&cli)

	cfg, err := config.New()
	ctx.FatalIfErrorf(err)

	err = ctx.Run(&context{
		cfg: cfg,
	})
	ctx.FatalIfErrorf(err)
}

func (cmd *runserverCmd) Run(ctx *context) error {
	s, err := server.New(ctx.cfg)
	if err != nil {
		return err
	}
	return s.ListenAndServe("127.0.0.1:8080")
}
