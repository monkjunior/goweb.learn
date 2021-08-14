package views

const (
	AlertLvError   = "danger"
	AlertLvWarning = "warning"
	AlertLvInfo    = "info"
	AlertLvSuccess = "success"
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
	Yield interface{}
}
