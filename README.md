# SMTP2HTTP - SMTP to HTTP Gateway

**English**🌎 | [**简体中文**🀄](./README_zh.md)

## Major Update

Added scheduled task functionality supporting cron expressions to periodically fetch API responses and send emails. Supports JSON response parsing with customizable templates.

## Project Description

An SMTP-to-HTTP gateway service that helps developers:
- Eliminate hard-coded SMTP configurations
- Send emails dynamically via REST API
- Visually configure multiple SMTP providers
- Utilize scheduled tasks to reduce email-related code

## Key Features
- 🚀 Web UI for SMTP server configuration
- 📦 Supports GET/POST/JSON request formats
- 🔒 Token-based API authentication
- 📎 Multi-file attachment support (Base64/Form-data)
- 🐳 Ready-to-use Docker image
- 📊 Email sending history & statistics
- 📅 Cron-based scheduled tasks
- 📧 Email templating support

## Community
Join our WeChat group or Telegram channel for discussions:

WeChat Group:

<img src="https://github.com/chuccp/smtp2http/blob/main/image/WeChat.png?raw=true" alt="WebChat" width="200">

Telegram:

https://t.me/+JClG9_DojaM0ZGE1

## Quick Start

### Direct Execution

```bash
# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/chuccp/smtp2http/releases/latest/download/smtp2http-windows-amd64.tar.gz" -OutFile "smtp2http-windows-amd64.tar.gz"
tar -zxvf smtp2http-windows-amd64.tar.gz
.\smtp2http.exe

# Linux
wget https://github.com/chuccp/smtp2http/releases/latest/download/smtp2http-linux-amd64.tar.gz
tar -zxvf smtp2http-linux-amd64.tar.gz
chmod +x smtp2http
./smtp2http
```

### Docker
```bash

docker pull cooge123/smtp2http

docker run -p 12566:12566 -p 12567:12567 -it --rm cooge123/smtp2http

```

## Configuration
Initial startup generates `config.ini`:

```ini
[core]
init = true
cachePath = .cache
dbType = sqlite

[sqlite]
filename = d-mail.db

[manage]
port = 12566
username = 111111
password = 111111
webPath = web

[api]
port = 12566

[mysql]
host = 127.0.0.1
port = 3306
dbname = d_mail
charset = utf8
username = root
password = 123456
```

## API Documentation

### Send Email
`POST /sendMail`

**Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| token | string | Yes | Auth token |
| subject | string | No | Email subject |
| content | string | Yes | Email body |
| recipients | []string | No | CC recipients |
| files | []File | No | Attachments |

**Success Response**:
```json
ok
```

### Examples
**JSON with Attachment**
```bash
curl -X POST 'http://127.0.0.1:12567/sendMail' \
--header 'Content-Type: application/json' \
--data-raw '{
  "token": "{{token}}",
  "subject": "test",
  "content": "this is a test",
  "recipients": ["ops@example.com"],
  "files": [{
    "name": "alert.log",
    "data": "{{base64_content}}"
  }]
}'
```

**Form-data with Multiple Files**
```bash
curl -X POST 'http://127.0.0.1:12567/sendMail' \
--form 'token={{token}}' \
--form 'subject=test' \
--form 'content=this is a test' \
--form 'recipients=finance@example.com,sales@example.com' \
--form 'files=@"/data/reports/sales.pdf"' \
--form 'files=@"/data/reports/expenses.xlsx"'
```

**GET Request**
```bash
curl 'http://127.0.0.1:12567/sendMail?token={{token}}&subject=test&content=this%20is%20a%20test&recipients=aaa@mail.com,bbb@mail.com'
```

## Building
Requires building the frontend [d-mail-view](https://github.com/chuccp/d-mail-view) first.
```