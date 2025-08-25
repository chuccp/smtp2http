# SMTP2HTTP - SMTP to HTTP  Tool

**[English](README.md)** 🌎 | **简体中文** 🀄

## Important Updates

Added scheduled task functionality, which can periodically read other APIs or links and send mail. It also supports template configuration, converting JSON data into text and sending it as an email.

## Project Description

A gateway service that converts SMTP protocol to HTTP interface, helping developers:

- Avoid hardcoding SMTP configurations in code
- Dynamically send emails through REST API
- Visually configure multiple SMTP service providers
- Use scheduled tasks to reduce email configuration in projects

## Key Features

- 🚀 Configure SMTP servers and recipient emails via Web UI
- 📦 Support for GET/POST/JSON request formats
- 🔒 Token-based API access control
- 📎 Multi-file attachment support (Base64 encoding/form upload)
- 🐳 Ready-to-use Docker image
- 📊 Send record query and statistics
- 📅 Support for scheduled tasks, requesting links, and sending emails
- 📧 Email template support

## Quick Start

### Direct Execution

```bash
# Windows
curl -LO https://github.com/chuccp/smtp2http/releases/latest/download/smtp2http-windows-amd64.tar.gz
./smtp2http.exe

# Linux
curl -LO https://github.com/chuccp/smtp2http/releases/latest/download/smtp2http-linux-amd64.tar.gz
chmod +x smtp2http
./smtp2http
```


### Docker Deployment

```bash
docker run -d \
  -p 12566:12566 \
  -p 12567:12567 \
  cooge123/smtp2http:latest
```


## Configuration Instructions

The configuration file  config.ini is generated after the first startup:

```ini
[core]
init      = true   ## Initialization switch, automatically becomes true after first startup
cachePath = .cache  ## Email attachment cache directory
dbType    = sqlite  ## Database type, supports sqlite and mysql

[sqlite]
filename = d-mail.db  ## Database path

[manage]
port     = 12566      ## Management port
username = 111111     ## Management username
password = 111111     ## Management password
webPath  = web        ## Management page path

[api]
port = 12566          ## API port

[mysql]
host     = 127.0.0.1   ## Database address
port     = 3306         ## Database port
dbname   = d_mail      ## Database name
charset  = utf8        ## Database charset
username = root        ## Database username
password = 123456      ## Database password
```


## API Documentation

### Send Email Interface

`POST /sendMail`

**Parameters**:

| Parameter Name | Type       | Required | Description           |
|----------------|------------|----------|-----------------------|
| token          | string     | Yes      | Authorization token   |
| subject        | string     | No       | Email subject         |
| content        | string     | Yes      | Email content         |
| recipients     | []string   | No       | Additional recipients |
| files          | []File     | No       | Attachment file list  |

**Success Response**:

```json
ok
```


### Complete Request Examples

**JSON Format (with attachments)**

```bash
curl -X POST 'http://127.0.0.1:12567/sendMail' \
--header 'Content-Type: application/json' \
--data-raw '{
  "token": "{{token}}",
  "subject": "Urgent System Notification",
  "content": "test",
  "recipients": ["ops@example.com"],
  "files": [
    {
      "name": "alert.log",
      "data": "{{base64_content}}"
    }
  ]
}'
```


**Form Submission (with multiple attachments)**

```bash
curl -X POST 'http://127.0.0.1:12567/sendMail' \
--form 'token={{token}}' \
--form 'subject=Urgent System Notification' \
--form 'content=test' \
--form 'recipients=finance@example.com,sales@example.com' \
--form 'files=@"/data/reports/sales.pdf"' \
--form 'files=@"/data/reports/expenses.xlsx"'
```


## Build Instructions

1. Clone the repository:
2. Compile the frontend interface (requires building [d-mail-view](https://github.com/chuccp/d-mail-view) first):
3. Compile the binary file: