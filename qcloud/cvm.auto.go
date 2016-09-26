package qcloud

type CvmService interface {
	DescribeInstances(status *int, zoneId *int) (*DescribeInstancesResponse, error)
}

type cvmService struct {
	QcloudEngine
	host string
	logger Logger
}

func NewCvmService(engine QcloudEngine) CvmService {
	return &cvmService{
		QcloudEngine: engine,
		host: "cvm.api.qcloud.com",
		logger: engine.GetLogger(),
	}
}


type DescribeInstancesResponse struct{
	Code int `json:"code,omitempty"`
	CodeDesc string `json:"codeDesc,omitempty"`
	Message string `json:"message,omitempty"`
}

// DescribeInstancess 查看实例列表
// 本接口 (DescribeInstances) 用于获取一个或多个实例的详细信息。
//	status: 过滤条件：实例的状态
//	zoneId: 过滤条件：可用区ID
func (svc *cvmService) DescribeInstances(status *int, zoneId *int) (*DescribeInstancesResponse, error) {
	paramMap := map[string]interface{}{}
	if status != nil {
		paramMap["status"] =  *status
	}
	if zoneId != nil {
		paramMap["zoneId"] =  *zoneId
	}
	rspObj := &DescribeInstancesResponse{}
	err := svc.DoRequest(svc.host,  "DescribeInstances", paramMap, rspObj)
	if err != nil {
		svc.logger.Error("DescribeInstances failed!", "error", err)
		return nil, err
	}
	return rspObj, nil
}

