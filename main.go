package main

import (
	"github.com/graphicweave/injun/cmd"
	"github.com/graphicweave/injun/mail"
	"github.com/GeertJohan/go.rice"
	"github.com/sirupsen/logrus"
)

func main() {
	var err error
	mail.RiceBox , err =  rice.FindBox("mail-templates")
	if err != nil {
		logrus.Fatal("failed to find rice box.", err)
	}
	cmd.RootCmd.Execute()
}