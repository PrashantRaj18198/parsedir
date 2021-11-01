package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

type FileInfo struct {
	// Path is the path to the file, can be template string
	Path string
	// Content in the content stored in the file, can be template string
	Content string
}

var Seperator = "\n---seperator---\n"

type MissingVariable string

const (
	// MissingVariableOff will neither show error nor warn when path parsing fails
	MissingVariableOff MissingVariable = "off"
	// MissingVariableWan will warn when the path parsing fails (logs to stderr)
	MissingVariableWarn MissingVariable = "warn"
	// MissingVariable will return error if the path parsing fails
	MissingVariableError MissingVariable = "error"
)

var (
	// MissingVariableCurr can be updated to user choice when treating errors
	// when unable to parse path variable. You can set them to "warn", "error", or "off".
	// Default is set below
	MissingVariableCurr MissingVariable = "warn"
)

func RecurseThroughDir(base string) ([]string, error) {
	files := []string{}
	fss, err := ioutil.ReadDir(base)
	if err != nil {
		return files, err
	}
	for _, fs := range fss {
		fullName := filepath.Join(base, fs.Name())
		if fs.IsDir() {
			out, err := RecurseThroughDir(fullName)
			if err != nil {
				return files, err
			}
			files = append(files, out...)
		} else {
			files = append(files, fullName)
		}
	}
	return files, nil
}

func ReadAllFiles(files []string) ([]*FileInfo, error) {
	fileInfos := []*FileInfo{}
	for _, path := range files {

		bdata, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		fileInfos = append(fileInfos, &FileInfo{
			Path:    path,
			Content: string(bdata),
		})
	}
	return fileInfos, nil
}

var re = regexp.MustCompile(`(?m)(\{\{\s*range\s*[^\s]+?\s*?\}\})`)

func PopulateAllFiles(files []*FileInfo, data interface{}) ([]*FileInfo, error) {
	out := []*FileInfo{}
	for _, f := range files {
		fmt.Fprintf(os.Stdout, "Parsing path '%s'\n", f.Path)
		matches := re.FindAllString(f.Path, -1)
		currOut, err := populateFile(f, data, matches...)
		if err != nil && MissingVariableCurr == MissingVariableError {
			return out, err
		}
		out = append(out, currOut...)

	}
	return out, nil
}

// populateFile returns single or multiple file in a slice of fileinfo
// if the file is a ranged file, pass a string to the rangedString
// otherwise, pass a nil value inplace of rangedString pointer
func populateFile(f *FileInfo, data interface{}, rangedStrings ...string) ([]*FileInfo, error) {
	out := []*FileInfo{}
	filePathWithoutRange := re.ReplaceAllString(f.Path, "")

	var err error
	var parsedPath string
	if len(rangedStrings) > 0 {
		parsedPath, err = parseRangedString(filePathWithoutRange, data, rangedStrings...)
	}
	if len(rangedStrings) == 0 {
		parsedPath, err = parseTemplate(filePathWithoutRange, data)
	}
	if err != nil {
		if MissingVariableCurr == MissingVariableWarn {
			fmt.Fprintf(os.Stderr, "[Warn] Unable to parse filepath %s. Error: %v.", f.Path, err)
		}
		if MissingVariableCurr == MissingVariableError {
			fmt.Fprintf(os.Stderr, "[Error] Unable to parse filepath %s. Error: %v.", f.Path, err)
		}
		return out, err
	}
	var parsedContent string
	if len(rangedStrings) > 0 {

		parsedContent, err = parseRangedString(f.Content, data, rangedStrings...)
	}
	if len(rangedStrings) == 0 {
		parsedContent, err = parseTemplate(f.Content, data)
	}
	if err != nil {
		return out, err
	}
	out, _ = getFileInfoFromParsed(parsedPath, parsedContent)

	return out, err
}

func parseRangedString(tmplString string, data interface{}, rangedStrings ...string) (string, error) {

	rangedStringsJoined := strings.Join(rangedStrings, "\n")
	ends := []string{}
	for range rangedStrings {
		ends = append(ends, "{{end}}")
	}
	endsJoined := strings.Join(ends, "\n")

	parseString := fmt.Sprintf("%s\n%s%s\n%s", rangedStringsJoined, tmplString, Seperator, endsJoined)
	return parseTemplate(parseString, data)
}

func parseTemplate(tmplString string, data interface{}) (string, error) {
	tmpl := template.New("template")
	tmpl, err := tmpl.Parse(tmplString)
	if err != nil {
		return "", err
	}
	writer := bytes.NewBufferString("")
	err = tmpl.Execute(writer, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing. `%v`", tmplString)
		return "", err
	}
	return writer.String(), nil
}

func getFileInfoFromParsed(path string, content string) ([]*FileInfo, error) {
	out := []*FileInfo{}
	paths := strings.Split(path, Seperator)
	contents := strings.Split(content, Seperator)
	if len(paths) != len(contents) {
		fmt.Printf("wrong path lens")
	}
	for i := range paths {
		f := FileInfo{
			Path:    strings.TrimRight(strings.TrimLeft(paths[i], "\n"), "\n"),
			Content: strings.TrimRight(strings.TrimLeft(contents[i], "\n"), "\n"),
		}
		if f.Path == "" && f.Content == "" {
			continue
		}
		out = append(out, &f)

	}
	return out, nil

}
