package gogenerator

import "fmt"

import (
	"io"
	def "github.com/bulletRush/qcloud-sdk/generator"
	"bytes"
	"strings"
)

type GoGenerator struct {
	writer io.WriteCloser
	pkgName string
	clsName string
	errList []error
}

func NewGoGenerator(writer io.WriteCloser, pkgName string, clsName string) *GoGenerator {
	return &GoGenerator{
		writer: writer,
		pkgName: pkgName,
		clsName: clsName,
	}
}

func (this *GoGenerator) GenFuncInputParams(infDef def.InterfaceDefinition) string {
	argList := []string{}
	for _, param := range infDef.InputParamList {
		var star string
		if param.Optional {
			star = "*"
		}
		argList = append(argList, fmt.Sprintf("%s %s%s", param.Name, star, param.Type))
	}
	return strings.Join(argList, ", ")
}

func (this *GoGenerator) GenFuncOutputParams(infDef def.InterfaceDefinition) ([]string, error) {
	// TODO
	return []string{}, nil
}

func (this *GoGenerator) GenCheckCall(paramDef def.ParamDefinition) string {
	typ := strings.Title(paramDef.Type)
	buf := bytes.NewBuffer(nil)
	var sym string
	if paramDef.Optional {
		sym = "*"
	}
	buf.Write([]byte(fmt.Sprintf("helper.Check%s(%s%s, \"%s\")", typ, sym, paramDef.Name, paramDef.Rule)))
	return buf.String()
}

func (this *GoGenerator) GenInputParamCheck(paramDef def.ParamDefinition) error {
	buf := bytes.NewBuffer(nil)
	_, err := buf.WriteTo(this.writer)
	return err
}

func star(optional bool) string {
	if optional {
		return "*"
	}
	return ""
}

func (this *GoGenerator) GenFuncDoc(infDef def.InterfaceDefinition) string {
	this.writer.Write([]byte(fmt.Sprintf("// %ss %s\n", infDef.Name, infDef.Brief)))
	this.writer.Write([]byte(fmt.Sprintf("// %s\n", infDef.Describe)))
	for _, paramDef := range infDef.InputParamList {
		this.writer.Write([]byte(fmt.Sprintf("//\t%s: %s\n", paramDef.Name, paramDef.Describe)))
	}
	return ""
}

func (this *GoGenerator) GenFunc(infDef def.InterfaceDefinition) error {
	this.GenFuncDoc(infDef)
	this.writer.Write([]byte(fmt.Sprintf("func (this *%s) %s(%s) error {\n", this.clsName, infDef.Name, this.GenFuncInputParams(infDef))))
	this.writer.Write([]byte("\tparamMap := map[string]interface{}{}\n"))
	for _, paramDef := range infDef.InputParamList {
		if paramDef.Optional {
			this.writer.Write([]byte(fmt.Sprintf("\tif %s != nil {\n", paramDef.Name)))
			this.writer.Write([]byte(fmt.Sprintf("\t\tif !%s {\n", this.GenCheckCall(paramDef))))
			this.writer.Write([]byte(fmt.Sprintf("\t\t\treturn helper.NewParamError(\"%s\")\n", paramDef.Name)))
			this.writer.Write([]byte(fmt.Sprintf("\t\t}\n")))
			this.writer.Write([]byte(fmt.Sprintf("\t\thelper.Process%s(paramMap, %s%s)\n", strings.Title(paramDef.Type), star(paramDef.Optional), paramDef.Name)))
			this.writer.Write([]byte(fmt.Sprintf("\t}\n")))
		} else {
			this.writer.Write([]byte(fmt.Sprintf("\tif !%s {\n", this.GenCheckCall(paramDef))))
			this.writer.Write([]byte(fmt.Sprintf("\t\treturn helper.NewParamError(\"%s\")\n", paramDef.Name)))
			this.writer.Write([]byte(fmt.Sprintf("\t}\n")))
			this.writer.Write([]byte(fmt.Sprintf("\thelper.Process%s(paramMap, %s%s)\n", strings.Title(paramDef.Type), star(paramDef.Optional), paramDef.Name)))
		}
	}
	this.writer.Write([]byte("\treturn helper.Do(paramMap)\n}\n"))
	return nil
}

