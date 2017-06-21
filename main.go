package main

import (
	"os"
	"github.com/Bourne-ID/beats-forwarder/forwarder"
	cfg "github.com/Bourne-ID/beats-forwarder/config"
	"github.com/Sirupsen/logrus"
	"flag"
)

var config = cfg.Config{}

func main() {


	flag.Parse()
	debug := flag.Lookup("d").Value.String()

	if (debug == "true") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// read the configuration
	err := cfg.Read(&config, "")
	if err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}

	// launch the forwarder
	err = forwarder.Run(&config)
	if err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}

}
