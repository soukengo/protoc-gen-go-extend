package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/soukengo/protoc-gen-go-extend/setters"
	"io/ioutil"
	"os"
)

func init() {
	generator.RegisterPlugin(new(setters.Plugin))
}

func main() {
	g := generator.New()

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		g.Error(err, "reading input")
	}

	if err := proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err, "parsing input proto")
	}
	for _, file := range g.Request.FileToGenerate {
		g1 := generator.New()
		g1.Request.FileToGenerate = []string{file}
		g1.Request.Parameter = g.Request.Parameter
		g1.Request.ProtoFile = g.Request.ProtoFile
		g1.Request.CompilerVersion = g.Request.CompilerVersion
		generate(g1)
	}

}

func generate(g *generator.Generator) {
	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}
	g.CommandLineParameters(g.Request.GetParameter())

	// Create a wrapped version of the Descriptors and EnumDescriptors that
	// point to the file that defines them.
	g.WrapTypes()

	g.SetPackageNames()
	g.BuildTypeNameMap()

	g.GenerateAllFiles()

	// Send back the results.
	data, err := proto.Marshal(g.Response)
	if err != nil {
		g.Error(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		g.Error(err, "failed to write output proto")
	}
}
