package Common

import (
	"fmt"
	"net/url"
	//	"Common/logger"
)

var (
	SMS_UUID = "80341"                         // uid
	SMS_AUTH = GetMd5String("alxx" + "alxx88") // md5(code + pwd)
)

//func init() {
//	sms := &SmsNotify{}
//	//	sms.SetPhone("18621669012") //zoulifeng
//	//	sms.SetPhone("18158139355") //lili
//	//	sms.SetPhone("15221975446") //zhouxueshi
//	logger.BackupNohup(sms)
//}

type SmsNotify struct {
	phoneList []string
}

func (s *SmsNotify) SetPhone(phone string) {
	s.phoneList = append(s.phoneList, phone)
}

func (s SmsNotify) Notify(info string) {
	for _, phone := range s.phoneList {
		SendSMS(phone, info)
	}
}

func SendSMS(phone, info string) error {
	v := url.Values{}
	v.Set("uid", SMS_UUID)
	v.Set("auth", SMS_AUTH)
	v.Set("mobile", phone)
	v.Set("expid", "0")
	v.Set("encode", "utf-8")
	v.Set("msg", info)
	req := v.Encode()

	resp, err := SendHttpReq([]byte(req), "http://sms.10690221.com:9011/hy/", Http_req_get, nil)
	if err != nil {
		fmt.Printf("SendSMS : ", err.Error())
		return err
	}
	fmt.Printf("SendSMS : phone=%s, resp=%s\n", phone, string(resp))
	return nil
}
