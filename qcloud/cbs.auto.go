package qcloud

type CbsService interface {
	DescribeCbsStorages(diskType *string, payMode *string) (*DescribeCbsStoragesResponse, error)
}

type cbsService struct {
	QcloudEngine
	host string
	logger Logger
}

func NewCbsService(engine QcloudEngine) CbsService {
	return &cbsService{
		QcloudEngine: engine,
		host: "cbs.api.qcloud.com",
		logger: engine.GetLogger(),
	}
}


type DescribeCbsStoragesResponse struct{
	Code int `json:"code,omitempty"`
	CodeDesc string `json:"codeDesc,omitempty"`
	Message string `json:"message,omitempty"`
}

// DescribeCbsStoragess 查询云硬盘信息
// 本接口（DescribeCbsStorages）用于查询云硬盘的详细信息。可根据云硬盘ID、云硬盘状态，云硬盘类型等对结果进行过滤。对于过滤条件，不同条件之间为与(AND)的关系，如果不传入则不以此条件过滤。
//	diskType: 标准值：root代表系统盘，data代表数据盘
//	payMode: 付费方式。标准值为包年包月：prePay和按量计费：postPay
func (svc *cbsService) DescribeCbsStorages(diskType *string, payMode *string) (*DescribeCbsStoragesResponse, error) {
	paramMap := map[string]interface{}{}
	if diskType != nil {
		paramMap["diskType"] =  *diskType
	}
	if payMode != nil {
		paramMap["payMode"] =  *payMode
	}
	rspObj := &DescribeCbsStoragesResponse{}
	err := svc.DoRequest(svc.host,  "DescribeCbsStorages", paramMap, rspObj)
	if err != nil {
		svc.logger.Error("DescribeCbsStorages failed!", "error", err)
		return nil, err
	}
	return rspObj, nil
}

