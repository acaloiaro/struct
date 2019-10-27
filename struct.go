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

var (
	fieldsFlag fields // a comma-separated list of fields provided by the user
	outFlag    string // the format (either string or json) in which to structure data
	sepFlag    string // a string on which to split lines of input
)

type fields []string

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

func init() {

	flagset := flag.CommandLine
	flagset.Var(&fieldsFlag, "fields", fieldsUsage)
	flagset.StringVar(&sepFlag, "separator", defaultFieldSeparator, separatorUsage)
	flagset.StringVar(&outFlag, "output", defaultOutput, outputUsage)

	flagset.Usage = func() {
		fmt.Println("\nCreate structured output from unstructured input")

		flagset.PrintDefaults()
	}

	flagset.Parse(os.Args[1:])
}

func setupFlags(f *flag.FlagSet) {

}

var sb strings.Builder
var jsonMap = make(map[string]string)

func main() {
	fieldsProvided := len(fieldsFlag)

	scanner := bufio.NewScanner(os.Stdin)

	var i int
	var field, in, v string

	// perform a line-wise scan over stdin until EOF
	for scanner.Scan() {
		in = scanner.Text()

		// if the separator string is overridden, split on it instead of the default of splitting on space
		var inputFields []string
		if sepFlag != defaultFieldSeparator {
			inputFields = strings.Split(in, sepFlag)
		} else {
			inputFields = strings.Fields(in)
		}

		// for every field to be retained from the input, build a representation of it in the output
		for i, v = range inputFields {
			if i < fieldsProvided {

				field = fieldsFlag[i]
			} else {
				field = ""
			}

			if err := buildOutput(field, v); err != nil {
				fmt.Println(err.Error())
				return
			}
		}

		fmt.Println(out())
	}
}

func out() (out string) {

	switch outFlag {
	case "string":
		out = sb.String()
		sb.Reset()
	case "json":
		o, _ := json.Marshal(jsonMap)
		out = string(o)
		for k := range jsonMap {
			delete(jsonMap, k)
		}
	default:
		out = ""
	}

	return
}

func buildOutput(field, value string) (err error) {
	switch outFlag {
	case "string":
		buildString(field, value)
	case "json":
		buildJSON(field, value)
	default:
		err = fmt.Errorf("invalid -output option: %s", outFlag)
		return
	}

	return
}

func buildString(field, value string) {
	if field != "" {
		sb.WriteString(field)
		sb.WriteString(":")
		sb.WriteString(value)
		sb.WriteString(defaultFieldSeparator)
	} else {
		sb.WriteString(value)
		sb.WriteString(defaultFieldSeparator)
	}
}

func buildJSON(field, value string) {
	if field == "" {
		return
	}

	jsonMap[field] = value
}
