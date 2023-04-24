package main

import "encoding/xml"

func DecodeMultiStatus(data []byte) (*MultiStatus, error) {
	var input multiStatusDecode
	if err := xml.Unmarshal(data, &input); err != nil {
		return nil, err
	}

	var output MultiStatus
	output.NS = "DAV:"
	output.Responses = make([]StatusResponse, len(input.Responses))
	for i, r := range input.Responses {
		output.Responses[i].URI = r.URI
		output.Responses[i].PropStat.Prop.GetEtag = r.PropStat.Prop.GetEtag
		output.Responses[i].PropStat.Status = r.PropStat.Status
	}

	return &output, nil
}

type MultiStatus struct {
	XMLName   xml.Name         `xml:"ns0:multistatus"`
	NS        string           `xml:"xmlns:ns0,attr"`
	Responses []StatusResponse `xml:"ns0:response"`
}

type StatusResponse struct {
	XMLName  xml.Name `xml:"ns0:response"`
	URI      string   `xml:"ns0:href"`
	PropStat PropStat `xml:"ns0:propstat"`
}

type PropStat struct {
	XMLName xml.Name `xml:"ns0:propstat"`
	Prop    Prop     `xml:"ns0:prop"`
	Status  string   `xml:"ns0:status"`
}

type Prop struct {
	XMLName xml.Name `xml:"ns0:prop"`
	GetEtag string   `xml:"ns0:getetag"`
}

type multiStatusDecode struct {
	XMLName   xml.Name               `xml:"multistatus"`
	Responses []statusResponseDecode `xml:"response"`
}

type statusResponseDecode struct {
	XMLName  xml.Name       `xml:"response"`
	URI      string         `xml:"href"`
	PropStat propStatDecode `xml:"propstat"`
}

type propStatDecode struct {
	XMLName xml.Name   `xml:"propstat"`
	Prop    propDecode `xml:"prop"`
	Status  string     `xml:"status"`
}

type propDecode struct {
	XMLName xml.Name `xml:"prop"`
	GetEtag string   `xml:"getetag"`
}
