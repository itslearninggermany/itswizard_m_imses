package itswizard_m_imses

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"strings"
)

const x = "http://schemas.xmlsoap.org/soap/envelope/"
const ims = "http://www.imsglobal.org/services/pms/xsd/imsPersonManMessSchema_v1p0"
const ims1 = "http://www.imsglobal.org/services/common/imsMessBindSchema_v1p0"
const ims2 = "http://www.imsglobal.org/services/common/imsCommonSchema_v1p0"
const ims3 = "http://www.imsglobal.org/services/pms/xsd/imsPersonManDataSchema_v1p0"
const ims4 = "http://www.imsglobal.org/services/gms/xsd/imsGroupManMessSchema_v1p0"
const ims5 = "http://www.imsglobal.org/services/gms/xsd/imsGroupManDataSchema_v1p0"
const ims6 = "http://www.imsglobal.org/services/mms/xsd/imsMemberManMessSchema_v1p0"
const ims7 = "http://www.imsglobal.org/services/mms/xsd/imsMemberManDataSchema_v1p0"

/*
information for the Request
*/
type Request struct {
	URL      string
	Username string
	Password string
	body     HTTPBody
}

/*
...
*/
type HTTPBody struct {
	XMLName xml.Name   `xml:"x:Envelope"`
	X       string     `xml:"xmlns:x,attr"`
	Ims     string     `xml:"xmlns:ims,attr"`
	Ims1    string     `xml:"xmlns:ims1,attr"`
	Ims2    string     `xml:"xmlns:ims2,attr"`
	Ims3    string     `xml:"xmlns:ims3,attr"`
	Ims4    string     `xml:"xmlns:ims4,attr"`
	Ims5    string     `xml:"xmlns:ims5,attr"`
	Ims6    string     `xml:"xmlns:ims6,attr"`
	Ims7    string     `xml:"xmlns:ims7,attr"`
	Header  soapHeader `xml:"x:Header"`
	Body    SoapBody   `xml:"x:Body"`
}

/*
...
*/
type SoapBody struct {
	XMLName xml.Name `xml:"x:Body"`
	Request interface{}
}

type NewImsesServiceInput struct {
	Username string
	Password string
	Url      string
}

/*
This is for the Header in the SOAP
*/
type soapHeader struct {
	XMLName               xml.Name `xml:"x:Header"`
	SyncRequestHeaderInfo struct {
		MessageIdentifier string `xml:"ims1:messageIdentifier"`
	} `xml:"ims1:syncRequestHeaderInfo"`
	Security struct {
		Text          string `xml:",chardata"`
		Wsse          string `xml:"xmlns:wsse,attr"`
		Wsu           string `xml:"xmlns:wsu,attr"`
		UsernameToken struct {
			Username string `xml:"wsse:Username"`
			Password struct {
				Text string `xml:",chardata"`
				Type string `xml:"Type,attr"`
			} `xml:"wsse:Password"`
		} `xml:"wsse:UsernameToken"`
	} `xml:"wsse:Security"`
}

/*
initiate the Header
*/
func (p *Request) initHeader(identifier string) {
	a := new(soapHeader)
	a.SyncRequestHeaderInfo.MessageIdentifier = identifier
	a.Security.Wsse = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"
	a.Security.Wsu = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd"
	a.Security.UsernameToken.Username = p.Username
	a.Security.UsernameToken.Password.Text = p.Password
	a.Security.UsernameToken.Password.Type = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText"
	p.body.Header = *a
}

/*
Creates a new Request
*/
func NewImsesService(input NewImsesServiceInput) *Request {
	a := new(Request)
	a.Username = input.Username
	a.Password = input.Password
	a.URL = input.Url
	a.body.X = x
	a.body.Ims = ims
	a.body.Ims1 = ims1
	a.body.Ims2 = ims2
	a.body.Ims3 = ims3
	a.body.Ims4 = ims4
	a.body.Ims5 = ims5
	a.body.Ims6 = ims6
	a.body.Ims7 = ims7
	a.initHeader(uuid.New().String())
	return a
}

//

type calloutput struct {
	status    string
	header    http.Header
	resBody   []byte
	messageID string
	err       error
	sendData  []string
}

