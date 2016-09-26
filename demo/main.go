package main

import (
	"github.com/bulletRush/qcloud-sdk/qcloud"
	"fmt"
)
func main() {
	engine := qcloud.NewQcloudEngine()
	cbs := qcloud.NewCbsService(engine)
	rsp, err := cbs.DescribeCbsStorages(nil, nil)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}
	fmt.Printf("rsp: %#v\n", rsp)
}
