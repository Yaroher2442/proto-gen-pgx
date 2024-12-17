package domain

import (
	customDescriptor "github.com/Yaroher2442/protoc-gen-gopgx/internal/domain/descriptorpb"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Method struct {
	Name            string
	ExecQuery       *Query
	ExecTransaction *Transaction
}

func NewMethod(method *protogen.Method) (*Method, error) {
	methodOptions := proto.GetExtension(
		method.Desc.Options().(*descriptorpb.MethodOptions),
		customDescriptor.E_Exec,
	).(*customDescriptor.ServiceMethodExec)
	meth := &Method{
		Name: method.GoName,
	}
	switch methodOptions.Exec.(type) {
	case *customDescriptor.ServiceMethodExec_Query:
		q := methodOptions.Exec.(*customDescriptor.ServiceMethodExec_Query).Query
		qSQL, err := NewQuery(q)
		if err != nil {
			return nil, err
		}
		meth.ExecQuery = qSQL
	case *customDescriptor.ServiceMethodExec_Transaction:
		tx := methodOptions.Exec.(*customDescriptor.ServiceMethodExec_Transaction).Transaction
		txSQL, err := NewTransaction(tx)
		if err != nil {
			return nil, err
		}
		meth.ExecTransaction = txSQL
	}
	return meth, nil
}

type Service struct {
	Name    string
	Methods []*Method
}

func NewService(service *protogen.Service) (*Service, error) {
	methods := make([]*Method, len(service.Methods))
	for i, method := range service.Methods {
		m, err := NewMethod(method)
		if err != nil {
			return nil, err
		}
		methods[i] = m
	}
	return &Service{Name: service.GoName}, nil
}

func NewServices(services []*protogen.Service) ([]*Service, error) {
	svcs := make([]*Service, len(services))
	for i, service := range services {
		svc, err := NewService(service)
		if err != nil {
			return nil, err
		}
		svcs[i] = svc
	}
	return svcs, nil
}
