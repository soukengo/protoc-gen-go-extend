package setters

import (
	"fmt"
	"go/token"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	name           = "setters"
	defaultVarName = "v"
)

var tpl = `func (x * %s) Set%s(%s %s) {
	if x == nil {
		return
	}
	x.%s = %s
}
`

type Plugin struct {
}

func (p *Plugin) Name() string {
	return name
}

func (p *Plugin) Generate(g *protogen.GeneratedFile, file *protogen.File) {
	messages := file.Messages
	for _, message := range messages {
		p.genMessageCode(g, message)
	}
}

func (p *Plugin) genMessageCode(g *protogen.GeneratedFile, m *protogen.Message) {
	fields := m.Fields
	for _, field := range fields {
		var typeName = fieldGoType(g, field)
		var fieldName = field.GoName
		var varName = field.Desc.JSONName()
		if t := token.Lookup(varName); t.IsKeyword() {
			varName = defaultVarName
		}
		g.P(fmt.Sprintf(tpl, m.GoIdent.GoName, fieldName, varName, typeName, fieldName, varName))
	}
}

func fieldGoType(g *protogen.GeneratedFile, field *protogen.Field) (goType string) {
	if field.Desc.IsWeak() {
		return "struct{}"
	}
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		goType = "bool"
	case protoreflect.EnumKind:
		goType = g.QualifiedGoIdent(field.Enum.GoIdent)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		goType = "int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		goType = "uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		goType = "int64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		goType = "uint64"
	case protoreflect.FloatKind:
		goType = "float32"
	case protoreflect.DoubleKind:
		goType = "float64"
	case protoreflect.StringKind:
		goType = "string"
	case protoreflect.BytesKind:
		goType = "[]byte"
	case protoreflect.MessageKind, protoreflect.GroupKind:
		goType = "*" + g.QualifiedGoIdent(field.Message.GoIdent)
	}
	switch {
	case field.Desc.IsList():
		return "[]" + goType
	case field.Desc.IsMap():
		keyType := fieldGoType(g, field.Message.Fields[0])
		valType := fieldGoType(g, field.Message.Fields[1])
		return fmt.Sprintf("map[%v]%v", keyType, valType)
	}
	return goType
}
