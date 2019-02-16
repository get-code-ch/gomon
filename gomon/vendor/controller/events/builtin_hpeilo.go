package events

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

type IloHealthAtGlance struct {
	BiosHardware  string
	Fans          string
	Temperature   string
	PowerSupplies string
	Processor     string
	Memory        string
	Network       string
	Storage       string
}

func (hg *IloHealthAtGlance) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type StatusAttr struct {
		Status string `xml:"STATUS,attr"`
	}
	var health struct {
		XMLName xml.Name   `xml:"HEALTH_AT_A_GLANCE"`
		BH      StatusAttr `xml:"BIOS_HARDWARE"`
		F       StatusAttr `xml:"FANS"`
		T       StatusAttr `xml:"TEMPERATURE"`
		PS      StatusAttr `xml:"POWER_SUPPLIES"`
		P       StatusAttr `xml:"PROCESSOR"`
		M       StatusAttr `xml:"MEMORY"`
		N       StatusAttr `xml:"NETWORK"`
		S       StatusAttr `xml:"STORAGE"`
	}
	err := d.DecodeElement(&health, &start)
	if err == nil {
		hg.BiosHardware = health.BH.Status
		hg.Fans = health.F.Status
		hg.Temperature = health.T.Status
		hg.PowerSupplies = health.PS.Status
		hg.Processor = health.P.Status
		hg.Memory = health.M.Status
		hg.Network = health.N.Status
		hg.Storage = health.S.Status
	}
	return err
}

func GetIloHealth(url string, username string, password string) (state string, message string, metric string) {
	var health IloHealthAtGlance
	var re = regexp.MustCompile(`(?mi)<HEALTH_AT_A_GLANCE>[\s\S]*</HEALTH_AT_A_GLANCE>`)

	var request = "<?xml version=\"1.0\"?>\n"
	request += "<RIBCL VERSION=\"2.23\">\n"
	request += "<LOGIN USER_LOGIN=\"" + username + "\" PASSWORD=\"" + password + "\">"
	request += "<SERVER_INFO MODE=\"read\">"
	request += "<GET_EMBEDDED_HEALTH />"
	request += "</SERVER_INFO>"
	request += "</LOGIN>"
	request += "</RIBCL>"

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Post(url+"/ribcl", "application/xml", bytes.NewBuffer([]byte(request)))
	if err != nil {
		return "CRITICAL", err.Error(), ""
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		/*
					body = []byte(`
						<RESPONSE
			        		STATUS="0x0000"
			        		MESSAGE='No error'
						/>
					`)
					/*
					body = []byte(`
				    <HEALTH_AT_A_GLANCE>
			    	    <BIOS_HARDWARE STATUS= "OK" />
			        	<FANS STATUS= "OK"/>
			        	<TEMPERATURE STATUS= "OK"/>
			        	<POWER_SUPPLIES STATUS= "FAILED"/>
			        	<PROCESSOR STATUS= "OK"/>
			        	<MEMORY STATUS= "OK"/>
			        	<NETWORK STATUS= "OK"/>
			        	<STORAGE STATUS= "WARNING"/>
			    	</HEALTH_AT_A_GLANCE>
					`)
					/**/

		haag := re.FindSubmatch(body)
		if haag == nil {
			return "CRITICAL", "Unexpected server response", ""
		}
		err = xml.Unmarshal(haag[0], &health)
		if err == nil {
			ok := ""
			critical := ""

			e := reflect.ValueOf(&health).Elem()
			for i := 0; i < e.NumField(); i++ {
				Name := e.Type().Field(i).Name
				Value := e.Field(i).Interface()

				if strings.ToUpper(Value.(string)) == "OK" {
					if ok != "" {
						ok += " / "
					}
					ok += Name + " : " + Value.(string)
				} else {
					if critical != "" {
						critical += " / "
					}
					critical += Name + " : " + Value.(string)
				}
			}
			if critical != "" {
				return "CRITICAL", critical + " - " + ok, ""
			} else {
				return "OK", ok, ""
			}
		} else {
			return "CRITICAL", "Error parsing server response" + err.Error(), ""
		}
	}
	return "CRITICAL", "", ""
}
