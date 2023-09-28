package itswizard_m_imses

import (
	"encoding/xml"
	"errors"
	"strings"
)

type UpdatePersonRequest struct {
	XMLName   xml.Name `xml:"ims:updatePersonRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims:sourcedId"`
	Person PersonPSR `xml:"ims:person"`
}

type PersonPSR struct {
	FormatName struct {
		Xsi string `xml:"xmlns:xsi,attr"`
		Nil string `xml:"xsi:nil,attr"`
	} `xml:"ims3:formatName"`
	Name struct {
		PartName []PartName `xml:"ims3:partName"`
	} `xml:"ims3:name"`
	Email  string `xml:"ims2:email"`
	UserId struct {
		UserIdValue string `xml:"ims2:userIdValue"`
		PassWord    string `xml:"ims2:passWord"`
	} `xml:"ims3:userId"`
	Address struct {
		Locality string   `xml:"ims3:locality"`
		Postcode string   `xml:"ims3:postcode"`
		Street   []string `xml:"ims3:street"`
	} `xml:"ims3:address"`
	/*
		//Todo: Every Birthday mus be filled out, when it is empty HTTP.Status: 500
		Demographics struct {
			Bday string `xml:"ims3:bday"`
		} `xml:"ims3:demographics"`
	*/
	InstitutionRole struct {
		InstitutionRoleType string `xml:"ims3:institutionRoleType"`
		PrimaryRoleType     string `xml:"ims3:primaryRoleType"`
	} `xml:"ims3:institutionRole"`
	Tel       []Tel `xml:"ims3:tel"`
	Extension struct {
		Relationship []RelationshipData `xml:"ims2:relationship"`
	} `xml:"ims3:extension"`
}

type RelationshipData struct {
	Relation string `xml:"ims2:relation"`
	SourceId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims2:sourceId"`
}

func (p *Request) CreateParentChildRelaionship(parentId string, childIds []string) error {
	out := p.ReadPerson(parentId)
	if out.Err != nil {
		return out.Err
	}
	person := out.Person

	cnp := UpdatePersonRequest{}

	cnp.Person.FormatName.Nil = "true"
	cnp.Person.FormatName.Xsi = "http://www.w3.org/2001/XMLSchema-instance"

	cnp.SourcedId.Identifier = person.SyncPersonKey

	names := []PartName{{
		NamePartType:  "First",
		NamePartValue: person.FirstName,
	}, {
		NamePartType:  "Last",
		NamePartValue: person.LastName,
	}}

	cnp.Person.Name.PartName = names

	cnp.Person.UserId.UserIdValue = person.Username
	cnp.Person.UserId.PassWord = person.Password
	cnp.Person.InstitutionRole.InstitutionRoleType = person.Profile
	cnp.Person.InstitutionRole.PrimaryRoleType = "true"

	cnp.Person.Email = person.Email

	telephones := []Tel{{
		TelType:  "Voice",
		TelValue: person.Phone,
	}, {
		TelType:  "Mobile",
		TelValue: person.Mobile,
	}}

	cnp.Person.Tel = telephones

	streets := []string{person.Street1, person.Street2}

	cnp.Person.Address.Street = streets
	cnp.Person.Address.Locality = person.City
	cnp.Person.Address.Postcode = person.Postcode

	var tmp []RelationshipData

	for i := 0; i < len(childIds); i++ {
		tmp = append(tmp, RelationshipData{
			Relation: "Child",
			SourceId: struct {
				Identifier string `xml:"ims2:identifier"`
			}{childIds[i]},
		})
	}

	cnp.Person.Extension.Relationship = tmp

	p.body.Body.Request = cnp

	resp := p.call("http://www.imsglobal.org/soap/pms/updatePerson", nil).GetBody()

	if strings.Contains(resp, "success") {
		return nil
	} else {
		return errors.New("See resp for more details:" + resp)
	}
}
