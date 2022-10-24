package setters

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"strings"
)

func getFieldType(g *generator.Generator, field *descriptor.FieldDescriptorProto) string {
	var varType string
	switch *field.Type {
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		varType = getMessageFieldType(g, field)
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		desc := g.ObjectNamed(field.GetTypeName())
		varType = g.TypeName(desc)
	default:
		varType = strings.ToLower(strings.Split(field.Type.String(), "_")[1])
	}
	return varType
}

func getMessageFieldType(g *generator.Generator, field *descriptor.FieldDescriptorProto) string {
	desc := g.ObjectNamed(field.GetTypeName())
	if d, ok := desc.(*generator.Descriptor); ok && d.GetOptions().GetMapEntry() {
		keyField, valField := d.Field[0], d.Field[1]
		keyType, _ := g.GoType(d, keyField)
		valType, _ := g.GoType(d, valField)
		return fmt.Sprintf("map[%s]%s", keyType, valType)
	}
	return "*" + g.TypeName(desc)
}

func upperFirstLatter(s string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) == 1 {
		return strings.ToUpper(string(s[0]))
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}
