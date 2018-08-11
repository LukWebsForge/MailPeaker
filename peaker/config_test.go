package peaker

import (
	testAssert "github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test_parseConfig(t *testing.T) {
	assert := testAssert.New(t)

	_, err := parseConfig("thisIsaWrongPath")
	assert.Error(err)

	// config/wrong{i}.toml

	for i := 1; i <= 2; i++ {
		path := "../test/config/wrong" + strconv.Itoa(i) + ".toml"
		_, err = parseConfig(path)

		if !strings.Contains(err.Error(), "The system cannot find the file specified") {
			assert.Error(err)
		} else {
			t.Fatalf("No testing asset file: '%s'", path)
		}

	}

	// config/correct1.toml

	config, err := parseConfig("../test/config/correct1.toml")
	assert.NoError(err)

	assert.Equal(&Config{
		Out: &ServerConfig{
			Server:   "smtp.example.com",
			Port:     587,
			Email:    "hey@example.com",
			Password: "123456",
		},
		In: []AccountConfig{
			{
				Name: "example",
				Server: &ServerConfig{
					Server:   "imap.example.com",
					Port:     993,
					Email:    "hey@example.com",
					Password: "123456",
				},
				Mailboxes: []string{"INBOX", "2. INBOX"},
				Template:  "de",
			},
		},
		Interval: 300 * time.Second,
		Template: "en",
		Dev:      true,
	}, config)

	// config/correct2.toml

	config, err = parseConfig("../test/config/correct2.toml")
	assert.NoError(err)

	assert.Equal(&Config{
		Out: &ServerConfig{
			Server:   "smtp.example.com",
			Port:     587,
			Email:    "hey@example.com",
			Password: "123456",
		},
		In: []AccountConfig{
			{
				Name: "example",
				Server: &ServerConfig{
					Server:   "imap.example.com",
					Port:     993,
					Email:    "hey@example.com",
					Password: "123456",
				},
				Mailboxes: []string{"INBOX"},
				Template:  "en",
			},
		},
		Interval: 300 * time.Second,
		Template: "en",
		Dev:      false,
	}, config)

}
