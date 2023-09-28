package itswizard_m_imses

import (
	"encoding/xml"
	"errors"
	"github.com/itslearninggermany/itswizard_m_basic"
	"strings"
)

type CreateCourse struct {
	XMLName   xml.Name `xml:"ims4:createGroupRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims4:sourcedId"`
	Group struct {
		GroupType struct {
			Scheme    string `xml:"ims5:scheme"`
			TypeValue struct {
				Type string `xml:"ims5:type"`
			} `xml:"ims5:typeValue"`
		} `xml:"ims5:groupType"`
		Relationship struct {
			Relation string `xml:"ims5:relation"`
			SourceId struct {
				Identifier string `xml:"ims2:identifier"`
			} `xml:"ims5:sourceId"`
		} `xml:"ims5:relationship"`
		Description struct {
			DescShort string `xml:"ims5:descShort"`
		} `xml:"ims5:description"`
		Extension struct {
			ExtensionField []OneExtensionField `xml:"ims2:extensionField"`
		} `xml:"ims5:extension"`
	} `xml:"ims4:group"`
}

type OneExtensionField struct {
	FieldName  string `xml:"ims2:fieldName"`
	FieldType  string `xml:"ims2:fieldType"`
	FieldValue string `xml:"ims2:fieldValue"`
}

func (p *Request) CreateCourse(group itswizard_m_basic.DbGroup15) (resp string, err error) {

	cnp := CreateCourse{}

	cnp.SourcedId.Identifier = group.SyncID
	cnp.Group.GroupType.Scheme = "ItslearningOrganisationTypes"

	cnp.Group.GroupType.TypeValue.Type = "Course"

	cnp.Group.Relationship.Relation = "Parent"
	cnp.Group.Relationship.SourceId.Identifier = group.ParentGroupID
	cnp.Group.Description.DescShort = group.Name

	var extensionFields []OneExtensionField

	extensionFields = append(extensionFields, OneExtensionField{
		FieldName:  "course",
		FieldType:  "String",
		FieldValue: group.Name,
	})

	extensionFields = append(extensionFields, OneExtensionField{
		FieldName:  "course/code",
		FieldType:  "String",
		FieldValue: group.Name,
	})

	cnp.Group.Extension.ExtensionField = extensionFields

	p.body.Body.Request = cnp

	resp = p.call("http://www.imsglobal.org/soap/gms/createGroup", nil).GetBody()

	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}

}
