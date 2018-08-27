package peaker

import (
	"errors"
	"fmt"
	"text/template"
	"unicode"
)

func FindTemplate(path string) (*template.Template, error) {
	packaged := unicode.IsLetter([]rune(path)[0])

	var readyTmpl *template.Template
	templateName := "peaker-" + path

	if packaged {
		// Searching for the bundled template
		content, err := bundledTemplateFor(path)
		if err != nil {
			// Bundled template not found -> err
			return nil, err
		}
		// Bundled template found -> Parsing
		parse, err := template.New(templateName).Parse(content)
		if err != nil {
			// Error while parsing
			return nil, fmt.Errorf("can't open %v", err)
		}
		// Everything is fine
		readyTmpl = parse
	} else {
		// Parsing the file
		parse, err := template.ParseFiles(path) // Maybe append nil as parameter
		if err != nil {
			// Error while parsing -> More errors
			return nil, fmt.Errorf("can't open template file '%s': %v", path, err)
		}
		// Everything is fine
		readyTmpl = parse
	}

	return readyTmpl, nil
}

// The first line of the template is the subject of the email
const enTemplate = `

`

const deTemplate = `

`

func bundledTemplateFor(name string) (string, error) {
	switch name {
	case "en":
		return enTemplate, nil
	case "de":
		return deTemplate, nil
	default:
		return "", errors.New("no bundled template '" + name + "' found")
	}
}
