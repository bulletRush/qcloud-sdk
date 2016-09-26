package qcloud

/*
const (
	CBS_API_URI = "cbs.api.qcloud.com"
)
var (
	DISKTYPE_ROOT = "root"
	DISKTYPE_DATA = "data"
)

type DescribeCbsStoragesRequest struct {
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

type CbsStorage struct {
	Attached int `json:"attached"`
	CreateTime string `json:"createTime"`
	DeadlineTime string `json:"deadlineTime"`
	DiskType string `json:"diskType"`
	PayMode string `json:"payMode"`
	Portable int `json:"portable"`
	ProjectId int `json:"projectId"`
	SnapshotAbility int `json:"snapshotAbility"`
	StorageId string `json:"storageId"`
	StorageName string `json:"storageName"`
	StorageSize int `json:"storageSize"`
	StorageStatus string `json:"storageStatus"`
	StorageType string `json:"storageType"`
	UInstanceId string `json:"uInstanceId"`
	ZoneId int `json:"zoneId"`
}

type DescribeCbsStoragesResponse struct {
	Code int
	message string
	TotalCount int
	StorageSet []CbsStorage
}

func checkEnum(v interface{}, a ...interface{}) bool {
	for _, k := range a {
		if reflect.DeepEqual(k, v) {
			return true
		}
	}
	return false
}

func (this *QcloudEngine) DescribeCbsStorages(param DescribeCbsStoragesRequest) error {
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
	rspObj := DescribeCbsStoragesResponse{}
	err := this.DoRequest(CBS_API_URI, "DescribeCbsStorages", args, &rspObj)
	if err != nil {
		return err
	}
	this.logger.Debug("json decode succ", "rspObj", rspObj)
	return nil
}
*/
