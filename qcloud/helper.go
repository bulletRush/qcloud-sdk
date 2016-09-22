package qcloud

import (
	"reflect"
	"fmt"
)

const (
	CBS_API_URI = "cbs.api.qcloud.com"
)
var (
	DISKTYPE_ROOT = "root"
	DISKTYPE_DATA = "data"
)

type DescribeCbsStiragesParam struct {
	DiskType *string
	PayMode *string
	Portable *int
	ProjectId *int
	StorageIds []string
	StorageType *string
	StorageStatus []string
	ZoneId *int
	Offset *int
	Limit *int
}

func checkEnum(v interface{}, a ...interface{}) bool {
	for _, k := range a {
		if reflect.DeepEqual(k, v) {
			return true
		}
	}
	return false
}

func (this *QcloudEngine) DescribeCbsStorages(param DescribeCbsStiragesParam) error {
	args := map[string]interface{}{}
	if param.DiskType != nil {
		if !checkEnum(*param.DiskType, DISKTYPE_DATA, DISKTYPE_ROOT) {
			this.logger.Error("diskType not in enum")
			return NewParamError("diskType not in enum")
		}
		args["diskType"] = *param.DiskType
	}
	if len(param.StorageIds) > 0 {
		for idx, storageId := range param.StorageIds {
			key := fmt.Sprintf("storageIds.%d", idx)
			args[key] = storageId
		}
	}
	err := this.DoRequest(CBS_API_URI, "DescribeCbsStorages", args)
	if err != nil {
		return err
	}
	return nil
}
