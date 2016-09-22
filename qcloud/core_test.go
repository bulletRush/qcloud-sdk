package qcloud

import (
	"testing"
)

func TestQcloudEngine_DoRequest(t *testing.T) {
	engine := NewQcloudEngine().WithRegion("gz")
	componentUrl := "cbs.api.qcloud.com"
	var rspObj interface{}
	err := engine.DoRequest(
		componentUrl, "DescribeCbsStorages", map[string]interface{}{}, rspObj,
	)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestQcloudEngine_GenerateSignature(t *testing.T) {
	engine := NewQcloudEngine().WithRegion("gz").WithSecret(
		"AKIDz8krbsJ5yKBZQpn74WFkmLPx3gnPhESA",
		"Gu5t9xGARNpq86cd98joQYCN3Cozk1qA",
	)
	componentUrl := "cvm.api.qcloud.com"
	content := map[string]interface{}{
		"Action": "DescribeInstances",
		"Nonce": 11886,
		"Region": "gz",
		"SecretId": "AKIDz8krbsJ5yKBZQpn74WFkmLPx3gnPhESA",
		"Timestamp": 1465185768,
		"instanceIds.0": "ins-09dx96dg",
		"limit": 20,
		"offset": 0,
	}
	signature := engine.GenerateSignature(componentUrl, content)
	correctSignature := "NSI3UqqD99b/UJb4tbG/xZpRW64="
	if signature != correctSignature {
		t.Errorf("not equal: %s != %s",
			signature, correctSignature)
		return
	}
}

