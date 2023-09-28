package itswizard_m_imses

import (
	"encoding/xml"
	"github.com/itslearninggermany/itswizard_m_basic"
)

type CreatePersonsRequest struct {
	XMLName         xml.Name `xml:"ims:createPersonsRequest"`
	PersonIdPairSet struct {
		PersonIdPair []PersonIdPair
	} `xml:"ims:personIdPairSet"`
}

type PersonIdPair struct {
	XMLName   xml.Name `xml:"ims:personIdPair"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims:sourcedId"`
	Person Person `xml:"ims:person"`
}

func (p *Request) CreatePersons(persons []itswizard_m_basic.DbPerson15) Output {

	var personIdPairs []PersonIdPair

	for i := 0; i < len(persons); i++ {

		var pip PersonIdPair

		pip.Person.FormatName.Nil = "true"
		pip.Person.FormatName.Xsi = "http://www.w3.org/2001/XMLSchema-instance"

		pip.SourcedId.Identifier = persons[i].SyncPersonKey

		names := []PartName{{
			NamePartType:  "First",
			NamePartValue: persons[i].FirstName,
		}, {
			NamePartType:  "Last",
			NamePartValue: persons[i].LastName,
		}}

		pip.Person.Name.PartName = names

		pip.Person.UserId.UserIdValue = persons[i].Username
		pip.Person.UserId.PassWord = persons[i].Password
		pip.Person.InstitutionRole.InstitutionRoleType = persons[i].Profile
		pip.Person.InstitutionRole.PrimaryRoleType = "true"

		pip.Person.Email = persons[i].Email

		telephones := []Tel{{
			TelType:  "Voice",
			TelValue: persons[i].Phone,
		}, {
			TelType:  "Mobile",
			TelValue: persons[i].Mobile,
		}}

		pip.Person.Tel = telephones

		streets := []string{persons[i].Street1, persons[i].Street2}

		pip.Person.Address.Street = streets
		pip.Person.Address.Locality = persons[i].City
		pip.Person.Address.Postcode = persons[i].Postcode

		personIdPairs = append(personIdPairs, pip)

	}

	var createPersonsRequest CreatePersonsRequest
	createPersonsRequest.PersonIdPairSet.PersonIdPair = personIdPairs

	p.body.Body.Request = createPersonsRequest

	var all []string
	for i := 0; i < len(persons); i++ {
		all = append(all, "SyncKey: "+persons[i].SyncPersonKey+" Username: "+persons[i].Username)
	}

	return p.call("http://www.imsglobal.org/soap/pms/createPersons", all).GetOutput()

}
