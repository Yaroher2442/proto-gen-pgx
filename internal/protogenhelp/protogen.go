package protogenhelp

import (
	"errors"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
)

var ExtTypes = new(protoregistry.Types)

func RegisterAllExtensions(descs interface {
	Messages() protoreflect.MessageDescriptors
	Extensions() protoreflect.ExtensionDescriptors
}) error {
	mds := descs.Messages()
	for i := 0; i < mds.Len(); i++ {
		err := RegisterAllExtensions(mds.Get(i))
		if err != nil {
			return err
		}
	}

	xds := descs.Extensions()
	for i := 0; i < xds.Len(); i++ {
		if err := ExtTypes.RegisterExtension(dynamicpb.NewExtensionType(xds.Get(i))); err != nil {
			return err
		}
	}

	return nil
}

func RegisterALLFilesExtension(plugin *protogen.Plugin) error {
	for _, file := range plugin.Files {
		if err := RegisterAllExtensions(file.Desc); err != nil {
			return err
		}
	}
	return nil
}

type Refixed interface {
	Reset()
	proto.Message
}

func RefixOptions(m Refixed, extTypes *protoregistry.Types) (Refixed, error) {
	if m == nil {
		return nil, errors.New("message is nil")
	}
	b, err := proto.Marshal(m)
	if err != nil {
		panic(err)
	}
	m.Reset()
	err = proto.UnmarshalOptions{Resolver: extTypes}.Unmarshal(b, m)
	if err != nil {
		panic(err)
	}
	return m, nil
}

func WithRefixOptions[T Refixed](m T) (T, error) {
	opts, err := RefixOptions(m, ExtTypes)
	return opts.(T), err
}
