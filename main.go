package main

import (
	"flag"
	cfg "github.com/Bourne-ID/beats-forwarder/config"
	"github.com/Bourne-ID/beats-forwarder/forwarder"
	"github.com/Sirupsen/logrus"
	"os"
	"path"
	"log"
)

var config = cfg.Config{}

func main() {

	flag.Parse()
	debug := flag.Lookup("d").Value.String()

	if debug == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}
	ex, err := os.Executable()
	if err != nil { log.Fatal(err) }
	dir := path.Dir(ex)

	var f *os.File
	if f, err = os.Create(dir + "/log.log"); err != nil {
		logrus.Error("log file creation failed", err)
		os.Exit(1)
	}
	logrus.SetOutput(f)
	defer func(f *os.File) {
		f.Sync() //look here
		f.Close()
	}(f)

	// read the configuration
	err = cfg.Read(&config)
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
