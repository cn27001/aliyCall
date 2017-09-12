package thirdParty

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gotoolkits/backend/utils"
	"github.com/gotoolkits/gorequest"
)

var aliyunApi = "http://dyvmsapi.aliyuncs.com/"
var accessSecret string
var sortedQueryString string

type AliyunVoiceSysParams struct {
	AccessKeyId      string
	Timestamp        string //格式为：yyyy-MM-dd’T’HH:mm:ss’Z’；时区为：GMT
	Format           string //没传默认为JSON，可选填值：XML
	SignatureMethod  string //建议固定值：HMAC-SHA1
	SignatureVersion string //建议固定值：1.0
	SignatureNonce   string //用于请求的防重放攻击，每次请求唯一，JAVA语言建议用：java.util.UUID.randomUUID()生成即可
	Signature        string //最终生成的签名结果值
}

type AliyunVoiceBussParams struct {
	Action               string //API的命名，固定值，如发送短信API的值为：SendSms / SingleCallByVoice
	Version              string //API的版本，固定值，如短信API的值为：2017-05-25
	RegionId             string //API支持的RegionID，如短信API的值为：cn-hangzhou
	CalledShowNumber     string //语音通知显示号码，必填
	CalledNumber         string //语音通知的被叫号码，必填。
	VoiceCode            string //语音通知的媒体文件编码,必填
	BusinessId           string //设置业务请求流水号
	OutId                string
	ResourceOwnerId      string
	OwnerId              string
	ResourceOwnerAccount string
}

var au AuthInfo

func NewAliyunVmsApi(phone string) url.Values {

	err := LoadJsonConf(&au)
	if err != nil {
		fmt.Println("load the config json failed ,please check the configuration.", err)
		return nil
	}

	accessSecret = au.accessSecret + "&"

	local, err := time.LoadLocation("GMT")
	if err != nil {
		fmt.Println(err)
	}
	t := time.Now()
	tm := t.In(local).Format("2006-01-02T15:04:05Z")

	sysParams := AliyunVoiceSysParams{
		AccessKeyId:      au.accessKeyId,
		Timestamp:        tm,
		Format:           "json",
		SignatureMethod:  "HMAC-SHA1",
		SignatureVersion: "1.0",
		SignatureNonce:   rrutils.NewV4().String(),
	}

	bussParams := AliyunVoiceBussParams{
		Action:           "SingleCallByVoice",
		Version:          "2017-05-25",
		RegionId:         "cn-hangzhou",
		CalledShowNumber: au.CalledShowNumber,
		CalledNumber:     phone,
		VoiceCode:        au.VoiceCode,
		OutId:            rrutils.NewV4().String(),
	}

	v := url.Values{}

	v.Add("AccessKeyId", sysParams.AccessKeyId)
	v.Add("Format", sysParams.Format)
	v.Add("SignatureMethod", sysParams.SignatureMethod)
	v.Add("SignatureNonce", sysParams.SignatureNonce)
	v.Add("SignatureVersion", sysParams.SignatureVersion)
	v.Add("Timestamp", sysParams.Timestamp)

	v.Add("Action", bussParams.Action)
	v.Add("Version", bussParams.Version)
	v.Add("RegionId", bussParams.RegionId)

	v.Add("CalledShowNumber", bussParams.CalledShowNumber)
	v.Add("CalledNumber", bussParams.CalledNumber)
	v.Add("VoiceCode", bussParams.VoiceCode)
	v.Add("OutId", bussParams.OutId)

	return v
}

func SignaturedAndUriArgs(v url.Values) string {
	return aliyunApi + "?" + "Signature=" + reqEncodeAndSign(v) + "&" + sortedQueryString
}

func PostMsgToAliyun(v url.Values) ([]byte, error) {

	url := SignaturedAndUriArgs(v)
	request := gorequest.New()
	resp, body, errs := request.Get(url).End()

	if errs != nil {
		fmt.Println(resp.StatusCode, errs)
	}
	fmt.Println(body)
	return nil, nil

}

func specialUrlEncode(v url.Values) string {
	urlEn := v.Encode()
	urlEn = strings.Replace(urlEn, "+", "%20", -1)
	urlEn = strings.Replace(urlEn, "*", "%2A", -1)
	urlEn = strings.Replace(urlEn, "%7E", "~", -1)
	return urlEn
}

func specialUrlEncodeToStr(v string) string {

	urlEn := url.QueryEscape(v)
	urlEn = strings.Replace(urlEn, "+", "%20", -1)
	urlEn = strings.Replace(urlEn, "*", "%2A", -1)
	urlEn = strings.Replace(urlEn, "%7E", "~", -1)
	return urlEn
}

func stringToSign(en string) string {
	var HTTPMethod = "GET"
	str := HTTPMethod + "&" + url.QueryEscape("/") + "&"
	return str + specialUrlEncodeToStr(en)
}

func sign(accessSecret, stringToSign string) string {
	mac := hmac.New(sha1.New, []byte(accessSecret))
	mac.Write([]byte(stringToSign))
	signData := mac.Sum(nil)

	return base64.StdEncoding.EncodeToString(signData)
}

func reqEncodeAndSign(v url.Values) string {
	sortedQueryString = specialUrlEncode(v)
	toSign := stringToSign(sortedQueryString)
	signature := sign(accessSecret, toSign)

	return specialUrlEncodeToStr(signature)
}
