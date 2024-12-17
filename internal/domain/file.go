package domain

import (
	"errors"
	customDescriptor "github.com/Yaroher2442/protoc-gen-gopgx/internal/domain/descriptorpb"
	"github.com/hashicorp/go-version"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

type File struct {
	File          *protogen.File
	Version       *version.Version
	CommonQueries map[string]*Query
	Services      []*Service
	Tables        []*Table
	Enums         []*Enum
}

func (f *File) RenderMigrations(file *protogen.GeneratedFile) error {
	for _, table := range f.Tables {
		err := table.RenderSchema(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *File) RenderQueries(file *protogen.GeneratedFile) error {
	return nil
}

var ErrFileNotSql = errors.New("file is not sql")

func NewFile(file *protogen.File) (*File, error) {
	fileOptions, ok := proto.GetExtension(
		file.Desc.Options(),
		customDescriptor.E_FileOptions,
	).(*customDescriptor.FileOptions)
	if !ok || fileOptions == nil {
		return nil, ErrFileNotSql
	}
	queries := make(map[string]*Query)
	if fileOptions.GetCommonQueries() != nil {
		for _, q := range fileOptions.GetCommonQueries().GetQueries() {
			queries[q.GetName()] = NewSqlQuery(q.GetSql())
		}
	}
	fileV, err := version.NewSemver(fileOptions.GetVersion())
	if err != nil {
		return nil, err
	}
	services, err := NewServices(file.Services)
	if err != nil {
		return nil, err
	}
	tables, err := NewTables(file.Messages)
	if err != nil {
		return nil, err
	}
	return &File{
		File:          file,
		CommonQueries: queries,
		Version:       fileV,
		Services:      services,
		Tables:        tables,
		Enums:         NewEnums(file.Enums),
	}, nil
}

func NewFilesFromProtoPlugin(plugin *protogen.Plugin) ([]*File, error) {
	files := make([]*File, 0)
	for _, f := range plugin.Files {
		file, err := NewFile(f)
		if err != nil {
			if errors.Is(err, ErrFileNotSql) {
				continue
			}
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}
