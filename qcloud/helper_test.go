package qcloud

import (
	"testing"
)

func TestQcloudEngine_DescribeCbsStorages(t *testing.T) {
	engine := NewQcloudEngine()
	param := DescribeCbsStiragesParam{
		DiskType: &DISKTYPE_DATA,
	}
	err := engine.DescribeCbsStorages(param)
	if err != nil {
		t.Error(err)
		return
	}
}
