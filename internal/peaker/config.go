package peaker

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ServerConfig struct {
	Server   string
	Port     int
	Email    string
	Password string
}

type AccountConfig struct {
	Name      string
	Server    *ServerConfig
	Mailboxes []string
	Template  string
}

type Config struct {
	Out      *ServerConfig
	Template string
	In       []*AccountConfig
	Dev      bool
	Interval time.Duration
}

const defaultConfig string = `
[base]
    # Use a bundled template with         "en" or "de"
    # OR use a path to a custom template  "./templates/custom.template"
    template = "en"
	# Whether to enable the dev mode. This starts the reading process instantly
    dev = true
    # The interval how often the peaker should check the mails (in seconds)
    interval = 300 # 60 * 5 => 5 Minutes
	# The outbound SMPT server to send to the mails
    server = "smtp.example.com"
    port = 587
    email = "hey@example.com"
    password = "123456"

[accounts]
    [accounts.example]
        # The IMAP server to check
        server = "imap.example.com"
        port = 993
        email = "hey@example.com"
        password = "123456"
        # Optional, defaults to "INBOX"
        mailboxes = ["INBOX", "2. INBOX"]
        # Optional, defaults to [base.template]
        template = "de"
`

func ReadConfig() (*Config, error) {
	filePath, err := findConfigFile()
	if err != nil {
		return nil, err
	}

	config, err := parseConfig(filePath)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func findConfigFile() (string, error) {
	filePath := "config.toml"

	if envPath := os.Getenv("PEAKER_CONFIG"); envPath != "" {
		// Using the environment variable PEAKER_CONFIG to change the config location
		filePath = envPath
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// The config file doesn't exists, so we copy it
		// File Permissions: 6 => The owner can edit it, 4 => the group can read it, 0 => no one else can read it
		err := ioutil.WriteFile(filePath, []byte(defaultConfig), 0640)
		if err != nil {
			return "", fmt.Errorf("can't write the config file to '%s'. "+
				"Create it manually and start the program again", filePath)
			panic(err)
		} else {
			return "", fmt.Errorf("created a new config file '%s'. "+
				"Please edit it and start the program again", filePath)
		}
	}

	return filePath, nil
}

func parseConfig(filePath string) (*Config, error) {
	// The config file exists (now), so we'll read it
	tree, err := toml.LoadFile(filepath.FromSlash(filePath))
	if err != nil {
		return nil, fmt.Errorf("can't parse the config file '%s': %v", filePath, err)
	}

	if err := checkKeys(tree, "base", "accounts", "base.template", "base.interval"); err != nil {
		return nil, err
	}

	outServerConfig, err := parseServerConfig(tree.Get("base").(*toml.Tree))
	if err != nil {
		return nil, err
	}

	defaultTemplate := tree.Get("base.template").(string)
	configAccounts, err := parseConfigAccounts(tree.Get("accounts").(*toml.Tree), defaultTemplate)
	if err != nil {
		return nil, err
	}

	config := Config{
		Out:      outServerConfig,
		In:       configAccounts,
		Interval: time.Duration(tree.Get("base.interval").(int64)) * time.Second,
		Template: defaultTemplate,
		Dev:      tree.GetDefault("base.dev", false).(bool),
	}

	return &config, nil
}

func parseConfigAccounts(tree *toml.Tree, defaultTemplate string) ([]*AccountConfig, error) {
	configs := make([]*AccountConfig, len(tree.Keys()))

	for index, key := range tree.Keys() {
		element := tree.Get(key).(*toml.Tree)

		serverConfig, err := parseServerConfig(element)
		if err != nil {
			return nil, err
		}

		nameSplit := strings.Split(key, ".")
		name := nameSplit[len(nameSplit)-1]

		config := AccountConfig{
			Name:      name,
			Server:    serverConfig,
			Template:  element.GetDefault("template", defaultTemplate).(string),
			Mailboxes: parseMailboxes(element),
		}

		configs[index] = &config
	}

	return configs, nil
}

func parseServerConfig(tree *toml.Tree) (*ServerConfig, error) {
	if err := checkKeys(tree, "server", "port", "email", "password"); err != nil {
		return nil, err
	}

	config := ServerConfig{
		Server:   tree.Get("server").(string),
		Port:     int(tree.Get("port").(int64)),
		Email:    tree.Get("email").(string),
		Password: tree.Get("password").(string),
	}

	return &config, nil
}

func parseMailboxes(tree *toml.Tree) []string {
	mailboxesRaw, ok := tree.Get("mailboxes").([]interface{})
	mailboxes := make([]string, len(mailboxesRaw))

	if ok && mailboxesRaw != nil {
		for indexMail, mailbox := range mailboxesRaw {
			mailboxes[indexMail] = mailbox.(string)
		}
	}

	if len(mailboxes) == 0 {
		mailboxes = append(mailboxes, "INBOX")
	}

	return mailboxes
}

func checkKeys(tree *toml.Tree, entries ...string) error {
	for _, element := range entries {

		if tree.Has(element) {
			continue
		}

		return fmt.Errorf("the config file has no entry for the key '%s'", element)
	}

	return nil
}
