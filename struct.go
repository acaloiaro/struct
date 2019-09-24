package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type fields []string

var fieldsFlag fields
var sepFlag string
var outFlag string

var in string
var fieldsProvided int

func (f *fields) String() string {
	return fmt.Sprint(*f)
}

func (f *fields) Set(value string) error {
	if len(*f) > 0 {
		return errors.New("fields switch already set")
	}

	for _, field := range strings.Split(value, ",") {
		*f = append(*f, field)
	}

	return nil
}

const (
	defaultFieldSeparator = " "
	defaultOutput         = "string"
	fieldsUsage           = "comma-separated list of fields corresponding to the input data's fields\n\n" +
		"Example\n\n" +
		"$ echo \"foo1 bar1\" | struct -fields foo,bar\n" +
		"$ foo:foo1 bar:bar1\n"
	outputUsage    = "the desired output format: json or string"
	separatorUsage = "the character that separates fields in the input data"
)

func init() {

	flag.Var(&fieldsFlag, "fields", fieldsUsage)
	flag.StringVar(&sepFlag, "separator", defaultFieldSeparator, separatorUsage)
	flag.StringVar(&outFlag, "output", defaultOutput, outputUsage)
	setupFlags(flag.CommandLine)

	flag.Parse()
}

func setupFlags(f *flag.FlagSet) {
	f.Usage = func() {
		fmt.Println("\nCreate structured data input strings")

		f.PrintDefaults()
	}
}

var sb strings.Builder
var jsonMap = make(map[string]string)

func main() {
	fieldsProvided = len(fieldsFlag)

	scanner := bufio.NewScanner(os.Stdin)

	var i int
	var field, v string

	for scanner.Scan() {
		in = scanner.Text()

		inputFields := strings.Fields(in)

		for i, v = range inputFields {
			if i < fieldsProvided {

				field = fieldsFlag[i]
			} else {
				field = ""
			}

			buildOutput(field, v)
		}

		fmt.Println(out())
	}
}

func out() (out string) {
	if outFlag == "string" {
		out = sb.String()
		sb.Reset()
	} else if outFlag == "json" {
		o, _ := json.Marshal(jsonMap)
		out = string(o)
		for k := range jsonMap {
			delete(jsonMap, k)
		}
	} else {
		out = ""
	}

	return
}

func buildOutput(field, value string) error {
	if outFlag == "string" {
		buildString(field, value)
	} else if outFlag == "json" {
		buildJSON(field, value)
	} else {
		return errors.New("invalid output option")
	}

	return nil
}

func buildString(field, value string) {
	if field != "" {
		sb.WriteString(field)
		sb.WriteString(":")
		sb.WriteString(value)
		sb.WriteString(sepFlag)
	} else {
		sb.WriteString(value)
		sb.WriteString(sepFlag)
	}
}

func buildJSON(field, value string) {
	if field == "" {
		return
	}

	jsonMap[field] = value
}
