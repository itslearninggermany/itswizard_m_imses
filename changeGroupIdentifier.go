package itswizard_m_imses

import (
	"encoding/xml"
	"errors"
)

type ChangeGroupsIdentifierRequest struct {
	XMLName          xml.Name `xml:"ims4:changeGroupsIdentifierRequest"`
	PairSourcedIdSet struct {
		IdentifierPair struct {
			FirstId  string `xml:"ims2:firstId"`
			SecondId string `xml:"ims2:secondId"`
		} `xml:"ims2:identifierPair"`
	} `xml:"ims4:pairSourcedIdSet"`
}

func (p *Request) ChangeGroupsIdentifierRequest(old, new string) error {
	cpi := ChangeGroupsIdentifierRequest{}
	cpi.PairSourcedIdSet.IdentifierPair.FirstId = old
	cpi.PairSourcedIdSet.IdentifierPair.SecondId = new
	p.body.Body.Request = cpi
	var out ResponseChangeIdentGroup
	err := xml.Unmarshal(p.call("http://www.imsglobal.org/soap/gms/changeGroupsIdentifier", nil).ResponseBody(), &out)
	if err != nil {
		return err
	}
	if out.Header.SyncResponseHeaderInfo.StatusInfo.Description.Text.Text != "" {
		return errors.New(out.Header.SyncResponseHeaderInfo.StatusInfo.Description.Text.Text)
	}
	return nil
}

type ResponseChangeIdentGroup struct {
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
				Text      string `xml:",chardata"`
				CodeMajor string `xml:"codeMajor"`
				Severity  string `xml:"severity"`
				CodeMinor struct {
					Text           string `xml:",chardata"`
					CodeMinorField struct {
						Text           string `xml:",chardata"`
						CodeMinorName  string `xml:"codeMinorName"`
						CodeMinorValue string `xml:"codeMinorValue"`
					} `xml:"codeMinorField"`
				} `xml:"codeMinor"`
				MessageIdRef string `xml:"messageIdRef"`
				Description  struct {
					Chardata string `xml:",chardata"`
					Language struct {
						Text  string `xml:",chardata"`
						Xmlns string `xml:"xmlns,attr"`
					} `xml:"language"`
					Text struct {
						Text  string `xml:",chardata"`
						Xmlns string `xml:"xmlns,attr"`
					} `xml:"text"`
				} `xml:"description"`
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
		Text                           string `xml:",chardata"`
		Xsi                            string `xml:"xsi,attr"`
		Xsd                            string `xml:"xsd,attr"`
		ChangeGroupsIdentifierResponse struct {
			Text  string `xml:",chardata"`
			Xmlns string `xml:"xmlns,attr"`
		} `xml:"changePersonIdentifierResponse"`
	} `xml:"Body"`
}
