package mail

import (
	"github.com/GeertJohan/go.rice"
	"html/template"
	log "github.com/sirupsen/logrus"
)

var RiceBox *rice.Box

type Compiler struct {
}

func (c Compiler) Compile(box *rice.Box, boxName, templateName string) (*template.Template, error) {

	templateString, err := box.String(boxName)

	if err != nil {
		return nil, err
	}

	return template.New(templateName).Parse(templateString)
}

func (c Compiler) MustCompile(box *rice.Box, boxName, templateName string) (*template.Template) {

	templateString, err := RiceBox.String(boxName)

	if err != nil {
		log.Errorln("error occurred while find rice box with name", boxName)
		log.Errorln("error is: ", err)
		return nil
	}

	tmpl, err := template.New(templateName).Parse(templateString)

	if err != nil {
		log.Errorln("error occurred while parsing template with name", templateName)
		log.Errorln("error is: ", err)
		return nil
	}

	return tmpl
}