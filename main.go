package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type options struct {
	verbose     bool
	enumToCheck bool
	srcMySQL    string
	dstCrDB     string
}

var opts options

type verboseLogger struct{}

func (v *verboseLogger) Log(err error) {
	if opts.verbose {
		log.Println(err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags)

	// TODO: Look at using cobra cmd for this.
	flag.Usage = usage
	flag.BoolVar(&opts.verbose, "v", false, "verbose logging mode")
	flag.BoolVar(&opts.enumToCheck, "enum-to-check", false, "convert enums to check constraints")
	flag.Parse()

	if len(flag.Args()) < 2 {
		usage()
		os.Exit(1)
	}

	opts.srcMySQL = flag.Args()[0]
	if len(opts.srcMySQL) == 0 {
		usage()
		os.Exit(1)
	}

	opts.dstCrDB = flag.Args()[1]
	if len(opts.dstCrDB) == 0 {
		usage()
		os.Exit(1)
	}

	logger := verboseLogger{}

	srcInfo, err := os.Stat(opts.srcMySQL)
	if err != nil {
		logger.Log(fmt.Errorf("could not stat source file %q: %s", opts.srcMySQL, err))

		fmt.Println("Hmm, couldn't open the source mysql file for reading. Make sure the path and permissions are correct.")
		os.Exit(1)
	}

	// TODO: Does reading and parsing line-by-line make more sense?
	srcData, err := ioutil.ReadFile(opts.srcMySQL)
	if err != nil {
		logger.Log(fmt.Errorf("could not read file %q: %s", opts.srcMySQL, err))

		fmt.Println("Hmm, couldn't read the source mysql file for reading. Make sure the path and permissions are correct.")
		os.Exit(1)
	}

	roacher := newRoacher(srcData)

	crdbData, err := roacher.roach()
	if err != nil {
		logger.Log(fmt.Errorf("could not roach mysql data: %s", err))

		fmt.Println("Damn, couldn't convert the mysql dump to crdb. Make sure the source file was prepared correctly.")
		os.Exit(1)
	}

	err = ioutil.WriteFile(opts.dstCrDB, crdbData, srcInfo.Mode())
	if err != nil {
		logger.Log(fmt.Errorf("could not save the crdb data to file %q: %s", opts.dstCrDB, err))

		fmt.Println("Oh dear, couldn't write the crdb data. Make sure you have the correct permissions.")
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Usage: radroach [FLAGS...] SOURCE_MYSQL_DUMP DESTINATION_CRDB_DUMP")
	flag.PrintDefaults()
}
