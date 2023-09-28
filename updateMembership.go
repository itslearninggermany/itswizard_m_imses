package itswizard_m_imses

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

type UpdateMembershipRequest struct {
	XMLName   xml.Name `xml:"ims6:updateMembershipRequest"`
	Text      string   `xml:",chardata"`
	SourcedId struct {
		Text       string `xml:",chardata"`
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims6:sourcedId"`
	Membership struct {
		Text           string `xml:",chardata"`
		GroupSourcedId struct {
			Text       string `xml:",chardata"`
			Identifier string `xml:"ims2:identifier"`
		} `xml:"ims7:groupSourcedId"`
		Member struct {
			Text            string `xml:",chardata"`
			MemberSourcedId struct {
				Text       string `xml:",chardata"`
				Identifier string `xml:"ims2:identifier"`
			} `xml:"ims7:memberSourcedId"`
			Role struct {
				Text     string `xml:",chardata"`
				RoleType string `xml:"ims7:roleType"`
			} `xml:"ims7role"`
		} `xml:"ims7:member"`
	} `xml:"ims6:membership"`
}

func (p *Request) UpdateMembership(input Membership) {
	a := UpdateMembershipRequest{}

	a.SourcedId.Identifier = input.ID
	a.Membership.Member.MemberSourcedId.Identifier = input.PersonID
	a.Membership.GroupSourcedId.Identifier = input.GroupID
	a.Membership.Member.Role.RoleType = input.Profile

	p.body.Body.Request = a
	outputCall := p.call("http://www.imsglobal.org/soap/mms/updateMembership", nil)

	fmt.Println(string(outputCall.resBody))
	return
}

func (p *Request) UpdateMembershipNew(groupID string, personID string, profile string) (resp string, err error) {
	a := UpdateMembershipRequest{}

	a.SourcedId.Identifier = personID + groupID
	a.Membership.GroupSourcedId.Identifier = groupID
	a.Membership.Member.MemberSourcedId.Identifier = personID
	a.Membership.Member.Role.RoleType = profile

	p.body.Body.Request = a

	resp = p.call("http://www.imsglobal.org/soap/mms/updateMembership", nil).GetBody()

	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}

	return
}
