package peaker

import (
	"fmt"
	"github.com/emersion/go-imap"
	iClient "github.com/emersion/go-imap/client"
)

func ReadMails(config *ServerConfig) {
	client, err := connect(config)
	if err != nil {
		panic(err)
	}

	// Don't forget to logout
	defer client.Logout()

	mailboxes, err := listMailboxes(client)
	if err != nil {
		panic(err)
	}

	for _, m := range mailboxes {
		println("Mailbox: " + m.Name)
	}

	mails, err := listMails("INBOX", client)
	if err != nil {
		panic(err)
	}

	for _, m := range mails {
		println("Mail: " + m.Envelope.Subject + " at " + m.Envelope.Date.String())
		println(m.Body)
	}

}

func connect(config *ServerConfig) (client *iClient.Client, err error) {
	// Connect to the server
	serverURL := fmt.Sprintf("%s:%d", config.Server, config.Port)

	client, err = iClient.DialTLS(serverURL, nil)
	if err != nil {
		return nil, err
	}

	// Logging in
	if err := client.Login(config.Email, config.Password); err != nil {
		return nil, err
	}

	return client, nil
}

func listMailboxes(client *iClient.Client) (mailboxes []*imap.MailboxInfo, err error) {
	boxChannel := make(chan *imap.MailboxInfo)
	done := make(chan error, 1)
	go func() {
		done <- client.List("", "*", boxChannel)
	}()

	mailboxes = make([]*imap.MailboxInfo, 0)
	for m := range boxChannel {
		mailboxes = append(mailboxes, m)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return mailboxes, nil
}

func listMails(boxName string, client *iClient.Client) (mails []*imap.Message, err error) {
	mbox, err := client.Select(boxName, false)
	if err != nil {
		return nil, err
	}

	// Checking the last 40 emails
	// Using a uint32, because mbox.Messages also returns one
	rangeMin := uint32(1)
	rangeMax := mbox.Messages
	if mbox.Messages+1 > 40 {
		rangeMin = mbox.Messages - 40
	}

	selectRange := new(imap.SeqSet)
	selectRange.AddRange(rangeMin, rangeMax)

	msgChannel := make(chan *imap.Message, 50)
	done := make(chan error, 1)
	go func() {
		done <- client.Fetch(selectRange, []imap.FetchItem{imap.FetchEnvelope}, msgChannel)
	}()

	messages := make([]*imap.Message, 0)
	for m := range msgChannel {
		messages = append(messages, m)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return messages, nil
}
