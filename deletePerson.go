package itswizard_m_imses

import (
	"encoding/xml"
	"errors"
	"strings"
)

type DeletePersonRequest struct {
	XMLName   xml.Name `xml:"ims:deletePersonRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims:sourcedId"`
}

func (p *Request) DeletePerson(syncID string) (resp string, err error) {
	d := DeletePersonRequest{}
	d.SourcedId.Identifier = syncID
	p.body.Body.Request = d
	resp = p.call("http://www.imsglobal.org/soap/pms/deletePerson", nil).GetBody()

	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}

}
