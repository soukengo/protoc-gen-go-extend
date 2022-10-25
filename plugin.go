package main

import "google.golang.org/protobuf/compiler/protogen"

type Plugin interface {
	Name() string
	Generate(g *protogen.GeneratedFile, file *protogen.File)
}

var activePlugins = map[string]Plugin{}

func registerPlugin(p Plugin) {
	activePlugins[p.Name()] = p
}
