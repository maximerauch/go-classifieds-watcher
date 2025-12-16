package email

import (
	"context"
	"fmt"
	"strings"

	"gopkg.in/gomail.v2"

	"github.com/maximerauch/go-classifieds-watcher/internal/config"
	"github.com/maximerauch/go-classifieds-watcher/internal/core"
)

type EmailNotifier struct {
	cfg config.EmailConfig
}

func NewEmailNotifier(cfg config.EmailConfig) *EmailNotifier {
	return &EmailNotifier{cfg: cfg}
}

// Send sends a single email for a specific item.
func (n *EmailNotifier) Send(ctx context.Context, item core.Item) error {
	// Respect context cancellation
	if ctx.Err() != nil {
		return ctx.Err()
	}

	m := gomail.NewMessage()
	m.SetHeader("From", n.cfg.From)
	m.SetHeader("To", n.cfg.To...)

	// Subject: e.g., "🔔 New Dog: Rex - Male"
	subject := fmt.Sprintf("🔔 New Item: %s", item.Title)
	m.SetHeader("Subject", subject)

	// Body construction
	body := n.buildBody(item)
	m.SetBody("text/html", body)

	// SMTP Configuration
	d := gomail.NewDialer(n.cfg.SMTPHost, n.cfg.SMTPPort, n.cfg.SMTPUser, n.cfg.SMTPPassword)

	// Send
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email for item %s: %w", item.ID, err)
	}

	return nil
}

func (n *EmailNotifier) buildBody(item core.Item) string {
	var sb strings.Builder

	sb.WriteString("<h2>New Item Discovered!</h2>")
	sb.WriteString("<ul>")

	// Title
	sb.WriteString(fmt.Sprintf("<li><strong>Title:</strong> %s</li>", item.Title))

	// Price (if applicable)
	if item.Price > 0 {
		sb.WriteString(fmt.Sprintf("<li><strong>Price:</strong> %.2f %s</li>", item.Price, item.Currency))
	}

	// Description
	sb.WriteString(fmt.Sprintf("<li><strong>Details:</strong> %s</li>", item.Description))
	sb.WriteString("</ul>")

	// Call-to-Action Button
	sb.WriteString(fmt.Sprintf(`
		<br/>
		<a href="%s" style="background-color: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; font-family: Arial, sans-serif;">
			View Item on Website
		</a>
		<br/><br/>
	`, item.Url))

	sb.WriteString("<p style='color: #888; font-size: 12px;'>Sent by Go Watcher</p>")

	return sb.String()
}
