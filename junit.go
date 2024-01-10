package main

import (
	"encoding/xml"
	"io"
	"os"
	"path"

	"github.com/bmatcuk/doublestar"
)

type TestCase struct {
	XMLName xml.Name `xml:"testcase"`
	File    string   `xml:"file,attr"`
	Time    float64  `xml:"time,attr"`
}

type TestSuite struct {
	TestCases []TestCase `xml:"testcase"`
}

func loadJUnitXML(reader io.Reader) *TestSuite {
	var suite TestSuite

	decoder := xml.NewDecoder(reader)
	err := decoder.Decode(&suite)
	if err != nil {
		fatalMsg("failed to parse junit xml: %v\n", err)
	}

	return &suite
}

func addFileTimesFromIOReader(fileTimes map[string]float64, reader io.Reader) {
	junitXML := loadJUnitXML(reader)
	printMsg("found %d test cases\n", len(junitXML.TestCases))

	for _, testCase := range junitXML.TestCases {
		filePath := path.Clean(testCase.File)
		fileTimes[filePath] += testCase.Time
	}
}

func getFileTimesFromJUnitXML(fileTimes map[string]float64) {
	if junitXMLPath != "" {
		filenames, err := doublestar.Glob(junitXMLPath)
		if err != nil {
			fatalMsg("failed to match jUnit filename pattern: %v", err)
		}
		for _, junitFilename := range filenames {
			file, err := os.Open(junitFilename)
			if err != nil {
				fatalMsg("failed to open junit xml: %v\n", err)
			}
			printMsg("using test times from JUnit report %s\n", junitFilename)
			addFileTimesFromIOReader(fileTimes, file)
			file.Close()
		}
	} else {
		printMsg("using test times from JUnit report at stdin\n")
		addFileTimesFromIOReader(fileTimes, os.Stdin)
	}
}
