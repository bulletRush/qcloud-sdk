package generator

import (
	"testing"
	"fmt"
)

func TestLoadServiceDefinitionFromFile(t *testing.T) {
	svc, err := LoadServiceDefinitionFromFile("qcloud.json")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%#v\n", svc)
}
