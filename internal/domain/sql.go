package domain

import (
	"errors"
	customDescriptor "github.com/Yaroher2442/protoc-gen-gopgx/internal/domain/descriptorpb"
	"github.com/Yaroher2442/protoc-gen-gopgx/internal/helpers"
	"github.com/samber/lo"
	"github.com/sttk/stringcase"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
)

type Query struct {
	SQL string
	Ref string
}

func NewSqlQuery(sql string) *Query {
	return &Query{SQL: sql, Ref: ""}
}

func NewQuery(query *customDescriptor.MethodQuery) (*Query, error) {
	switch query.Query.(type) {
	case *customDescriptor.MethodQuery_Sql:
		if query.GetSql() == "" {
			return nil, errors.New("empty sql")
		}
		return &Query{SQL: query.GetSql(), Ref: ""}, nil
	case *customDescriptor.MethodQuery_Ref:
		if query.GetRef() == "" {
			return nil, errors.New("empty ref")
		}
		return &Query{Ref: query.GetRef(), SQL: ""}, nil
	default:
		return nil, errors.New("unknown query type")
	}
}
func (q *Query) IsRef() bool {
	return q.Ref != ""
}

func (q *Query) IsSQL() bool {
	return q.SQL != ""
}

type Transaction struct {
	Level   customDescriptor.Transaction_Level
	Queries []*Query
}

func NewTransaction(tx *customDescriptor.Transaction) (*Transaction, error) {
	qs := make([]*Query, 0)
	for _, q := range tx.GetQuerySet() {
		query, err := NewQuery(q)
		if err != nil {
			return nil, err
		}
		qs = append(qs, query)
	}
	if len(qs) == 0 {
		return nil, errors.New("empty transaction")
	}
	return &Transaction{
		Level:   tx.GetLevel(),
		Queries: make([]*Query, 0),
	}, nil
}

type Enum struct {
	Enum          *protogen.Enum
	StoreAsString customDescriptor.EnumStoredAs
}

func NewEnumField(enum *protogen.Enum) *Enum {
	option, ok := proto.GetExtension(
		enum.Desc.Options(),
		customDescriptor.E_StoredAs,
	).(customDescriptor.EnumStoredAs)
	if !ok {
		option = customDescriptor.EnumStoredAs_ENUM_STORED_AS_INT
	}
	return &Enum{
		Enum:          enum,
		StoreAsString: option,
	}
}

func NewEnums(enums []*protogen.Enum) []*Enum {
	var res []*Enum
	for _, e := range enums {
		res = append(res, NewEnumField(e))
	}
	return res
}

type PGType interface {
}

type PGInt struct {
}

type PGBool struct {
}

type PGText struct {
}

type TableRef struct {
	TargetTable *Table
}

type Field struct {
	Filed      *protogen.Field
	PgType     PGType
	Constraint []customDescriptor.Constraint
	TableRef   *TableRef
}

func (f *Field) RenderSql(file *protogen.GeneratedFile) string {
	constraints := strings.Join(lo.Map(f.Constraint, func(c customDescriptor.Constraint, _ int) string {
		return strings.ReplaceAll(c.String(), "CONSTRAINT_", "")
	}), ",")
	if f.TableRef != nil {
		return "table_ref_here"
	}
	return strings.Join(
		[]string{
			string(f.Filed.Desc.Name()),
			f.Filed.Desc.Kind().String(),
			constraints,
		},
		" ",
	)
}

func NewField(field *protogen.Field) (*Field, error) {
	option, ok := proto.GetExtension(
		field.Desc.Options(),
		customDescriptor.E_Column,
	).(*customDescriptor.Column)
	if !ok {
		return nil, errors.New("options error in field")
	}
	//var pgType PGType = nil
	var tableRef *TableRef
	switch field.Desc.Kind() {
	case protoreflect.MessageKind:
		tableRef = &TableRef{
			TargetTable: nil,
		}
	}
	//if pgType==nil{
	//	return nil, errors.New("options error in field")
	//}
	return &Field{
		Filed: field,
		//PgType:     NewFieldPgType(option),
		Constraint: option.GetConstraints(),
		TableRef:   tableRef,
	}, nil
}

type Table struct {
	Table       *protogen.Message
	Annotations []string
	Fields      []*Field
}

func (t *Table) RenderSchema(file *protogen.GeneratedFile) error {
	file.P("CREATE TABLE IF NOT EXISTS ", stringcase.SnakeCase(string(t.Table.Desc.FullName())), " (")
	for _, field := range t.Fields {
		file.P(helpers.ToAnyList(field.RenderSql(file))...)
	}
	file.P(strings.Join(lo.Map(t.Fields, func(a *Field, _ int) string {
		return field.RenderSql(file)
	})), ",")
	return nil
}

func NewTable(table *protogen.Message) (*Table, error) {
	fields := make([]*Field, 0)
	for _, field := range table.Fields {
		f, err := NewField(field)
		if err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}
	t := &Table{
		Table:  table,
		Fields: fields,
	}
	tableOptions := proto.GetExtension(
		table.Desc.Options(),
		customDescriptor.E_Schema,
	).(*customDescriptor.SchemaOption)
	if tableOptions != nil {
		t.Annotations = tableOptions.GetAnnotations()
	}
	return t, nil
}

func NewTables(messages []*protogen.Message) ([]*Table, error) {
	tables := make([]*Table, 0)
	for _, message := range messages {
		table, err := NewTable(message)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}
