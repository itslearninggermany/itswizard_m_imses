package itswizard_m_imses

import "encoding/xml"

type DeletePersonsRequest struct {
	XMLName      xml.Name `xml:"ims:deletePersonsRequest"`
	SourcedIdSet struct {
		Identifier []string `xml:"ims2:identifier"`
	} `xml:"ims:sourcedIdSet"`
}

func (p *Request) DeletePersons(syncKeys []string) Output {
	dpr := DeletePersonsRequest{}
	dpr.SourcedIdSet.Identifier = syncKeys
	p.body.Body.Request = dpr
	return p.call("http://www.imsglobal.org/soap/pms/deletePersons", syncKeys).GetOutput()
}
