package thirdParty

import (
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type AuthInfo struct {
	vendor           string
	accessKeyId      string
	accessSecret     string
	CalledShowNumber string
	VoiceCode        string
	timestamp        int
	QueueName        string
}

func (p AuthInfo) GetKeyId() string {
	return p.accessKeyId
}
func (p *AuthInfo) SetKeyId(id string) {
	p.accessKeyId = id
}

func (p *AuthInfo) SetSecret(sec string) {
	p.accessSecret = sec
}

func (p AuthInfo) GetToken() string {
	return p.accessSecret
}

func (p *AuthInfo) GetTimeSec() string {

	timestamp, _ := strconv.Atoi(time.Now().Format("20060102150405"))
	if timestamp-p.timestamp > 300 {
		p.timestamp = timestamp
	}
	return strconv.Itoa(p.timestamp)
}

func LoadJsonConf(au *AuthInfo) error {

	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/aliyCall/")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	} else {
		au.accessKeyId = viper.GetString("accessKeyId")
		au.accessSecret = viper.GetString("accessSecret")
		au.CalledShowNumber = viper.GetString("CalledShowNumber")
		au.VoiceCode = viper.GetString("VoiceCode")
		au.QueueName = viper.GetString("QueueName")
	}
	return nil
}
