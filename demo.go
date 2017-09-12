package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gotoolkits/aliyCall/thirdParty"
)

func main() {

	var phone string
	flag.StringVar(&phone, "p", "", "Receive phone numbers")
	flag.Parse()

	if len(phone) < 1 {
		fmt.Println("Usage: -p [PhoneNumber]       Receive phone numbers")
		os.Exit(1)
	}

	session := thirdParty.NewAliyunVmsApi(phone)
	thirdParty.PostMsgToAliyun(session)

	thirdParty.NewAliyunBaseApi()

}
