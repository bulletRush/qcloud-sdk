package main

import (
	"os"
	"fmt"
	def "github.com/bulletRush/qcloud-sdk/generator/common"
	"github.com/bulletRush/qcloud-sdk/generator/gogen"
)

func main() {
	svcDef, err := def.LoadServiceDefinitionFromFile("../qcloud/cbs.json")
	if err != nil {
		fmt.Sprint(err)
		return
	}
	fout, err := os.OpenFile("../qcloud/cbs.auto.go", os.O_WRONLY | os.O_CREATE, 0644)
	gen := gogen.NewGoGenerator(
		fout, "qcloud", "CbsService",
	)
	gen.GenService(*svcDef)
}
