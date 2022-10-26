package main

import (
    "crypto/tls"
    "fmt"
    "log"
    "net/smtp"
)

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

func main() {
    e := NewEmailer([]string{"xxxx@receiver.com"})
    auth := smtp.PlainAuth("", e.Username, e.PassWord, e.Address)
    if err := e.SendEmail(auth, "测试邮件"); err != nil {
        fmt.Println("sending email failed!")
    }
}

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
