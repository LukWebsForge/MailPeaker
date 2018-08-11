package peaker

import (
	"bytes"
	testAssert "github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestFindTemplate(t *testing.T) {
	assert := testAssert.New(t)

	// Bundled

	_, err := FindTemplate("en")
	assert.NoError(err)

	_, err = FindTemplate("de")
	assert.NoError(err)

	_, err = FindTemplate("gloryToArstotzka")
	assert.Error(err)

	// External files

	template, err := FindTemplate("../test/template/correct.txt.tmpl")
	assert.NoError(err)

	buffer := bytes.NewBufferString("")
	err = template.Execute(buffer, map[string]interface{}{
		"Account": "hey@example.com",
		"Email":   "hello@some.one",
	})
	wantedText := "New mail on hey@example.com\n\nHey you've got a new mail from hello@some.one"

	assert.NoError(err)
	assert.Equal(wantedText, strings.Replace(buffer.String(), "\r", "", -1))

	_, err = FindTemplate("../test/template/wrong.txt.tmpl")
	assert.Error(err)

	_, err = FindTemplate("./coolItDoesntExists")
	assert.Error(err)
}
