package setters

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"go/token"
)

const (
	name      = "setters"
	varSuffix = "_"
)

var tpl = `func (x * %s) Set%s(%s %s) {
	if x == nil {
		return
	}
	x.%s = %s
}
`

type Plugin struct {
	gen *generator.Generator
}

func (p *Plugin) Name() string {
	return name
}

func (p *Plugin) Init(g *generator.Generator) {
	p.gen = g
}

func (p *Plugin) Generate(file *generator.FileDescriptor) {
	messages := file.MessageType
	for _, message := range messages {
		p.genMessageCode(message)
	}
}

func (p *Plugin) GenerateImports(file *generator.FileDescriptor) {

}

func (p *Plugin) genMessageCode(m *descriptor.DescriptorProto) {
	fields := m.Field
	for _, field := range fields {
		var typeName = getFieldType(p.gen, field)
		var fieldName = generator.CamelCase(*field.JsonName)
		var varName = *field.JsonName
		if t := token.Lookup(varName); t.IsKeyword() {
			varName = varName + varSuffix
		}
		p.gen.P(fmt.Sprintf(tpl, *m.Name, fieldName, varName, typeName, fieldName, varName))
	}
}
