package itswizard_m_imses

import (
	"encoding/xml"
	"errors"
	"strings"
)

type DeleteGroupRequest struct {
	XMLName   xml.Name `xml:"ims4:deleteGroupRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims4:sourcedId"`
}

func (p *Request) DeleteGroup(syncID string) (resp string, err error) {
	d := DeleteGroupRequest{}
	d.SourcedId.Identifier = syncID
	p.body.Body.Request = d

	resp = p.call("http://www.imsglobal.org/soap/gms/deleteGroup", nil).GetBody()

	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}

}
