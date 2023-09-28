package itswizard_m_imses

import (
	"encoding/xml"
	"errors"
	"github.com/itslearninggermany/itswizard_m_basic"
	"strconv"
	"strings"
)

type CreatePersonRequest struct {
	XMLName   xml.Name `xml:"ims:createPersonRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims:sourcedId"`
	Person Person `xml:"ims:person"`
}

type PartName struct {
	NamePartType  string `xml:"ims3:namePartType"`
	NamePartValue string `xml:"ims3:namePartValue"`
}

type Tel struct {
	XMLName  xml.Name `xml:"ims3:tel"`
	TelType  string   `xml:"ims3:telType"`
	TelValue string   `xml:"ims3:telValue"`
}

type Person struct {
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
	Tel []Tel `xml:"ims3:tel"`
}

func (p *Request) CreatePerson(person itswizard_m_basic.DbPerson15) (resp string, err error) {
	cnp := CreatePersonRequest{}

	cnp.Person.FormatName.Nil = "true"
	cnp.Person.FormatName.Xsi = "http://www.w3.org/2001/XMLSchema-instance"

	cnp.SourcedId.Identifier = person.SyncPersonKey

	names := []PartName{{
		NamePartType:  "First",
		NamePartValue: shortStringToLength(person.FirstName, 64),
	}, {
		NamePartType:  "Last",
		NamePartValue: shortStringToLength(person.LastName, 64),
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

	p.body.Body.Request = cnp

	resp = p.call("http://www.imsglobal.org/soap/pms/createPerson", nil).GetBody()

	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}
}

type changeNames struct {
	XMLName   xml.Name `xml:"ims:updatePersonRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims:sourcedId"`
	Person struct {
		FormatName struct {
			Xsi string `xml:"xmlns:xsi,attr"`
			Nil string `xml:"xsi:nil,attr"`
		} `xml:"ims3:formatName"`
		Name struct {
			PartName []PartName `xml:"ims3:partName"`
		} `xml:"ims3:name"`
		InstitutionRole struct {
			InstitutionRoleType string `xml:"ims3:institutionRoleType"`
			PrimaryRoleType     string `xml:"ims3:primaryRoleType"`
		} `xml:"ims3:institutionRole"`
	} `xml:"ims:person"`
}

type changeName struct {
	XMLName   xml.Name `xml:"ims:updatePersonRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims:sourcedId"`
	Person struct {
		FormatName struct {
			Xsi string `xml:"xmlns:xsi,attr"`
			Nil string `xml:"xsi:nil,attr"`
		} `xml:"ims3:formatName"`
		Name struct {
			PartName struct {
				NamePartType  string `xml:"ims3:namePartType"`
				NamePartValue string `xml:"ims3:namePartValue"`
			} `xml:"ims3:partName"`
		} `xml:"ims3:name"`
	} `xml:"ims:person"`
}

func (p *Request) UpdateFirstName(SyncPersonKey string, firstName string) (resp string, err error) {
	cnp := changeName{}

	cnp.Person.FormatName.Nil = "true"
	cnp.Person.FormatName.Xsi = "http://www.w3.org/2001/XMLSchema-instance"

	cnp.SourcedId.Identifier = SyncPersonKey

	cnp.Person.Name.PartName.NamePartType = "First"
	cnp.Person.Name.PartName.NamePartValue = firstName

	p.body.Body.Request = cnp
	resp = p.call("http://www.imsglobal.org/soap/pms/updatePerson", nil).GetBody()
	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}
}

func (p *Request) UpdateLastName(SyncPersonKey string, lastName string) (resp string, err error) {
	cnp := changeName{}

	cnp.Person.FormatName.Nil = "true"
	cnp.Person.FormatName.Xsi = "http://www.w3.org/2001/XMLSchema-instance"

	cnp.SourcedId.Identifier = SyncPersonKey

	cnp.Person.Name.PartName.NamePartType = "Last"
	cnp.Person.Name.PartName.NamePartValue = lastName

	p.body.Body.Request = cnp
	resp = p.call("http://www.imsglobal.org/soap/pms/updatePerson", nil).GetBody()
	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}
}

func (p *Request) UpdateNames(SyncPersonKey string, lastName string, firstName string, profile string) (resp string, err error) {
	cnp := changeNames{}

	cnp.Person.FormatName.Nil = "true"
	cnp.Person.FormatName.Xsi = "http://www.w3.org/2001/XMLSchema-instance"

	cnp.SourcedId.Identifier = SyncPersonKey

	names := []PartName{{
		NamePartType:  "First",
		NamePartValue: firstName,
	}, {
		NamePartType:  "Last",
		NamePartValue: lastName,
	}}

	cnp.Person.Name.PartName = names

	cnp.Person.InstitutionRole.InstitutionRoleType = profile
	cnp.Person.InstitutionRole.PrimaryRoleType = "true"

	p.body.Body.Request = cnp
	resp = p.call("http://www.imsglobal.org/soap/pms/updatePerson", nil).GetBody()
	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}
}

type changeEmail struct {
	XMLName   xml.Name `xml:"ims:updatePersonRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims:sourcedId"`
	Person struct {
		Email string `xml:"ims2:email"`
	} `xml:"ims:person"`
}

