package itswizard_m_imses

import (
	"encoding/xml"
)

type ChildParentOut struct {
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
		Text                        string `xml:",chardata"`
		Xsi                         string `xml:"xsi,attr"`
		Xsd                         string `xml:"xsd,attr"`
		ReadPersonsForGroupResponse struct {
			Text            string `xml:",chardata"`
			Xmlns           string `xml:"xmlns,attr"`
			PersonIdPairSet struct {
				Text         string `xml:",chardata"`
				PersonIdPair []struct {
					Text      string `xml:",chardata"`
					SourcedId struct {
						Text       string `xml:",chardata"`
						Identifier struct {
							Text  string `xml:",chardata"`
							Xmlns string `xml:"xmlns,attr"`
						} `xml:"identifier"`
					} `xml:"sourcedId"`
					Person struct {
						Text       string `xml:",chardata"`
						FormatName struct {
							Text  string `xml:",chardata"`
							Nil   string `xml:"nil,attr"`
							Xmlns string `xml:"xmlns,attr"`
						} `xml:"formatName"`
						Name struct {
							Text     string `xml:",chardata"`
							Xmlns    string `xml:"xmlns,attr"`
							PartName []struct {
								Text          string `xml:",chardata"`
								NamePartType  string `xml:"namePartType"`
								NamePartValue string `xml:"namePartValue"`
							} `xml:"partName"`
						} `xml:"name"`
						Email struct {
							Text  string `xml:",chardata"`
							Xmlns string `xml:"xmlns,attr"`
						} `xml:"email"`
						UserId struct {
							Text        string `xml:",chardata"`
							Xmlns       string `xml:"xmlns,attr"`
							UserIdValue struct {
								Text  string `xml:",chardata"`
								Xmlns string `xml:"xmlns,attr"`
							} `xml:"userIdValue"`
						} `xml:"userId"`
						Address struct {
							Text     string `xml:",chardata"`
							Xmlns    string `xml:"xmlns,attr"`
							Extadd   string `xml:"extadd"`
							Locality string `xml:"locality"`
							Postcode string `xml:"postcode"`
							Street   string `xml:"street"`
						} `xml:"address"`
						Demographics struct {
							Text   string `xml:",chardata"`
							Xmlns  string `xml:"xmlns,attr"`
							Gender string `xml:"gender"`
							Bday   string `xml:"bday"`
						} `xml:"demographics"`
						InstitutionRole struct {
							Text                string `xml:",chardata"`
							Xmlns               string `xml:"xmlns,attr"`
							InstitutionRoleType string `xml:"institutionRoleType"`
							PrimaryRoleType     string `xml:"primaryRoleType"`
						} `xml:"institutionRole"`
						Tel []struct {
							Text     string `xml:",chardata"`
							Xmlns    string `xml:"xmlns,attr"`
							TelType  string `xml:"telType"`
							TelValue string `xml:"telValue"`
						} `xml:"tel"`
						Extension struct {
							Text         string `xml:",chardata"`
							Xmlns        string `xml:"xmlns,attr"`
							Relationship []struct {
								Text     string `xml:",chardata"`
								Xmlns    string `xml:"xmlns,attr"`
								Relation string `xml:"relation"`
								SourceId struct {
									Text       string `xml:",chardata"`
									Identifier string `xml:"identifier"`
								} `xml:"sourceId"`
							} `xml:"relationship"`
						} `xml:"extension"`
					} `xml:"person"`
				} `xml:"personIdPair"`
			} `xml:"personIdPairSet"`
		} `xml:"readPersonsForGroupResponse"`
	} `xml:"Body"`
}

func (p *Request) ReadParenChildRelationship(groupSyncID string, organisationId, institutionId uint) (map[string][]string, error) {
	a := ReadPersonsForGroupRequest{}

	a.SourcedId.Identifier = groupSyncID

	p.body.Body.Request = a

	callOut := p.call("http://www.imsglobal.org/soap/pms/readPersonsForGroup", nil)

	if callOut.err != nil {
		return nil, callOut.err
	}
	var o ChildParentOut
	err := xml.Unmarshal(callOut.resBody, &o)
	if err != nil {
		return nil, err
	}

	psr := make(map[string][]string)

	for i := 0; i < len(o.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair); i++ {
		var tmp []string
		for s := 0; s < (len(o.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Extension.Relationship)); s++ {
			if o.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Extension.Relationship[s].Relation == "Child" {
				tmp = append(tmp, o.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Extension.Relationship[s].SourceId.Identifier)
			}
		}
		if len(tmp) > 0 {
			psr[o.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].SourcedId.Identifier.Text] = tmp
		}
	}
	return psr, nil
}
