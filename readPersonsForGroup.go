package itswizard_m_imses

import (
	"encoding/xml"
	"github.com/itslearninggermany/itswizard_m_basic"
	"net/http"
)

type ReadPersonsForGroupRequest struct {
	XMLName   xml.Name `xml:"ims:readPersonsForGroupRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims:groupSourcedId"`
}

type ReadPersonsForGroupOutput struct {
	Persons   []itswizard_m_basic.Person
	Status    string
	Header    http.Header
	MessageID string
	Err       error
}

type ReadPersonsForGroupOutputOld struct {
	Persons   []itswizard_m_basic.DbPerson15
	Status    string
	Header    http.Header
	MessageID string
	Err       error
}

type ReadPersonsForGroupResponseXML struct {
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
					} `xml:"person"`
				} `xml:"personIdPair"`
			} `xml:"personIdPairSet"`
		} `xml:"readPersonsForGroupResponse"`
	} `xml:"Body"`
}

/*
	Todo: Datum bearbeiten UCS Format

	Format of Date: "2006-01-02T15:04:05Z"
*/
//func (p * Request) ReadAllPersons (PageIndex int, PageSize int, CreatedFrom string, OnlyManuallyCreatedUsers bool, ConvertFromManual bool) (Persons []itswizard_basic.DbPerson15, Status string, Header http.Header,MessageID string, Err error)   {
func (p *Request) ReadPersonsForGroupOld(groupSyncID string) *ReadPersonsForGroupOutputOld {
	output := new(ReadPersonsForGroupOutputOld)
	a := ReadPersonsForGroupRequest{}

	a.SourcedId.Identifier = groupSyncID

	p.body.Body.Request = a

	callOut := p.call("http://www.imsglobal.org/soap/pms/readPersonsForGroup", nil)
	if callOut.err != nil {
		output.Err = callOut.err
		return output
	}
	//	fmt.Println(string(callOut.resBody))
	var response ReadPersonsForGroupResponseXML
	err := xml.Unmarshal(callOut.resBody, &response)
	if err != nil {
		output.Err = err
		return output
	}

	var persons []itswizard_m_basic.DbPerson15
	for i := 0; i < len(response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair); i++ {
		n := itswizard_m_basic.DbPerson15{
			ID:            "",
			SyncPersonKey: response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].SourcedId.Identifier.Text,
			Username:      response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.UserId.UserIdValue.Text,
			Profile:       response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.InstitutionRole.InstitutionRoleType,
			Email:         response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Email.Text,
			Street1:       response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Address.Street,
			Postcode:      response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Address.Postcode,
			City:          response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Address.Locality,
		}
		//First and LastName
		for in := 0; in < len(response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Name.PartName); in++ {
			if response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Name.PartName[in].NamePartType == "First" {
				n.FirstName = response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Name.PartName[in].NamePartValue
			}
			if response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Name.PartName[in].NamePartType == "Last" {
				n.LastName = response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Name.PartName[in].NamePartValue
			}
		}
		// Phone
		for in := 0; in < len(response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Tel); in++ {
			if response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Tel[in].TelType == "Voice" {
				n.Phone = response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Tel[in].TelValue
			}
			if response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Name.PartName[in].NamePartType == "Mobile" {
				n.Mobile = response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Tel[in].TelValue
			}
		}
		persons = append(persons, n)
	}

	output.Persons = persons
	output.MessageID = callOut.messageID
	output.Status = callOut.status
	output.Header = callOut.header
	output.Err = err
	return output
}

// func (p * Request) ReadAllPersons (PageIndex int, PageSize int, CreatedFrom string, OnlyManuallyCreatedUsers bool, ConvertFromManual bool) (Persons []itswizard_basic.DbPerson15, Status string, Header http.Header,MessageID string, Err error)   {
func (p *Request) ReadPersonsForGroup(groupSyncID string, organisationId, institutionId uint) *ReadPersonsForGroupOutput {
	output := new(ReadPersonsForGroupOutput)
	a := ReadPersonsForGroupRequest{}

	a.SourcedId.Identifier = groupSyncID

	p.body.Body.Request = a

	callOut := p.call("http://www.imsglobal.org/soap/pms/readPersonsForGroup", nil)
	if callOut.err != nil {
		output.Err = callOut.err
		return output
	}
	//ioutil.WriteFile("resp.xml",callOut.resBody,777)
	var response ReadPersonsForGroupResponseXML
	err := xml.Unmarshal(callOut.resBody, &response)
	if err != nil {
		output.Err = err
		return output
	}

	var persons []itswizard_m_basic.Person
	for i := 0; i < len(response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair); i++ {
		n := itswizard_m_basic.Person{
			PersonSyncKey:  response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].SourcedId.Identifier.Text,
			Username:       response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.UserId.UserIdValue.Text,
			Profile:        response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.InstitutionRole.InstitutionRoleType,
			Email:          response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Email.Text,
			Street1:        response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Address.Street,
			Street2:        response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Address.Extadd,
			Postcode:       response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Address.Postcode,
			City:           response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Address.Locality,
			Organisation15: organisationId,
			Institution15:  institutionId,
		}

		//First and LastName
		for in := 0; in < len(response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Name.PartName); in++ {
			if response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Name.PartName[in].NamePartType == "First" {
				n.FirstName = response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Name.PartName[in].NamePartValue
			}
			if response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Name.PartName[in].NamePartType == "Last" {
				n.LastName = response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Name.PartName[in].NamePartValue
			}
		}
		// Phone
		for in := 0; in < len(response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Tel); in++ {
			if response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Tel[in].TelType == "Voice" {
				n.Phone = response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Tel[in].TelValue
			}
			if response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Name.PartName[in].NamePartType == "Mobile" {
				n.Mobile = response.Body.ReadPersonsForGroupResponse.PersonIdPairSet.PersonIdPair[i].Person.Tel[in].TelValue
			}
		}
		persons = append(persons, n)
	}

	output.Persons = persons
	output.MessageID = callOut.messageID
	output.Status = callOut.status
	output.Header = callOut.header
	output.Err = err
	return output
}
