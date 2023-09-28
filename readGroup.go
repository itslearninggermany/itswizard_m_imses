package itswizard_m_imses

import (
	"encoding/xml"
	"errors"
	"github.com/itslearninggermany/itswizard_m_basic"
	"net/http"
	"strings"
)

type ReadGroupRequest struct {
	XMLName   xml.Name `xml:"ims4:readGroupRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims4:sourcedId"`
}

type ReadGroupResponse struct {
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
			StatusInfo        struct {
				Text         string `xml:",chardata"`
				CodeMajor    string `xml:"codeMajor"`
				Severity     string `xml:"severity"`
				MessageIdRef string `xml:"messageIdRef"`
			} `xml:"statusInfo"`
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
		Text              string `xml:",chardata"`
		Xsi               string `xml:"xsi,attr"`
		Xsd               string `xml:"xsd,attr"`
		ReadGroupResponse struct {
			Text  string `xml:",chardata"`
			Xmlns string `xml:"xmlns,attr"`
			Group struct {
				Text      string `xml:",chardata"`
				GroupType struct {
					Text      string `xml:",chardata"`
					Xmlns     string `xml:"xmlns,attr"`
					Scheme    string `xml:"scheme"`
					TypeValue struct {
						Text  string `xml:",chardata"`
						Type  string `xml:"type"`
						Level string `xml:"level"`
					} `xml:"typeValue"`
				} `xml:"groupType"`
				Relationship struct {
					Text     string `xml:",chardata"`
					Xmlns    string `xml:"xmlns,attr"`
					Relation string `xml:"relation"`
					SourceId struct {
						Text       string `xml:",chardata"`
						Identifier struct {
							Text  string `xml:",chardata"`
							Xmlns string `xml:"xmlns,attr"`
						} `xml:"identifier"`
					} `xml:"sourceId"`
					Label string `xml:"label"`
				} `xml:"relationship"`
				Description struct {
					Text      string `xml:",chardata"`
					Xmlns     string `xml:"xmlns,attr"`
					DescShort string `xml:"descShort"`
					DescFull  string `xml:"descFull"`
				} `xml:"description"`
				Extension struct {
					Text           string `xml:",chardata"`
					Xmlns          string `xml:"xmlns,attr"`
					ExtensionField struct {
						Text       string `xml:",chardata"`
						Xmlns      string `xml:"xmlns,attr"`
						FieldName  string `xml:"fieldName"`
						FieldType  string `xml:"fieldType"`
						FieldValue string `xml:"fieldValue"`
					} `xml:"extensionField"`
				} `xml:"extension"`
			} `xml:"group"`
		} `xml:"readGroupResponse"`
	} `xml:"Body"`
}

type ReadGroupOutput struct {
	Group     itswizard_m_basic.Group
	Status    string
	Header    http.Header
	MessageID string
	Err       error
}

func (p *Request) ReadGroup(syncKey string) ReadGroupOutput {
	a := ReadGroupRequest{}

	a.SourcedId.Identifier = syncKey

	p.body.Body.Request = a
	outputCall := p.call("http://www.imsglobal.org/soap/gms/readGroup", nil)

	var output ReadGroupOutput

	if strings.Contains(string(outputCall.resBody), "unknownobject") {
		output.Err = errors.New("unknownobject")
		return output
	}

	if outputCall.err != nil {
		output.Err = outputCall.err
		return output
	}

	output.Status = outputCall.status
	output.Header = outputCall.header
	output.MessageID = outputCall.messageID

	var response ReadGroupResponse
	err := xml.Unmarshal(outputCall.resBody, &response)

	if err != nil {
		output.Err = errors.New("Error Message: " + outputCall.err.Error() + " Html-Body: " + string(outputCall.resBody))
		return output
	}

	course := false
	if response.Body.ReadGroupResponse.Group.GroupType.TypeValue.Type == "Course" {
		course = true
	}

	n := itswizard_m_basic.Group{
		GroupSyncKey:  a.SourcedId.Identifier,
		Name:          response.Body.ReadGroupResponse.Group.Description.DescShort,
		ParentGroupID: response.Body.ReadGroupResponse.Group.Relationship.SourceId.Identifier.Text,
		IsCourse:      course,
	}

	output.Group = n
	return output
}
