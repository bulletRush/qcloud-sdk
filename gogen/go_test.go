package gogenerator

import (
	"testing"
	"os"
	def "github.com/bulletRush/qcloud-sdk/generator"
	"fmt"
)

var (
	diskTypeParamDef = def.ParamDefinition{
		Name: "diskType",
		Type: "string",
		Rule: "root|data",
		Optional: true,
		Describe: "disk usage type",
	}
	payModeParamDef = def.ParamDefinition{
		Name: "payMode",
		Type: "string",
		Rule: "prePay|postPay",
		Optional: true,
		Describe: "disk pay mode",
	}
	storageTypeParamDef = def.ParamDefinition{
		Name: "storageType",
		Type: "string",
		Rule: "cloudBasic|cloudSSD",
		Optional: false,
		Describe: "storage type",
	}
	interfaceDef = def.InterfaceDefinition{
		Name: "DescribeCbsStorages",
		Brief: "list cbs storages",
		Describe: "see xx for more",
		InputParamList: []def.ParamDefinition{
			diskTypeParamDef, payModeParamDef, storageTypeParamDef,
		},
	}
)

func TestNewGoGenerator(t *testing.T) {

}

func newGoGenerator() *GoGenerator {
	return NewGoGenerator(os.Stdout, "demo", "DemoService")
}

func xTestGoGenerator_GenInputParamCheck(t *testing.T) {
	gg := newGoGenerator()
	gg.GenInputParamCheck(diskTypeParamDef)
}

func xTestGoGenerator_GenCheckCall(t *testing.T) {
	gg := newGoGenerator()
	fmt.Println(gg.GenCheckCall(payModeParamDef))
}

func TestGoGenerator_GenFunc(t *testing.T) {
	gg := newGoGenerator()
	fmt.Println(gg.GenFunc(interfaceDef))
}
