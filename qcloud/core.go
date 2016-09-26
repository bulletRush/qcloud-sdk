package qcloud

import (
	"os"
	"time"
	"math/rand"
	"sort"
	"fmt"
	"strings"
	"encoding/base64"
	"net/http"
	"crypto/hmac"
	"crypto/sha1"
	"net/url"
	"io/ioutil"
	"encoding/json"
)

const (
	QCLOUD_SECRETID = "QCLOUD_SECRETID"
	QCLOUD_SECRETKEY = "QCLOUD_SECRETKEY"
)

type QcloudEngine interface {
	WithRegion(region string) QcloudEngine
	WithSecret(secretId, secretKey string) QcloudEngine
	WithLogger(logger Logger) QcloudEngine
	DoRequest(componetUrl, action string, content map[string]interface{}, rspObj interface{}) error
	GetLogger() Logger
}

type qcloudEngine struct {
	region string
	secretId string
	secretKey string
	urlPostfix string
	logger Logger
}

func NewQcloudEngine() QcloudEngine {
	return &qcloudEngine{
		secretKey: os.Getenv(QCLOUD_SECRETKEY),
		secretId: os.Getenv(QCLOUD_SECRETID),
		urlPostfix: "/v2/index.php",
		region: "gz",
		logger: NewLogger(os.Stdout, stdoutFormatter{}),
	}
}

func (this *qcloudEngine) WithSecret(secretId, secretKey string) QcloudEngine {
	this.secretId = secretId
	this.secretKey = secretKey
	return this
}

func (this *qcloudEngine) WithRegion(region string) QcloudEngine {
	this.region = region
	return this
}

func (this *qcloudEngine) WithLogger(logger Logger) QcloudEngine {
	this.logger = logger
	return this
}

func (this *qcloudEngine) GetLogger() Logger {
	return this.logger
}

func (this *qcloudEngine) generateGetParams(content map[string]interface{}) string {
	var keys []string
	for k := range content {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	secs := make([]string, 0)
	for _, k := range keys {
		secs = append(secs, fmt.Sprintf("%s=%v",
			url.QueryEscape(k), url.QueryEscape(fmt.Sprint(content[k])),
		))
	}
	return strings.Join(secs, "&")
}

func (this *qcloudEngine) generateSignature(
	componentUrl string, content map[string]interface{},
) string {
	raw_args := fmt.Sprintf("GET%s%s?%s", componentUrl, this.urlPostfix, this.generateGetParams(content))
	mac := hmac.New(sha1.New, []byte(this.secretKey))
	mac.Write([]byte(raw_args))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signature
}

func (this *qcloudEngine) DoRequest(
	componetUrl string, action string,
	content map[string]interface{}, rspObj interface{},
) error {
	content["Action"] = action
	content["Nonce"] = rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
	content["Region"] = this.region
	content["SecretId"] = this.secretId
	content["Timestamp"] = time.Now().UTC().Unix()
	content["Signature"] = this.generateSignature(
		componetUrl, content,
	)
	url := fmt.Sprintf("https://%s%s?%s",
		componetUrl, this.urlPostfix, this.generateGetParams(content),
	)
	rsp, err := http.Get(url)
	if err != nil {
		this.logger.Error("http get failed!", "url", url, "error", err)
		return err
	}
	defer rsp.Body.Close()
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		this.logger.Error("http read failed!", "url", url, "error", err)
		return err
	}
	this.logger.Debug("http get succ", "response", string(data), "url", url)
	if rspObj != nil {
		err = json.Unmarshal(data, rspObj)
		if err != nil {
			this.logger.Error("unmarshal json failed!", "error", err, "response", string(data))
			return err
		}
	}
	return nil
}