func (p *Request) call(soapAction string, sendData []string) *calloutput {
	by, err := xml.Marshal(p.body)

	o := new(calloutput)

	o.sendData = sendData
	q, err := http.NewRequest("POST", p.URL, bytes.NewBuffer(by))
	q.Header.Set("Content-Type", "text/xml; charset=utf-8")

	if p.body.Body.Request == nil {
		o.err = errors.New("No Soap Body Created")
		return o
	}

	q.Header.Set("SOAPAction", soapAction)

	client := &http.Client{}
	resp, err := client.Do(q)

	o.status = resp.Status
	o.header = resp.Header
	o.messageID = p.body.Header.SyncRequestHeaderInfo.MessageIdentifier

	if err != nil {
		o.err = err
		return o
	}
	defer resp.Body.Close()
	bodys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		o.err = err
		return o
	}
	o.resBody = bodys

	if !strings.Contains(resp.Status, "200") {
		o.err = errors.New(resp.Status)
		return o
	}
	return o
}

type Response struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	S       string   `xml:"s,attr"`
	U       string   `xml:"u,attr"`
	Header  struct {
		Text                   string `xml:",chardata"`
		SyncResponseHeaderInfo struct {
			Text              string        `xml:",chardata"`
			H                 string        `xml:"h,attr"`
			Xmlns             string        `xml:"xmlns,attr"`
			Xsi               string        `xml:"xsi,attr"`
			Xsd               string        `xml:"xsd,attr"`
			MessageIdentifier string        `xml:"messageIdentifier"`
			StatusInfoSet     StatusInfoSet `xml:"statusInfoSet"`
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
		Text  string `xml:",chardata"`
		Fault struct {
			Text      string `xml:",chardata"`
			Faultcode struct {
				Text string `xml:",chardata"`
				A    string `xml:"a,attr"`
			} `xml:"faultcode"`
			Faultstring struct {
				Text string `xml:",chardata"`
				Lang string `xml:"lang,attr"`
			} `xml:"faultstring"`
		} `xml:"Fault"`
	} `xml:"Body"`
}

type StatusInfoSet struct {
	XMLName    xml.Name `xml:"statusInfoSet"`
	Text       string   `xml:",chardata"`
	StatusInfo []struct {
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
}

func (p *calloutput) ResponseBody() []byte {
	return p.resBody
}

func (p *calloutput) GetBody() string {
	return string(p.resBody)
}

func (p *calloutput) GetOutput() Output {
	fmt.Println(string(p.resBody))
	out := new(Output)

	if p.err != nil {
		out.Err = p.err
		return *out
	}
	out.Header = p.header
	out.MessageID = p.messageID
	out.Status = p.status

	var resp Response
	err := xml.Unmarshal(p.resBody, &resp)
	if err != nil {
		out.Err = err
		return *out
	}

	/// Error Checking
	/*
		Multi Status
	*/
	if len(resp.Header.SyncResponseHeaderInfo.StatusInfoSet.StatusInfo) > 0 {
		fmt.Println(len(resp.Header.SyncResponseHeaderInfo.StatusInfoSet.StatusInfo))
		fmt.Println(p.sendData)

		if len(resp.Header.SyncResponseHeaderInfo.StatusInfoSet.StatusInfo) != len(p.sendData) {
			out.Err = errors.New("Length Problem: " + string(p.resBody))
			return *out
		}

		out.MultResult = make(map[string]string)

		for i := 0; i < len(resp.Header.SyncResponseHeaderInfo.StatusInfoSet.StatusInfo); i++ {
			if resp.Header.SyncResponseHeaderInfo.StatusInfo.Description.Text.Text != "" {
				out.MultResult[p.sendData[i]] = resp.Header.SyncResponseHeaderInfo.StatusInfoSet.StatusInfo[i].Description.Text.Text
			} else {
				out.MultResult[p.sendData[i]] = resp.Header.SyncResponseHeaderInfo.StatusInfoSet.StatusInfo[i].CodeMajor
			}
		}
		return *out
	}

	if resp.Body.Fault.Faultstring.Text != "" {
		out.Err = errors.New(resp.Body.Fault.Faultstring.Text)
		return *out
	}

	if resp.Header.SyncResponseHeaderInfo.StatusInfo.Severity == "Warning" {
		out.Warning = resp.Header.SyncResponseHeaderInfo.StatusInfo.Description.Text.Text
		return *out
	}

	if resp.Header.SyncResponseHeaderInfo.StatusInfo.CodeMajor == "Failure" {
		out.Failure = resp.Header.SyncResponseHeaderInfo.StatusInfo.Description.Text.Text
		out.Err = errors.New(out.Failure)
		return *out
	}

	if resp.Header.SyncResponseHeaderInfo.StatusInfo.CodeMajor != "success" {
		out.Err = errors.New("unknown error")
		return *out
	}
	return *out
}

type Output struct {
	Status     string
	Header     http.Header
	MessageID  string
	Err        error
	Warning    string
	Failure    string
	MultResult map[string]string
}
