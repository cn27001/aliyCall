package thirdParty

import (
	"fmt"
	"io/ioutil"

	"github.com/gogap/ali_mns"
)

type dybaseAuth struct {
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
	MessageType     string
	QueueName       string
}
type voiceReportMsg struct {
	call_id     string // 呼叫ID    100001616500^100001871490   必须
	start_time  string // 通话开始时间,未接通则为空   2017-06-01 10:00:00 可选
	end_time    string // 通话结束时间，未接通则为空   2017-06-01 10:00:00 可选
	duration    string // 通话时长，未接通为0  10  可选
	status_code string // 呼叫结果状态码 200010  必须
	status_msg  string // 结果描述    执行完成    可选
	out_id      string // 扩展字段回传，将调用api时传入的字段返回   123456  可选
	dtmf        string // DTMF按键  123456  可选
}

const (
	ENDPOINT = "http://1943695596114318.mns.cn-hangzhou.aliyuncs.com"
	REGIONID = "cn-hangzhou"
	MSGTYPE  = "VoiceReport" //短信回执：SmsReport，短息上行：SmsUp，语音呼叫：VoiceReport，流量直冲：FlowReport

)

var Delete bool

func NewAliyunBaseApi() {

	da := newDybaseAuth()

	client := ali_mns.NewAliMNSClient(da.Endpoint, da.AccessKeyId, da.AccessKeySecret)

	resource := fmt.Sprintf("queues/%s/%s?waitseconds=10", da.QueueName, "messages")

	fmt.Println(resource)

	resp, err := client.Send("GET", nil, nil, resource)

	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

}

func newDybaseAuth() dybaseAuth {
	da := dybaseAuth{}
	da.Endpoint = ENDPOINT
	da.MessageType = MSGTYPE
	da.QueueName = au.QueueName
	da.AccessKeyId = au.accessKeyId
	da.AccessKeySecret = au.accessSecret

	return da
}
