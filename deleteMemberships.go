package itswizard_m_imses

import (
	"encoding/xml"
	"errors"
	"strings"
)

type DeleteMembershipRequest struct {
	XMLName   xml.Name `xml:"ims6:deleteMembershipRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims6:sourcedId"`
}

func (p *Request) DeleteMembership(syncID string) (resp string, err error) {
	d := DeleteMembershipRequest{}
	d.SourcedId.Identifier = syncID
	p.body.Body.Request = d

	resp = p.call("http://www.imsglobal.org/soap/mms/deleteMembership", nil).GetBody()

	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}

}
