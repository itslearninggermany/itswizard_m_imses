package itswizard_m_imses

import "encoding/xml"

type ChangePersonsIdentifierRequest struct {
	XMLName          xml.Name `xml:"ims:changePersonsIdentifierRequest"`
	PairSourcedIdSet struct {
		IdentifierPairs []IdentifierPair
	} `xml:"ims:pairSourcedIdSet"`
}

type IdentifierPair struct {
	XMLName  xml.Name `xml:"ims2:identifierPair"`
	FirstId  string   `xml:"ims2:firstId"`
	SecondId string   `xml:"ims2:secondId"`
}

/*
map[oldId]newID
*/
func (p *Request) ChangePersonsIdentifier(input map[string]string) Output {

	cpi := ChangePersonsIdentifierRequest{}

	var identifierPairs []IdentifierPair
	for oldID, newID := range input {
		identifierPairs = append(identifierPairs, IdentifierPair{
			FirstId:  oldID,
			SecondId: newID,
		})
	}

	var all []string
	for i := 0; i < len(identifierPairs); i++ {
		all = append(all, "old ID: "+identifierPairs[i].FirstId+" new ID: "+identifierPairs[i].SecondId)
	}

	cpi.PairSourcedIdSet.IdentifierPairs = identifierPairs
	p.body.Body.Request = cpi

	return p.call("http://www.imsglobal.org/soap/pms/changePersonsIdentifier", all).GetOutput()
}
