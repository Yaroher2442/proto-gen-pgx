package main

import (
	"flag"
	"github.com/Yaroher2442/protoc-gen-gopgx/pkg/protocgen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	var flags flag.FlagSet
	sqlPath := flags.String("sql_path", "./sql", "--sql_path")
	opts := &protogen.Options{
		ParamFunc: flags.Set,
	}
	opts.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		generator, err := protocgen.NewGenerator(&protocgen.Config{
			SqlPath: *sqlPath,
		}, plugin)
		if err != nil {
			return err
		}
		return generator.Generate()
	})
}
