/* Maltego library for Go
 * stuart@sensepost.com //@noobiedog
 * glenn@sensepost.com // @glennzw
 *
 * Implemented almost verbatim from the Maltego.py library
 */

package maltegolocal

import (
	"encoding/xml"
	"strconv"
	"fmt"
)
// ############### STRUCTS #################

type MaltegoEntityObj struct {
	entityType         string
	value              string
	iconURL            string
	weight             int
	displayInformation [][]string
	AdditionalFields   [][]string
}

type MaltegoTransform struct {
	entities   []*MaltegoEntityObj
	exceptions [][]string
	UIMessages [][]string
}

type DisplayInformationa struct {
	Nameattrd string `xml:"Name,attr"`
	Typeattrd string `xml:"Type,attr"`
	Textd []byte `xml:",chardata"`    
}

type AdditionalFieldsa struct {
	Matchingrulea string `xml:"MatchingRule,attr"`
	Nameattra string `xml:"Name,attr"`
	Typeattra string `xml:"DisplayName,attr"`
	Texta string `xml:",chardata"`    
}

type Entitya struct {
	XMLName   xml.Name `xml:"Entity"`
	Type        string      `xml:"Type,attr"`
	Value string   `xml:"Value"`
	Weight  int   `xml:"Weight"`
	Dispinfo     DisplayInformationa      `xml:"DisplayInformation>Label,inline"`
	Addfield     []AdditionalFieldsa      `xml:"AdditionalFields>Field,inline"`
	IconURL    string  `xml:"IconURL"`
}

type LocalTransform struct {
	Value  string
	Values map[string]string
}

// ################# FUNCTIONS ##################

func MaltegoEntity(eT string, eV string) *MaltegoEntityObj {
	return &MaltegoEntityObj{entityType: eT, value: eV, weight: 100}
}

// Add data to Fields
func (m *MaltegoTransform) AddEntity(enType, enValue string) *MaltegoEntityObj {
	me := &MaltegoEntityObj{entityType: enType, value: enValue, weight: 100}
	m.entities = append(m.entities, me)
	return me
}

func (m *MaltegoEntityObj) AddProperty(fieldName, displayName, matchingRule, value string) {
	prop := []string{fieldName, displayName, matchingRule, value}
	m.AdditionalFields = append(m.AdditionalFields, prop)
}

func (m *MaltegoEntityObj) AddDisplayInformation(di, dl string) {
	info := []string{dl, di}
	m.displayInformation = append(m.displayInformation, info)
}

func (m *MaltegoTransform) AddUIMessage(message, messageType string) {
	msg := []string{messageType, message}
	m.UIMessages = append(m.UIMessages, msg)
}

func (m *MaltegoTransform) AddException(exceptionString, code string) {
	exc := []string{exceptionString, code}
	m.exceptions = append(m.exceptions, exc)
}

// Set Data Values
func (m *MaltegoEntityObj) SetType(eT string) {
	m.entityType = eT
}

func (m *MaltegoEntityObj) SetValue(eV string) {
	m.value = eV
}

func (m *MaltegoEntityObj) SetWeight(w int) {
	m.weight = w
}

func (m *MaltegoEntityObj) SetIconURL(iU string) {
	m.iconURL = iU
}

func (m *MaltegoEntityObj) SetLinkColor(color string) {
	m.AddProperty("link#maltego.link.color", "LinkColor", "", color)
}

func (m *MaltegoEntityObj) SetLinkStyle(style string) {
	m.AddProperty("link#maltego.link.style", "LinkStyle", "", style)
}

func (m *MaltegoEntityObj) SetLinkThickness(thick int) {
	thickInt := strconv.Itoa(thick)
	m.AddProperty("link#maltego.link.style", "LinkStyle", "", thickInt)
}

func (m *MaltegoEntityObj) SetLinkLabel(label string) {
	m.AddProperty("link#maltego.link.label", "Label", "", label)
}

func (m *MaltegoEntityObj) SetBookmark(bookmark string) {
	m.AddProperty("bookmark#", "Bookmark", "", bookmark)
}

func (m *MaltegoEntityObj) SetNote(note string) {
	m.AddProperty("notes#", "Notes", "", note)
}

func ParseLocalArguments(args []string) LocalTransform {
	Value := args[1]
	Vals := make(map[string]string)
	if len(args) > 2 {
		vars := strings.Split(args[2], "#")
		for _, x := range vars {
			kv := strings.Split(x, "=")
			Vals[kv[0]] = kv[1]
		}
	}
	return LocalTransform{Value: Value, Values: Vals}
}

// ################ XML OUTPUTS ###############
// Rework these into XML output, not printed
func (m *MaltegoTransform) ReturnOutput() string {
	r := "<MaltegoMessage>\n"
	r += "<MaltegoTransformResponseMessage>\n"
	r += "<Entities>\n"
	for _, e := range m.entities {
		r += e.ReturnEntity()
	}
	r += "</Entities>\n"
	r += "<UIMessages>\n"
	for _, e := range m.UIMessages {
		mType, mVal := e[0], e[1]
		r += "<UIMessage MessageType=\"" + mType + "\">" + mVal + "</UIMessage>\n"
	}
	r += "</UIMessages>\n"
	r += "</MaltegoTransformResponseMessage>\n"
	r += "</MaltegoMessage>\n"
	return r
}

func (m *MaltegoTransform) ThrowExceptions() string {
	r := "<MaltegoMessage>\n"
	r += "<MaltegoTransformExceptionMessage>\n"
	r += "<Exceptions>\n"
	for _, e := range m.exceptions {
		code, ex := e[0], e[1]
		r += "<Exception code='" + code + "'>" + ex + "</Exception>\n"
	}
	r += "</Exceptions>\n"
	r += "</MaltegoTransformExceptionMessage>\n"
	r += "</MaltegoMessage>\n"
	return r
}

func (m *MaltegoEntityObj) ReturnEntity() string {	
	r := &Entitya{Type: m.entityType, Value: m.value, Weight: m.weight, IconURL:m.iconURL}
	// Add DisplayInformation XML
	if len(m.displayInformation) > 0 {
		for _, e := range m.displayInformation {
			name_, type_ := e[0], e[1]
			r.Dispinfo = DisplayInformationa{Nameattrd:name_, Typeattrd:"text/html", Textd:[]byte("<![CDATA[" + type_ + "]]>")}
		}
	} else {
			r.Dispinfo = DisplayInformationa{Nameattrd:"", Typeattrd:"text/html", Textd:[]byte("")}
	}
	// Add AdditionalFields XML
	if len(m.AdditionalFields) > 0 {
		for _, e := range m.AdditionalFields {
			fieldName_, displayName_, matchingRule_, value_ := e[0], e[1], e[2], e[3]
			if matchingRule_ == "strict" {
				r.Addfield = append(r.Addfield,AdditionalFieldsa{Nameattra:fieldName_, Typeattra:displayName_, Texta:value_})
			} else {
				r.Addfield = append(r.Addfield,AdditionalFieldsa{Matchingrulea:matchingRule_,Nameattra:fieldName_, Typeattra:displayName_, Texta:value_})
			}
		}
	}
	output, err := xml.MarshalIndent(r, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return string(output)
}