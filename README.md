# Mailifica CLI (`mailifica`)

> Official Command-Line Interface for the **Mailifica** Transactional & Marketing Email Infrastructure Platform.

[![Go Report Card](https://goreportcard.com/badge/github.com/mailifica/mailifica-cli)](https://goreportcard.com/report/github.com/mailifica/mailifica-cli)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)

---

## ⚡ Features

- 🔑 **Quick Authentication**: Interactive login (`mailifica login`) saved to `~/.mailifica/config.json`.
- ✉️ **Send & Inspect Emails**: Send emails with HTML, text, attachments, templates, variables, tags, and scheduled delivery.
- 🌐 **Domain & DNS Management**: Add domains, list status, inspect SPF/DKIM/DMARC records, and trigger verification.
- 🛡️ **API Key Management**: Create and revoke workspace API keys.
- 🪝 **Built-in Webhook Tunnel Listener**: Forward webhooks in real time to local applications (`mailifica webhooks listen --forward-to http://localhost:8000/mailifica/webhook`) with live traffic inspector.
- 📄 **JSON Output**: Clean machine-readable output with `--json` for scripting and CI/CD pipelines.

---

## 📦 Installation

### Via `go install`:

```bash
go install github.com/mailifica/mailifica-cli@latest
```

### From Source:

```bash
git clone https://github.com/mailifica/mailifica-cli.git
cd mailifica-cli
make build
# Binary created at ./bin/mailifica
```

---

## 🚀 Quick Start

### 1. Authenticate

```bash
mailifica login
```

Or configure via environment variables:

```bash
export MAILIFICA_API_KEY="ma_live_your_api_key_here"
export MAILIFICA_BASE_URL="https://api.mailifica.com/v1"
```

### 2. Verify Status

```bash
mailifica whoami
```

---

## 📖 Command Reference

### Emails

```bash
# Send a plain text or HTML email
mailifica emails send \
  --to user@example.com \
  --from hello@yourdomain.com \
  --subject "Welcome to our platform" \
  --html "<h1>Welcome!</h1><p>Thanks for joining us.</p>" \
  --tag category=onboarding

# Send email using template and variables
mailifica emails send \
  --to user@example.com \
  --from hello@yourdomain.com \
  --template tmpl_welcome_01 \
  --variables '{"name":"Alex","plan":"Pro"}'

# Send with attachment
mailifica emails send \
  --to user@example.com \
  --from billing@yourdomain.com \
  --subject "Monthly Invoice" \
  --text "Please find your invoice attached." \
  --attach ./invoice.pdf

# Get email delivery details
mailifica emails get email_abc12345

# List recently sent emails
mailifica emails list
```

---

### Domains & DNS

```bash
# List all sending domains
mailifica domains list

# Add a new domain
mailifica domains add notifications.yourcompany.com

# Get DNS verification records (SPF, DKIM, DMARC, MX)
mailifica domains get dom_abc12345

# Trigger DNS verification check
mailifica domains verify dom_abc12345

# Delete a domain
mailifica domains delete dom_abc12345
```

---

### API Keys

```bash
# List workspace API keys
mailifica api-keys list

# Create a new API key
mailifica api-keys create --name "Production Server" --permission full_access

# Revoke an API key
mailifica api-keys revoke key_12345
```

---

### Templates

```bash
# List available email templates
mailifica templates list

# View template content
mailifica templates get tmpl_welcome_01
```

---

### Local Webhook Tunnel Listener

Forward incoming Mailifica webhooks directly to your local development backend:

```bash
mailifica webhooks listen --forward-to http://localhost:8000/mailifica/webhook
```

Sample output:
```text
  __  __       _ _ _  __ _           
 |  \/  |     (_) (_)/ _(_)          
 | \  / | __ _ _| |_| |_ _  ___ __ _ 
 | |\/| |/ _` | | |  _| |/ __/ _` |
 | |  | | (_| | | | | | | | (_| (_| |
 |_|  |_|\__,_|_|_|_|_| |_|\___\__,_|

● Webhook listener ready!
  Listening on: http://localhost:4000
  Forwarding to: http://localhost:8000/mailifica/webhook

14:22:01 ➔ WEBHOOK RECEIVED email.delivered
  Event ID: evt_64fae109
  Forwarded: [200 OK] ➔ http://localhost:8000/mailifica/webhook (18ms)
```

---

## 📄 License

MIT License. See [LICENSE](LICENSE) for details.
