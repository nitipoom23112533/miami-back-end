package mg

import (
	"github.com/mailgun/mailgun-go/v4"
)

const (
	domainName    = "domainName"
	privateAPIKey = "privateAPIKey"
)

// Client mailgun client instance
var Client *mailgun.MailgunImpl

// InitMailGunClient func
func InitMailGunClient() {
	Client = mailgun.NewMailgun(domainName, privateAPIKey)
}