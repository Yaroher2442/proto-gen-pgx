package protocgen

import (
	"github.com/Yaroher2442/protoc-gen-gopgx/internal/domain"
	"github.com/Yaroher2442/protoc-gen-gopgx/internal/protogenhelp"
	"google.golang.org/protobuf/compiler/protogen"
	"path/filepath"
	"time"
)

type Config struct {
	SqlPath string
}

type Generator struct {
	plugin *protogen.Plugin
	files  []*domain.File
	config *Config
}

func NewGenerator(config *Config, plugin *protogen.Plugin) (*Generator, error) {
	abs, err := filepath.Abs(config.SqlPath)
	if err != nil {
		return nil, err
	}
	config.SqlPath = abs
	err = protogenhelp.RegisterALLFilesExtension(plugin)
	if err != nil {
		return nil, err
	}
	files, err := domain.NewFilesFromProtoPlugin(plugin)
	if err != nil {
		return nil, err
	}
	return &Generator{
		plugin: plugin,
		files:  files,
		config: config,
	}, nil
}

func (g *Generator) RenderSql() error {
	gf := g.plugin.NewGeneratedFile(filepath.Join(
		"migrations", time.Now().Format("20060102150405")+".sql",
	), "")
	for _, file := range g.files {
		err := file.RenderMigrations(gf)
		gf.P()
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) Generate() error {
	err := g.RenderSql()
	if err != nil {
		return err
	}
	return nil
}
