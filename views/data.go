package views

import "github.com/monkjunior/goweb.learn/models"

const (
	AlertLvError   = "danger"
	AlertLvWarning = "warning"
	AlertLvInfo    = "info"
	AlertLvSuccess = "success"

	AlertMsgGeneric = "Something went wrong. Please try again, and contact us if the problem persists."
)

// Alert is used to render bootstrap alert messages in templates.
type Alert struct {
	Level   string
	Message string
}

// Data is the top level structure that views expect data
// to come in.
type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}

func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(PublicError); ok {
		d.AlertError(pErr.Public())
	} else {
		d.AlertError(AlertMsgGeneric)
	}
}

func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvError,
		Message: msg,
	}
}

type PublicError interface {
	error
	Public() string
}
