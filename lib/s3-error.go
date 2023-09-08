package lib

import (
	"encoding/xml"
	"fmt"
)

type S3Error struct {
	XMLName   xml.Name `xml:"Error"`
	Code      string   `xml:"Code"`
	Message   string   `xml:"Message"`
	HostId    string   `xml:"HostId"`
	RequestId string   `xml:"RequestId"`
}

func (s S3Error) Error() string {
	return fmt.Sprintf("%v - %v", s.Code, s.Message)
}

func (s S3Error) Empty() bool {
	return s.Message == "" && s.Code == ""
}
