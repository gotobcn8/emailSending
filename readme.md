```cpp
## go原生的smtp包发送邮件
go发送邮件的封装包有jordon,gomail等等，可以直接搜得到，这边展示一下如何用go原生的包进行邮件发送。

邮件发送目前绝大部分公司都已经不使用明文发送了，所以发送内容都是经过加密，加密的协议方式有SSL和TLS，不懂的可以搜一下，都是类似地加密协议，SSL和TLS有很多共同点，可以大致理解为TLS在SSL的基础上加强了。

发送邮件的域名服务器一般为smtp开头的，例如qq邮箱的发送端服务器为smtp.qq.com

邮箱发送的端口，一般465/587是SSL加密端口，其中587可以用来作为TLS加密端口。

放代码，先定义以下邮箱基本信息
```go
type EmailType string
type EmailConf struct {
    Address    string   `json:"address"`     //域名
    Port       uint16   `json:"port"`        //端口
    SecureMode string   `json:"secure_mode"` //安全模式
    Username   string   `json:"username"`    //用户名,发送邮箱地址
    PassWord   string   `json:"password"`    //客户端专用密码，一般为授权码
    Receivers  []string `json:"receiver"`    //接收邮箱地址
}

func NewEmailer(receivers []string) *EmailConf {
    return &EmailConf{
        Address:    "smtp.qq.com",
        Port:       587,
        SecureMode: "TLS",
        Username:   "yyyyy@qq.com",
        PassWord:   "shouquanma",
        Receivers:  receivers,
    }
}
```

TLS发送方式和SSL发送方式，因为go原生支持TLS，所以直接调用原生函数就好了。SSL官方未提供。但是建立连接的过程如下。

```go
func (e *EmailConf) SendEmail(auth smtp.Auth, msg string) (err error) {
    addr := fmt.Sprintf("%s:%d", e.Address, e.Port)
    if e.SecureMode == "TLS" {
        return smtp.SendMail(addr, auth, e.Username, e.Receivers, []byte(msg)) //这个是go原生支持的email发送方式
    }
    return e.SendWithSSL(addr, auth, msg)
}

func (e *EmailConf) SendWithSSL(addr string, auth smtp.Auth, msg string) error {
    // Merge the To, Cc, and Bcc fields
    c, err := e.dial(addr)
    if err != nil {
        log.Printf("Create smpt client error:%v", err)
        return err
    }
    defer c.Close()
    if auth != nil {
        if ok, _ := c.Extension("AUTH"); ok {
            if err = c.Auth(auth); err != nil {
                log.Printf("Error during AUTH: %v", err)
                return err
            }
        }
    }
    if err = c.Mail(e.Username); err != nil {
        return err
    }
    for _, addr := range e.Receivers {
        if err = c.Rcpt(addr); err != nil {
            return err
        }
    }
    w, err := c.Data()
    if err != nil {
        return err
    }
    if _, err = w.Write([]byte(msg)); err != nil {
        return err
    }
    if err = w.Close(); err != nil {
        return err
    }
    return c.Quit()
}

func (e *EmailConf) dial(addr string) (*smtp.Client, error) {
    conn, err := tls.Dial("tcp", addr, nil)
    if err != nil {
        log.Printf("Dialing Error: %v", err)
        return nil, err
    }
    return smtp.NewClient(conn, e.Address)
}

```

经过测试，国内许多邮箱都是不支持TLS加密方式的，这点离国外有很大差距，包括很知名的网易163，新浪邮箱等等都不支持TLS, 目前测试已知的有QQ邮箱支持TLS, 想要测试更多支持TLS邮箱的可以尝试使用gmail和icloud邮箱（貌似icloud邮箱不支持SSL了）。

```

```
