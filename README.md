在项目中，经常由用到邮件来通知异常日志，但这需要在项目中配置STMP，以及需要接收邮件的邮箱，但是当邮件地址发生变更的时候，就需要修改项目的配置文件。或者由于网络限制导致STMP无法配置。

本项目则是使用http接口的替代STMP，只需要在本项目中配置好STMP以及需要接收邮件的电子邮箱地址，就给可以实现邮件的发送

支持get以及post请求

get 请求 例子：

```powershell
curl  'http://127.0.0.1:12566/sendMail?token=99eaf30feb23e28057367431d820cf319915792921d9cf21b5f761fb75433225&content=this%20is%20a%20test'
```

post请求例子

```powershell
curl  'http://127.0.0.1:12567/sendMail?token=d6a1ee40c5bad981461643f5404a305a2e3f480cc6fcf65ba98efb63ce32d471&content=1212&subject=1212' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'token=99eaf30feb23e28057367431d820cf319915792921d9cf21b5f761fb75433225' \
--data-urlencode 'content=this is a test'
```

发送带文件的邮件

```powershell
curl  'http://127.0.0.1:12567/sendMail?token=d6a1ee40c5bad981461643f5404a305a2e3f480cc6fcf65ba98efb63ce32d471&content=1212&subject=1212' \
--form 'files=@"/C:/Users/cooge/Documents/111111.txt"' \
--form 'files=@"/C:/Users/cooge/Documents/22222222222222.txt"' \
--form 'token="d6a1ee40c5bad981461643f5404a305a2e3f480cc6fcf65ba98efb63ce32d471"' \
--form 'content="1212"'
```

其实就是简单的表单提交，这样，不同的语言，不同的平台都可以很好的使用本项目

参数说明

token：在管理界面通过手动添加得到，token是唯一值，与STMP，要接收的电子邮箱绑定

content：邮件内容

subject：邮件主题，在生成token的时候，也设置了一个subject，当参数给空的时候，就会使用设置token时的subject

files：邮件文件，支持多个文件