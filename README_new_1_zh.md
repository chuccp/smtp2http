# SMTP2HTTP - SMTP转HTTP 工具

**English**🌎 | [**简体中文**🀄](./README_zh.md)

## 重要更新

增加定时任务功能，可定时读取其它api或链接发送邮件 。并支持模板配置，将json数据转换成文本，作为邮件内容发送 。

## 项目描述

将SMTP协议转换为HTTP接口的网关服务，帮助开发者：

- 无需在代码中硬编码SMTP配置
- 通过REST API动态发送邮件
- 可视化配置多个SMTP服务商
- 使用定时任务，减少项目中邮件配置

## 主要特性

- 🚀 通过Web UI配置SMTP服务器和接收邮箱
- 📦 支持GET/POST/JSON多种请求格式
- 🔒 基于Token的API访问控制
- 📎 多文件附件支持（Base64编码/表单上传）
- 🐳 开箱即用的Docker镜像
- 📊 发送记录查询与统计
- 📅 支持定时任务，请求链接，发送邮件
- 📧 支持邮件模板

## 社区
欢迎加入微信群或者telegram，提供更多意见。

微信群：

<img src="https://github.com/chuccp/smtp2http/blob/main/image/WeChat.png?raw=true" alt="WebChat" width="200">

telegram：

https://t.me/+JClG9_DojaM0ZGE1

## 快速开始

### 直接运行

```bash
# Windows
curl -LO https://github.com/chuccp/smtp2http/releases/latest/download/smtp2http-windows-amd64.tar.gz
./smtp2http.exe

# Linux
curl -LO https://github.com/chuccp/smtp2http/releases/latest/download/smtp2http-linux-amd64.tar.gz
chmod +x smtp2http
./smtp2http
```

### Docker运行

```bash
docker run -d \
  -p 12566:12566 \
  -p 12567:12567 \
  cooge123/smtp2http:latest
```

## 配置说明

首次启动后生成配置文件 `config.ini`：

```ini
[core]
init      = true   ##初始化开关，初始化配置完成后变为true 
cachePath = .cache  ##邮件附件缓存目录
dbType    = sqlite  ##数据库类型，支持sqlite和mysql

[sqlite]
filename = d-mail.db  ##数据库路径

[manage]
port     = 12566      ##管理端口   
username = 111111     ##管理用户名    
password = 111111     ##管理密码
webPath  = web        ##管理页面路径

[api]
port = 12566          ##API端口    

[mysql]
host     = 127.0.0.1   ##数据库地址
port     = 3306         ##数据库端口
dbname   = d_mail      ##数据库名称
charset  = utf8        ##数据库字符集
username = root        ##数据库用户名
password = 123456      ##数据库密码
```

## API文档

### 发送邮件接口

`POST /sendMail`

**参数**：

| 参数名        | 类型       | 必填 | 说明      |
|------------|----------|----|---------|
| token      | string   | 是  | 授权令牌    |
| subject    | string   | 否  | 邮件主题    |
| content    | string   | 是  | 邮件内容    |
| recipients | []string | 否  | 额外收件人列表 |
| files      | []File   | 否  | 附件文件列表  |

**成功响应**：

```json
ok
```

### 完整请求示例

**JSON格式（含附件）**

```bash
curl -X POST 'http://127.0.0.1:12567/sendMail' \
--header 'Content-Type: application/json' \
--data-raw '{
  "token": "{{token}}",
  "subject": "紧急系统通知",
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

**表单提交（含多个附件）**

```bash
curl -X POST 'http://127.0.0.1:12567/sendMail' \
--form 'token={{token}}' \
--form 'subject=紧急系统通知' \
--form 'content=test' \
--form 'recipients=finance@example.com,sales@example.com' \
--form 'files=@"/data/reports/sales.pdf"' \
--form 'files=@"/data/reports/expenses.xlsx"'
```

## 构建说明

1. 克隆仓库：
2. 编译前端界面（需先构建[d-mail-view](https://github.com/chuccp/d-mail-view)）：
3. 编译二进制文件：







        