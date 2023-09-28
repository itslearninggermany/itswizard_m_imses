package itswizard_m_imses

import (
	"github.com/itslearninggermany/itswizard_m_basic"
	"log"
)

func (p *Request) ReadAllGroup(rootGroup string, organisationRoot string, OrganisationID uint, InstitutionID uint) (out []itswizard_m_basic.Group) {
	allGroupsTooMuch := make(map[string]bool)

	persons := p.ReadPersonsForGroup(rootGroup, OrganisationID, InstitutionID)
	for _, person := range persons.Persons {
		mems := p.ReadMembershipsForPerson(person.PersonSyncKey)
		for _, mem := range mems {
			allGroupsTooMuch[mem.GroupID] = true
		}
	}

	// Schaue ob im Stamm die RootGroup vorkommt

	allGroupsTooImport := make(map[string]itswizard_m_basic.Group)

	for group, _ := range allGroupsTooMuch {
		tocheck := group
		importG := false
		for {
			log.Println("gerade zu checken", tocheck)
			if tocheck == "" {
				break
			}
			if tocheck == rootGroup {
				break
			}
			gr := p.ReadGroup(tocheck).Group
			if tocheck == organisationRoot {
				importG = true
				break
			}
			if gr.ParentGroupID == organisationRoot {
				importG = true
				break
			}
			if gr.ParentGroupID == rootGroup {
				break
			}
			tocheck = gr.ParentGroupID
		}
		if importG {
			allGroupsTooImport[group] = p.ReadGroup(group).Group
		}
	}

	for _, v := range allGroupsTooImport {
		v.Organisation15 = OrganisationID
		v.Institution15 = InstitutionID
		out = append(out, v)
	}

	return
}
