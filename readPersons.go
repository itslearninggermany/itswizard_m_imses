package itswizard_m_imses

import (
	"encoding/xml"
	"errors"
	"github.com/itslearninggermany/itswizard_m_basic"
	"net/http"
)

type ReadPersonsRequest struct {
	XMLName      xml.Name `xml:"ims:readPersonsRequest"`
	SourcedIdSet struct {
		Identifier []string `xml:"ims2:identifier"`
	} `xml:"ims:sourcedIdSet"`
}

type ReadPersonsResponse struct {
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
				StatusInfo []struct {
					Text         string `xml:",chardata"`
					CodeMajor    string `xml:"codeMajor"`
					Severity     string `xml:"severity"`
					MessageIdRef string `xml:"messageIdRef"`
				} `xml:"statusInfo"`
			} `xml:"statusInfoSet"`
		} `xml:"syncResponseHeaderInfo"`
	} `xml:"Header"`
	Body struct {
		Text                string `xml:",chardata"`
		Xsi                 string `xml:"xsi,attr"`
		Xsd                 string `xml:"xsd,attr"`
		ReadPersonsResponse struct {
			Text      string `xml:",chardata"`
			Xmlns     string `xml:"xmlns,attr"`
			PersonSet struct {
				Text   string `xml:",chardata"`
				Person []struct {
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
			} `xml:"personSet"`
		} `xml:"readPersonsResponse"`
	} `xml:"Body"`
}

type ReadPersonsOutput struct {
	Persons   []itswizard_m_basic.DbPerson15
	Status    string
	Header    http.Header
	MessageID string
	Err       error
}

func (p *Request) ReadPersons(syncKeys []string) ReadPersonsOutput {

	a := ReadPersonsRequest{}

	a.SourcedIdSet.Identifier = syncKeys

	p.body.Body.Request = a
	outputCall := p.call("http://www.imsglobal.org/soap/pms/readPersons", nil)

	var output ReadPersonsOutput

	if outputCall.err != nil {
		output.Err = outputCall.err
		return output
	}

	output.Status = outputCall.status
	output.Header = outputCall.header
	output.MessageID = outputCall.messageID

	var response ReadPersonsResponse
	err := xml.Unmarshal(outputCall.resBody, &response)

	if err != nil {
		output.Err = errors.New("Error Message: " + outputCall.err.Error() + " Html-Body: " + string(outputCall.resBody))
		return output
	}

	var persons []itswizard_m_basic.DbPerson15

	for i := 0; i < len(response.Body.ReadPersonsResponse.PersonSet.Person); i++ {

		///
		n := itswizard_m_basic.DbPerson15{
			Username: response.Body.ReadPersonsResponse.PersonSet.Person[i].UserId.UserIdValue.Text,
			Profile:  response.Body.ReadPersonsResponse.PersonSet.Person[i].InstitutionRole.InstitutionRoleType,
			Email:    response.Body.ReadPersonsResponse.PersonSet.Person[i].Email.Text,
			Street1:  response.Body.ReadPersonsResponse.PersonSet.Person[i].Address.Street,
			Postcode: response.Body.ReadPersonsResponse.PersonSet.Person[i].Address.Postcode,
			City:     response.Body.ReadPersonsResponse.PersonSet.Person[i].Address.Locality,
		}
		//First and LastName
		for in := 0; in < len(response.Body.ReadPersonsResponse.PersonSet.Person[i].Name.PartName); in++ {
			if response.Body.ReadPersonsResponse.PersonSet.Person[i].Name.PartName[in].NamePartType == "First" {
				n.FirstName = response.Body.ReadPersonsResponse.PersonSet.Person[i].Name.PartName[in].NamePartValue
			}
			if response.Body.ReadPersonsResponse.PersonSet.Person[i].Name.PartName[in].NamePartType == "Last" {
				n.LastName = response.Body.ReadPersonsResponse.PersonSet.Person[i].Name.PartName[in].NamePartValue
			}
		}
		// Phone
		for in := 0; in < len(response.Body.ReadPersonsResponse.PersonSet.Person[i].Tel); in++ {
			if response.Body.ReadPersonsResponse.PersonSet.Person[i].Tel[in].TelType == "Voice" {
				n.Phone = response.Body.ReadPersonsResponse.PersonSet.Person[i].Tel[in].TelValue
			}
			if response.Body.ReadPersonsResponse.PersonSet.Person[i].Name.PartName[in].NamePartType == "Mobile" {
				n.Mobile = response.Body.ReadPersonsResponse.PersonSet.Person[i].Tel[in].TelValue
			}
		}
		persons = append(persons, n)
	}
	output.Persons = persons
	return output
}
