package itswizard_m_imses

import (
	"encoding/xml"
	"errors"
	"github.com/itslearninggermany/itswizard_m_basic"
	"strings"
)

type CreateGroupRequest struct {
	XMLName   xml.Name `xml:"ims4:createGroupRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims4:sourcedId"`
	Group struct {
		GroupType struct {
			Scheme    string `xml:"ims5:scheme"`
			TypeValue struct {
				Type  string `xml:"ims5:type"`
				Level string `xml:"ims5:level"`
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
	} `xml:"ims4:group"`
}

func (p *Request) CreateGroup(group itswizard_m_basic.DbGroup15, organisation bool) (resp string, err error) {

	cnp := CreateGroupRequest{}

	cnp.SourcedId.Identifier = group.SyncID
	cnp.Group.GroupType.Scheme = "ItslearningOrganisationTypes"

	if organisation {
		cnp.Group.GroupType.TypeValue.Type = "School"
		cnp.Group.GroupType.TypeValue.Level = "1"
	} else {
		cnp.Group.GroupType.TypeValue.Type = "Unspecified"
		cnp.Group.GroupType.TypeValue.Level = "-1"
	}

	cnp.Group.Relationship.Relation = "Parent"
	cnp.Group.Relationship.SourceId.Identifier = group.ParentGroupID
	// short Description
	cnp.Group.Description.DescShort = shortStringToLength(group.Name, 64)

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

func (p *Request) CreateGroupNewGroup(group itswizard_m_basic.Group, organisation bool) (resp string, err error) {

	cnp := CreateGroupRequest{}

	cnp.SourcedId.Identifier = group.GroupSyncKey
	cnp.Group.GroupType.Scheme = "ItslearningOrganisationTypes"

	if organisation {
		cnp.Group.GroupType.TypeValue.Type = "School"
		cnp.Group.GroupType.TypeValue.Level = "1"
	} else {
		cnp.Group.GroupType.TypeValue.Type = "Unspecified"
		cnp.Group.GroupType.TypeValue.Level = "-1"
	}

	cnp.Group.Relationship.Relation = "Parent"
	cnp.Group.Relationship.SourceId.Identifier = group.ParentGroupID
	cnp.Group.Description.DescShort = group.Name

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

type UpdateGroupRequest struct {
	XMLName   xml.Name `xml:"ims4:updateGroupRequest"`
	SourcedId struct {
		Identifier string `xml:"ims2:identifier"`
	} `xml:"ims4:sourcedId"`
	Group struct {
		GroupType struct {
			Scheme    string `xml:"ims5:scheme"`
			TypeValue struct {
				Type  string `xml:"ims5:type"`
				Level string `xml:"ims5:level"`
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
	} `xml:"ims4:group"`
}

func (p *Request) UpdateGroup(group itswizard_m_basic.Group, organisation bool) (resp string, err error) {

	cnp := UpdateGroupRequest{}

	cnp.SourcedId.Identifier = group.GroupSyncKey
	cnp.Group.GroupType.Scheme = "ItslearningOrganisationTypes"

	if organisation {
		cnp.Group.GroupType.TypeValue.Type = "School"
		cnp.Group.GroupType.TypeValue.Level = "1"
	} else {
		cnp.Group.GroupType.TypeValue.Type = "Unspecified"
		cnp.Group.GroupType.TypeValue.Level = "-1"
	}

	cnp.Group.Relationship.Relation = "Parent"
	cnp.Group.Relationship.SourceId.Identifier = group.ParentGroupID
	cnp.Group.Description.DescShort = group.Name

	p.body.Body.Request = cnp

	resp = p.call("http://www.imsglobal.org/soap/gms/updateGroup", nil).GetBody()

	if strings.Contains(resp, "success") {
		err = nil
		return
	} else {
		err = errors.New("See resp for more details")
		return
	}
}

func shortStringToLength(input string, length int) string {
	i := 0
	for len(input) > length {
		input = input[:len(input)-i]
		i++
	}
	return input
}
