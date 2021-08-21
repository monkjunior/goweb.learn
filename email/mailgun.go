package email

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

const (
	// TODO: make this configurable
	resetBaseURL   = "http://127.0.0.1:8080/reset"
	welcomeSubject = "Welcome to Goweb.learn!"
	welcomeText    = `Hi there!

Welcome to Goweb.learn! We really hope you enjoy using our application!

Best,
Ted`
	welcomeHTML = `Hi there!<br/>
<br/>
Welcome to Goweb.learn! We really hope you enjoy using our application!<br/>
<br/>
Best,<br/>
Ted`
	resetSubject  = "Instructions for resetting your password."
	resetTextTmpl = `Hi there!
It appears that you have requested a password reset. If this was you, please follow the link bel
%s
If you are asked for a token, please use the following value:
%s
If you didn't request a password reset you can safely ignore this email and your account will no
Best,
Goweb Learn Support
`
	resetHTMLTmpl = `Hi there!<br/>
<br/>
It appears that you have requested a password reset. If this was you, please follow the link bel
<br/>
<a href="%s">%s</a><br/>
<br/>
If you are asked for a token, please use the following value:<br/>
<br/>
%s<br/>
<br/>
If you didn't request a password reset you can safely ignore this email and your account will no
<br/>
Best,<br/>
Goweb Learn Support<br/>
`
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

func (c *Client) ResetPw(toEmail, token string) error {
	v := url.Values{}
	v.Set("token", token)
	resetUrl := resetBaseURL + "?" + v.Encode()
	resetText := fmt.Sprintf(resetTextTmpl, resetUrl, token)

	message := c.mg.NewMessage(c.sender, resetSubject, resetText, toEmail)
	resetHTML := fmt.Sprintf(resetHTMLTmpl, resetUrl, resetUrl, token)
	message.SetHtml(resetHTML)

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
