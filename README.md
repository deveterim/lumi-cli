# Lumi CLI (`lumi`)

> Official Command-Line Interface for the **Lumi** Transactional & Marketing Email Infrastructure Platform.

[![Go Report Card](https://goreportcard.com/badge/github.com/deveterim/lumi-cli)](https://goreportcard.com/report/github.com/deveterim/lumi-cli)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)

---

## ⚡ Features

- 🔑 **Quick Authentication**: Interactive login (`lumi login`) saved to `~/.lumi/config.json`.
- ✉️ **Send & Inspect Emails**: Send emails with HTML, text, attachments, templates, variables, tags, and scheduled delivery.
- 🌐 **Domain & DNS Management**: Add domains, list status, inspect SPF/DKIM/DMARC records, and trigger verification.
- 🛡️ **API Key Management**: Create and revoke workspace API keys.
- 🪝 **Built-in Webhook Tunnel Listener**: Forward webhooks in real time to local applications (`lumi webhooks listen --forward-to http://localhost:8000/lumi/webhook`) with live traffic inspector.
- 📄 **JSON Output**: Clean machine-readable output with `--json` for scripting and CI/CD pipelines.

---

## 📦 Installation

### Via `go install`:

```bash
go install github.com/deveterim/lumi-cli@latest
```

### From Source:

```bash
git clone https://github.com/deveterim/lumi-cli.git
cd lumi-cli
make build
# Binary created at ./bin/lumi
```

---

## 🚀 Quick Start

### 1. Authenticate

```bash
lumi login
```

Or configure via environment variables:

```bash
export LUMI_API_KEY="lu_live_your_api_key_here"
export LUMI_BASE_URL="https://api.lumi.email/v1"
```

### 2. Verify Status

```bash
lumi whoami
```

---

## 📖 Command Reference

### Emails

```bash
# Send a plain text or HTML email
lumi emails send \
  --to user@example.com \
  --from hello@yourdomain.com \
  --subject "Welcome to our platform" \
  --html "<h1>Welcome!</h1><p>Thanks for joining us.</p>" \
  --tag category=onboarding

# Send email using template and variables
lumi emails send \
  --to user@example.com \
  --from hello@yourdomain.com \
  --template tmpl_welcome_01 \
  --variables '{"name":"Alex","plan":"Pro"}'

# Send with attachment
lumi emails send \
  --to user@example.com \
  --from billing@yourdomain.com \
  --subject "Monthly Invoice" \
  --text "Please find your invoice attached." \
  --attach ./invoice.pdf

# Get email delivery details
lumi emails get email_abc12345

# List recently sent emails
lumi emails list
```

---

### Domains & DNS

```bash
# List all sending domains
lumi domains list

# Add a new domain
lumi domains add notifications.yourcompany.com

# Get DNS verification records (SPF, DKIM, DMARC, MX)
lumi domains get dom_abc12345

# Trigger DNS verification check
lumi domains verify dom_abc12345

# Delete a domain
lumi domains delete dom_abc12345
```

---

### API Keys

```bash
# List workspace API keys
lumi api-keys list

# Create a new API key
lumi api-keys create --name "Production Server" --permission full_access

# Revoke an API key
lumi api-keys revoke key_12345
```

---

### Templates

```bash
# List available email templates
lumi templates list

# View template content
lumi templates get tmpl_welcome_01
```

---

### Local Webhook Tunnel Listener

Forward incoming Lumi webhooks directly to your local development backend:

```bash
lumi webhooks listen --forward-to http://localhost:8000/lumi/webhook
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
  Forwarding to: http://localhost:8000/lumi/webhook

14:22:01 ➔ WEBHOOK RECEIVED email.delivered
  Event ID: evt_64fae109
  Forwarded: [200 OK] ➔ http://localhost:8000/lumi/webhook (18ms)
```

---

## 📄 License

MIT License. See [LICENSE](LICENSE) for details.
