package email

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

const (
	welcomeSubject = "Welcome to Goweb.learn!"
	welcomeText    = `Hi there!

Welcome to Goweb.learn! We really hope you enjoy using our application!

Best,
Jon`
	welcomeHTML = `Hi there!<br/>
<br/>
Welcome to Goweb.learn! We really hope you enjoy using our application!<br/>
<br/>
Best,<br/>
Jon`
)

type ClientConfig func(*Client)

func WithSender(name, email string) ClientConfig {
	return func(c *Client) {
		c.sender = buildEmail(name, email)
	}
}

func WithMailgun(domain, apiKey string) ClientConfig {
	return func(c *Client) {
		c.mg = mailgun.NewMailgun(domain, apiKey)
	}
}

func NewClient(opts ...ClientConfig) *Client {
	client := Client{
		// Set a default from email address
		sender: "support@goweb.com",
	}
	for _, opt := range opts {
		opt(&client)
	}
	return &client
}

type Client struct {
	sender string
	mg     mailgun.Mailgun
}

func (c *Client) Welcome(toName, toEmail string) error {
	message := c.mg.NewMessage(c.sender, welcomeSubject, welcomeText, buildEmail(toName, toEmail))
	message.SetHtml(welcomeHTML)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := c.mg.Send(ctx, message)
	return err
}

func buildEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
