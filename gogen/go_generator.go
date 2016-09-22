package gogenerator

import "fmt"

import (
	"io"
	def "github.com/bulletRush/qcloud-sdk/generator"
	"bytes"
	"strings"
)

const (
	MAX_ARGS_COUNT = 5
)

type GoGenerator struct {
	writer io.WriteCloser
	pkgName string
	clsName string
	errList []error
	constBuffer bytes.Buffer
	varBuffer bytes.Buffer
	funcBuffer bytes.Buffer
	indentCnt int
}

func NewGoGenerator(writer io.WriteCloser, pkgName string, clsName string) *GoGenerator {
	return &GoGenerator{
		writer: writer,
		pkgName: pkgName,
		clsName: clsName,
	}
}

func (gen *GoGenerator) Output() error {
	// TODO
	gen.writer.Write(gen.funcBuffer.Bytes())
	return nil
}

func (gen *GoGenerator) GenFuncInputParams(infDef def.InterfaceDefinition) string {
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

func (gen *GoGenerator) GenFuncOutputParams(infDef def.InterfaceDefinition) ([]string, error) {
	// TODO
	return []string{}, nil
}

func (gen *GoGenerator) GenCheckCall(paramDef def.ParamDefinition) string {
	typ := strings.Title(paramDef.Type)
	buf := bytes.NewBuffer(nil)
	var sym string
	if paramDef.Optional {
		sym = "*"
	}
	buf.Write([]byte(fmt.Sprintf("helper.Check%s(%s%s, \"%s\")", typ, sym, paramDef.Name, paramDef.Rule)))
	return buf.String()
}

func (gen *GoGenerator) indent() {
	gen.indentCnt++
}

func (gen *GoGenerator) unIndent() {
	gen.indentCnt--
	if gen.indentCnt < 0 {
		panic("inden less than zero!")
	}
}

func (gen *GoGenerator) GenInputParamCheck(paramDef def.ParamDefinition) error {
	buf := bytes.NewBuffer(nil)
	_, err := buf.WriteTo(gen.writer)
	return err
}

func star(optional bool) string {
	if optional {
		return "*"
	}
	return ""
}

func (gen *GoGenerator) repeatTab() string {
	return strings.Repeat("\t", gen.indentCnt)
}

func (gen *GoGenerator) fprintf(buf *bytes.Buffer, format string, ctx ...interface{}) (int, error) {
	n, err := fmt.Fprintf(buf, gen.repeatTab() + format, ctx...)
	return n, err
}

func (gen *GoGenerator) GenFuncDoc(infDef def.InterfaceDefinition) {
	b := &gen.funcBuffer
	gen.fprintf(b, "// %ss %s\n", infDef.Name, infDef.Brief)
	gen.fprintf(b, "// %s\n", infDef.Describe)
	for _, paramDef := range infDef.InputParamList {
		gen.fprintf(b, "//\t%s: %s\n", paramDef.Name, paramDef.Describe)
	}
}

func (gen *GoGenerator) genRequestType(infDef def.InterfaceDefinition) string {
	return infDef.Name + "Request"
}

// 生成接口入参的结构体定义
func (gen *GoGenerator) GenRequestDefinition(infDef def.InterfaceDefinition) error {
	b := &gen.funcBuffer
	typeDef := gen.genRequestType(infDef)
	gen.fprintf(b, "type %s struct {\n", typeDef)
	gen.indent()
	for _, paramDef := range infDef.InputParamList {
		if paramDef.Optional {
			gen.fprintf(b, "%s *%s\n", strings.Title(paramDef.Name), paramDef.Type)
		} else {
			gen.fprintf(b, "%s %s\n", strings.Title(paramDef.Name), paramDef.Type)
		}
	}
	gen.unIndent()
	gen.fprintf(b, "}\n")
	return nil
}

var (
	codeParam = def.ParamDefinition{
		Name: "code",
		Describe: "common code",
		Type: "int",
	}
	codeDescParam = def.ParamDefinition{
		Name: "codeDesc",
		Describe: "code desc",
		Type: "string",
	}
	messageParam = def.ParamDefinition{
		Name: "message",
		Describe: "message",
		Type: "string",
	}
	commonParamList = []def.ParamDefinition{
		codeParam, codeDescParam, messageParam,
	}
)

func (gen *GoGenerator) genResponseType(infDef def.InterfaceDefinition) string {
	return infDef.Name + "Response"
}

// 生成接口返回值的结构体定义，code/codeDesc/message自动填充
func (gen *GoGenerator) GenResponseDefinition(infDef def.InterfaceDefinition) error {
	b := &gen.funcBuffer
	typeDef := gen.genResponseType(infDef)
	gen.fprintf(b, "type %s struct{\n", typeDef)
	gen.indent()
	outParamList := commonParamList
	outParamList = append(outParamList, infDef.OutputParamList...)
	for _, paramDef := range outParamList {
		gen.fprintf(b, "%s %s%s `json:\"%s,omitempty\"`\n", strings.Title(paramDef.Name), star(paramDef.Optional), paramDef.Type, paramDef.Name)
	}
	gen.unIndent()
	gen.fprintf(b, "}\n\n")
	return nil
}

func (gen *GoGenerator) GenFunc(infDef def.InterfaceDefinition) error {
	// generate a struct if input arguments is great than MAX_ARGS_COUNT
	if len(infDef.InputParamList) > MAX_ARGS_COUNT {
		gen.GenRequestDefinition(infDef)
	}
	gen.GenResponseDefinition(infDef)
	b := &gen.funcBuffer
	gen.GenFuncDoc(infDef)
	gen.fprintf(b, "func (svc *%s) %s(%s) (*%s, error) {\n", gen.clsName, infDef.Name, gen.GenFuncInputParams(infDef), gen.genResponseType(infDef))
	gen.indent()
	gen.fprintf(b, "paramMap := map[string]interface{}{}\n")
	for _, paramDef := range infDef.InputParamList {
		if paramDef.Optional {
			gen.fprintf(b, "if %s != nil {\n", paramDef.Name)
			gen.indent()
			gen.fprintf(b, "paramMap[\"%s\"] =  *%s\n", paramDef.Name, paramDef.Name)
			gen.unIndent()
			gen.fprintf(b, "}\n")
		} else {
			gen.fprintf(b, "paramMap[\"%s\"] = %s\n", paramDef.Name, paramDef.Name)
		}
	}
	gen.fprintf(b, "rspObj := &%s{}\n", gen.genResponseType(infDef))
	gen.fprintf(b, "err := svc.DoRequest(paramMap, rspObj)\n")
	gen.fprintf(b, "if err != nil {\n")
	gen.indent()
	gen.fprintf(b, "svc.logger.Error(\"%s failed!\", \"error\", err)\n", infDef.Name)
	gen.fprintf(b, "return nil, err\n")
	gen.unIndent()
	gen.fprintf(b, "}\n")
	gen.fprintf(b, "return rspObj, nil\n")
	gen.unIndent()
	gen.fprintf(b, "}\n")
	return nil
}
