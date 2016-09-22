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
)

const (
	QCLOUD_SECRETID = "QCLOUD_SECRETID"
	QCLOUD_SECRETKEY = "QCLOUD_SECRETKEY"
)



type QcloudEngine struct {
	region string
	secretId string
	secretKey string
	urlPostfix string
	logger Logger
}

func NewQcloudEngine() *QcloudEngine {
	return &QcloudEngine{
		secretKey: os.Getenv(QCLOUD_SECRETKEY),
		secretId: os.Getenv(QCLOUD_SECRETID),
		urlPostfix: "/v2/index.php",
		region: "gz",
		logger: NewLogger(os.Stdout, stdoutFormatter{}),
	}
}

func (this *QcloudEngine) WithSecret(secretId, secretKey string) *QcloudEngine {
	this.secretId = secretId
	this.secretKey = secretKey
	return this
}

func (this *QcloudEngine) WithRegion(region string) *QcloudEngine {
	this.region = region
	return this
}

func (this *QcloudEngine) WithLogger(logger Logger) *QcloudEngine {
	this.logger = logger
	return this
}

func (this *QcloudEngine) GenerateGetParams(content map[string]interface{}) string {
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

func (this *QcloudEngine) GenerateSignature(
	componentUrl string, content map[string]interface{},
) string {
	raw_args := fmt.Sprintf("GET%s%s?%s", componentUrl, this.urlPostfix, this.GenerateGetParams(content))
	mac := hmac.New(sha1.New, []byte(this.secretKey))
	mac.Write([]byte(raw_args))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signature
}

func (this *QcloudEngine) DoRequest(
	componetUrl string, action string,
	content map[string]interface{},
) error {
	content["Action"] = action
	content["Nonce"] = rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
	content["Region"] = this.region
	content["SecretId"] = this.secretId
	content["Timestamp"] = time.Now().UTC().Unix()
	content["Signature"] = this.GenerateSignature(
		componetUrl, content,
	)
	url := fmt.Sprintf("https://%s%s?%s",
		componetUrl, this.urlPostfix, this.GenerateGetParams(content),
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
	return nil
}