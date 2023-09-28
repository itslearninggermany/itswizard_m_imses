package itswizard_m_imses

import (
	"encoding/xml"
	"fmt"
)

type ReadMembershipsForPersonRequest struct {
	XMLName         xml.Name `xml:"ims6:readMembershipsForPersonRequest"`
	Text            string   `xml:",chardata"`
	PersonSourcedId struct {
		Text       string `xml:",chardata"`
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims6:personSourcedId"`
}

type readMembershipsForPersonResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	S       string   `xml:"s,attr"`
	Header  struct {
		Text   string `xml:",chardata"`
		Action struct {
			Text           string `xml:",chardata"`
			MustUnderstand string `xml:"mustUnderstand,attr"`
			Xmlns          string `xml:"xmlns,attr"`
		} `xml:"Action"`
		SyncResponseHeaderInfo struct {
			Text              string `xml:",chardata"`
			Xmlns             string `xml:"xmlns,attr"`
			Xsi               string `xml:"xsi,attr"`
			Xsd               string `xml:"xsd,attr"`
			H                 string `xml:"h,attr"`
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
	} `xml:"Header"`
	Body struct {
		Text                             string `xml:",chardata"`
		Xsi                              string `xml:"xsi,attr"`
		Xsd                              string `xml:"xsd,attr"`
		ReadMembershipsForPersonResponse struct {
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
								Text      string `xml:",chardata"`
								RoleType  string `xml:"roleType"`
								Extension struct {
									Text           string `xml:",chardata"`
									ExtensionField struct {
										Text       string `xml:",chardata"`
										Xmlns      string `xml:"xmlns,attr"`
										FieldName  string `xml:"fieldName"`
										FieldType  string `xml:"fieldType"`
										FieldValue string `xml:"fieldValue"`
									} `xml:"extensionField"`
								} `xml:"extension"`
							} `xml:"role"`
						} `xml:"member"`
					} `xml:"membership"`
				} `xml:"membershipIdPair"`
			} `xml:"membershipIDPairSet"`
		} `xml:"readMembershipsForPersonResponse"`
	} `xml:"Body"`
}

func (p *Request) ReadMembershipsForPerson(syncKey string) (r []Membership) {
	a := ReadMembershipsForPersonRequest{}

	a.PersonSourcedId.Identifier = syncKey
	p.body.Body.Request = a
	outputCall := p.call("http://www.imsglobal.org/soap/mms/readMembershipsForPerson", nil)
	var out readMembershipsForPersonResponse
	err := xml.Unmarshal(outputCall.resBody, &out)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < len(out.Body.ReadMembershipsForPersonResponse.MembershipIDPairSet.MembershipIdPair); i++ {
		r = append(r, Membership{
			ID:       out.Body.ReadMembershipsForPersonResponse.MembershipIDPairSet.MembershipIdPair[i].SourcedId.Identifier.Text,
			GroupID:  out.Body.ReadMembershipsForPersonResponse.MembershipIDPairSet.MembershipIdPair[i].Membership.GroupSourcedId.Identifier.Text,
			PersonID: out.Body.ReadMembershipsForPersonResponse.MembershipIDPairSet.MembershipIdPair[i].Membership.Member.MemberSourcedId.Identifier.Text,
			Profile:  out.Body.ReadMembershipsForPersonResponse.MembershipIDPairSet.MembershipIdPair[i].Membership.Member.Role.RoleType,
		})
	}

	return
}