func (p *Request) UpdateEmail(SyncPersonKey string, eMail string) (resp string, err error) {
	cnp := changeEmail{}
	cnp.SourcedId.Identifier = SyncPersonKey
	cnp.Person.Email = eMail

	p.body.Body.Request = cnp
	resp = p.call("http://www.imsglobal.org/soap/pms/updatePerson", nil).GetBody()
	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}
}

type changeUsername struct {
	XMLName   xml.Name `xml:"ims:updatePersonRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims:sourcedId"`
	Person struct {
		UserId struct {
			UserIdValue string `xml:"ims2:userIdValue"`
		} `xml:"ims3:userId"`
	} `xml:"ims:person"`
}

func (p *Request) UpdateUsername(SyncPersonKey string, username string) (resp string, err error) {
	cnp := changeUsername{}
	cnp.SourcedId.Identifier = SyncPersonKey
	cnp.Person.UserId.UserIdValue = username
	p.body.Body.Request = cnp
	resp = p.call("http://www.imsglobal.org/soap/pms/updatePerson", nil).GetBody()
	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}
}

type changeAdress struct {
	XMLName   xml.Name `xml:"ims:updatePersonRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims:sourcedId"`
	Person struct {
		Address struct {
			Locality string   `xml:"ims3:locality"`
			Postcode string   `xml:"ims3:postcode"`
			Street   []string `xml:"ims3:street"`
		} `xml:"ims3:address"`
	} `xml:"ims:person"`
}

func (p *Request) UpdateAdresse(SyncPersonKey, postcode, city, street1, street2 string) (resp string, err error) {
	cnp := changeAdress{}
	cnp.SourcedId.Identifier = SyncPersonKey
	x := []string{street1, street2}
	cnp.Person.Address.Locality = city
	cnp.Person.Address.Postcode = postcode
	cnp.Person.Address.Street = x

	p.body.Body.Request = cnp
	resp = p.call("http://www.imsglobal.org/soap/pms/updatePerson", nil).GetBody()
	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}
}

type changeProfile struct {
	XMLName   xml.Name `xml:"ims:updatePersonRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims:sourcedId"`
	Person struct {
		InstitutionRole struct {
			InstitutionRoleType string `xml:"ims3:institutionRoleType"`
			PrimaryRoleType     string `xml:"ims3:primaryRoleType"`
		} `xml:"ims3:institutionRole"`
	} `xml:"ims:person"`
}

func (p *Request) UpdateProfile(SyncPersonKey, profile string) (resp string, err error) {
	cnp := changeProfile{}
	cnp.SourcedId.Identifier = SyncPersonKey
	cnp.Person.InstitutionRole.InstitutionRoleType = profile
	cnp.Person.InstitutionRole.PrimaryRoleType = "true"

	p.body.Body.Request = cnp
	resp = p.call("http://www.imsglobal.org/soap/pms/updatePerson", nil).GetBody()
	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}
}

type changePhone struct {
	XMLName   xml.Name `xml:"ims:updatePersonRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims:sourcedId"`
	Person struct {
		Tel struct {
			TelType  string `xml:"ims3:telType"`
			TelValue string `xml:"ims3:telValue"`
		} `xml:"ims3:tel"`
	} `xml:"ims:person"`
}

func (p *Request) UpdatePhone(SyncPersonKey, phonenumber string) (resp string, err error) {
	cnp := changePhone{}
	cnp.SourcedId.Identifier = SyncPersonKey
	cnp.Person.Tel.TelType = "Voice"
	cnp.Person.Tel.TelValue = phonenumber

	p.body.Body.Request = cnp
	resp = p.call("http://www.imsglobal.org/soap/pms/updatePerson", nil).GetBody()
	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}
}

func (p *Request) UpdateMobilePhone(SyncPersonKey, phonenumber string) (resp string, err error) {
	cnp := changePhone{}
	cnp.SourcedId.Identifier = SyncPersonKey
	cnp.Person.Tel.TelType = "Mobile"
	cnp.Person.Tel.TelValue = phonenumber

	p.body.Body.Request = cnp
	resp = p.call("http://www.imsglobal.org/soap/pms/updatePerson", nil).GetBody()
	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}
}

type changeBirthday struct {
	XMLName   xml.Name `xml:"ims:updatePersonRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims:sourcedId"`
	Person struct {
		Demographics struct {
			Bday string `xml:"ims3:bday"`
		} `xml:"ims3:demographics"`
	} `xml:"ims:person"`
}

func (p *Request) UpdateBirthday(SyncPersonKey string, day, month, year uint) (resp string, err error) {
	cnp := changeBirthday{}
	cnp.SourcedId.Identifier = SyncPersonKey
	cnp.Person.Demographics.Bday = strconv.Itoa(int(year)) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(int(day))

	p.body.Body.Request = cnp
	resp = p.call("http://www.imsglobal.org/soap/pms/updatePerson", nil).GetBody()
	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}
}

type changePassword struct {
	XMLName   xml.Name `xml:"ims:updatePersonRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims:sourcedId"`
	Person struct {
		UserId struct {
			PassWord string `xml:"ims2:passWord"`
		} `xml:"ims3:userId"`
	} `xml:"ims:person"`
}

func (p *Request) UpdatePassword(SyncPersonKey, password string) (resp string, err error) {
	cnp := changePassword{}
	cnp.SourcedId.Identifier = SyncPersonKey
	cnp.Person.UserId.PassWord = password

	p.body.Body.Request = cnp
	resp = p.call("http://www.imsglobal.org/soap/pms/updatePerson", nil).GetBody()
	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}
}
