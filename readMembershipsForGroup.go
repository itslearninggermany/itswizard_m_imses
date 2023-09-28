package itswizard_m_imses

import (
	"encoding/xml"
	"errors"
	"strings"
)

type ReadMembershipsForGroupRequest struct {
	XMLName        xml.Name `xml:"ims6:readMembershipsForGroupRequest"`
	GroupSourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims6:groupSourcedId"`
}

type ReadMembershipsForGroupResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	S       string   `xml:"s,attr"`
	U       string   `xml:"u,attr"`
	Header  struct {
		Text                   string `xml:",chardata"`
		SyncResponseHeaderInfo struct {
			Text              string `xml:",chardata"`
			H                 string `xml:"h,attr"`
			Xmlns             string `xml:"xmlns,attr"`
			Xsi               string `xml:"xsi,attr"`
			Xsd               string `xml:"xsd,attr"`
			MessageIdentifier string `xml:"messageIdentifier"`
			StatusInfoSet     struct {
				Text       string `xml:",chardata"`
				StatusInfo struct {
					Text         string `xml:",chardata"`
					CodeMajor    string `xml:"codeMajor"`
					Severity     string `xml:"severity"`
					MessageIdRef string `xml:"messageIdRef"`
				} `xml:"statusInfo"`
			} `xml:"statusInfoSet"`
		} `xml:"syncResponseHeaderInfo"`
		Security struct {
			Text           string `xml:",chardata"`
			MustUnderstand string `xml:"mustUnderstand,attr"`
			O              string `xml:"o,attr"`
			Timestamp      struct {
				Text    string `xml:",chardata"`
				ID      string `xml:"Id,attr"`
				Created string `xml:"Created"`
				Expires string `xml:"Expires"`
			} `xml:"Timestamp"`
		} `xml:"Security"`
	} `xml:"Header"`
	Body struct {
		Text                            string `xml:",chardata"`
		Xsi                             string `xml:"xsi,attr"`
		Xsd                             string `xml:"xsd,attr"`
		ReadMembershipsForGroupResponse struct {
			Text                string `xml:",chardata"`
			Xmlns               string `xml:"xmlns,attr"`
			MembershipIDPairSet struct {
				Text             string `xml:",chardata"`
				MembershipIdPair []struct {
					Text      string `xml:",chardata"`
					SourcedId struct {
						Text       string `xml:",chardata"`
						Identifier struct {
							Text  string `xml:",chardata"`
							Xmlns string `xml:"xmlns,attr"`
						} `xml:"identifier"`
					} `xml:"sourcedId"`
					Membership struct {
						Text           string `xml:",chardata"`
						GroupSourcedId struct {
							Text       string `xml:",chardata"`
							Xmlns      string `xml:"xmlns,attr"`
							Identifier struct {
								Text  string `xml:",chardata"`
								Xmlns string `xml:"xmlns,attr"`
							} `xml:"identifier"`
						} `xml:"groupSourcedId"`
						Member struct {
							Text            string `xml:",chardata"`
							Xmlns           string `xml:"xmlns,attr"`
							MemberSourcedId struct {
								Text       string `xml:",chardata"`
								Identifier struct {
									Text  string `xml:",chardata"`
									Xmlns string `xml:"xmlns,attr"`
								} `xml:"identifier"`
							} `xml:"memberSourcedId"`
							Role struct {
								Text     string `xml:",chardata"`
								RoleType string `xml:"roleType"`
							} `xml:"role"`
						} `xml:"member"`
					} `xml:"membership"`
				} `xml:"membershipIdPair"`
			} `xml:"membershipIDPairSet"`
		} `xml:"readMembershipsForGroupResponse"`
	} `xml:"Body"`
}

func (p *Request) ReadMembershipsForGroup(groupID string) (out []Membership, err error, resp string) {
	a := ReadMembershipsForGroupRequest{}

	a.GroupSourcedId.Identifier = groupID

	p.body.Body.Request = a

	resp = p.call("http://www.imsglobal.org/soap/mms/readMembershipsForGroup", nil).GetBody()

	var response ReadMembershipsForGroupResponse

	err = xml.Unmarshal([]byte(resp), &response)
	if err != nil {
		return
	}

	for i := 0; i < len(response.Body.ReadMembershipsForGroupResponse.MembershipIDPairSet.MembershipIdPair); i++ {
		out = append(out, Membership{
			ID:       response.Body.ReadMembershipsForGroupResponse.MembershipIDPairSet.MembershipIdPair[i].SourcedId.Identifier.Text,
			GroupID:  response.Body.ReadMembershipsForGroupResponse.MembershipIDPairSet.MembershipIdPair[i].Membership.GroupSourcedId.Identifier.Text,
			PersonID: response.Body.ReadMembershipsForGroupResponse.MembershipIDPairSet.MembershipIdPair[i].Membership.Member.MemberSourcedId.Identifier.Text,
			Profile:  response.Body.ReadMembershipsForGroupResponse.MembershipIDPairSet.MembershipIdPair[i].Membership.Member.Role.RoleType,
		})
	}

	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}

	return
}
