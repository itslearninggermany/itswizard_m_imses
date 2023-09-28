package itswizard_m_imses

import (
	"encoding/xml"
	"errors"
	"github.com/itslearninggermany/itswizard_m_basic"
	"net/http"
)

type ReadPersonRequest struct {
	XMLName   xml.Name `xml:"ims:readPersonRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims:sourcedId"`
}

type ReadPersonResponse struct {
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
			StatusInfo        struct {
				Text         string `xml:",chardata"`
				CodeMajor    string `xml:"codeMajor"`
				Severity     string `xml:"severity"`
				MessageIdRef string `xml:"messageIdRef"`
			} `xml:"statusInfo"`
		} `xml:"syncResponseHeaderInfo"`
	} `xml:"Header"`
	Body struct {
		Text               string `xml:",chardata"`
		Xsi                string `xml:"xsi,attr"`
		Xsd                string `xml:"xsd,attr"`
		ReadPersonResponse struct {
			Text   string `xml:",chardata"`
			Xmlns  string `xml:"xmlns,attr"`
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
		} `xml:"readPersonResponse"`
	} `xml:"Body"`
}

type ReadPersonOutput struct {
	Person    itswizard_m_basic.DbPerson15
	Status    string
	Header    http.Header
	MessageID string
	Err       error
}

func (p *Request) ReadPerson(syncKey string) ReadPersonOutput {
	a := ReadPersonRequest{}

	a.SourcedId.Identifier = syncKey

	p.body.Body.Request = a
	outputCall := p.call("http://www.imsglobal.org/soap/pms/readPerson", nil)

	var output ReadPersonOutput

	if outputCall.err != nil {
		output.Err = outputCall.err
		return output
	}

	output.Status = outputCall.status
	output.Header = outputCall.header
	output.MessageID = outputCall.messageID

	var response ReadPersonResponse
	err := xml.Unmarshal(outputCall.resBody, &response)

	if err != nil {
		output.Err = errors.New("Error Message: " + outputCall.err.Error() + " Html-Body: " + string(outputCall.resBody))
		return output
	}

	n := itswizard_m_basic.DbPerson15{
		Username: response.Body.ReadPersonResponse.Person.UserId.UserIdValue.Text,
		Profile:  response.Body.ReadPersonResponse.Person.InstitutionRole.InstitutionRoleType,
		Email:    response.Body.ReadPersonResponse.Person.Email.Text,
		Street1:  response.Body.ReadPersonResponse.Person.Address.Street,
		Postcode: response.Body.ReadPersonResponse.Person.Address.Postcode,
		City:     response.Body.ReadPersonResponse.Person.Address.Locality,
	}
	//First and LastName
	for in := 0; in < len(response.Body.ReadPersonResponse.Person.Name.PartName); in++ {
		if response.Body.ReadPersonResponse.Person.Name.PartName[in].NamePartType == "First" {
			n.FirstName = response.Body.ReadPersonResponse.Person.Name.PartName[in].NamePartValue
		}
		if response.Body.ReadPersonResponse.Person.Name.PartName[in].NamePartType == "Last" {
			n.LastName = response.Body.ReadPersonResponse.Person.Name.PartName[in].NamePartValue
		}
	}
	// Phone
	for in := 0; in < len(response.Body.ReadPersonResponse.Person.Tel); in++ {
		if response.Body.ReadPersonResponse.Person.Tel[in].TelType == "Voice" {
			n.Phone = response.Body.ReadPersonResponse.Person.Tel[in].TelValue
		}
		if response.Body.ReadPersonResponse.Person.Name.PartName[in].NamePartType == "Mobile" {
			n.Mobile = response.Body.ReadPersonResponse.Person.Tel[in].TelValue
		}
	}

	n.SyncPersonKey = syncKey
	output.Person = n
	return output
}
