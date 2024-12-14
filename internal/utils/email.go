package utils

import (
	"bytes"
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/FIY-pc/BBingyan/internal/infrastructure/logger"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

func GenerateHTMLMsg(email, user, nickname, subject, body string) []byte {
	msg := []byte("To: " + email + "\r\n" +
		"From: " + user + "\r\n" + "<" + nickname + ">\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
		body)
	return msg
}

func GetTemplatePath(filename string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	// 使用 filepath 包来处理路径
	path := strings.Split(dir, "BBingyan")[0]
	templatePath := filepath.Join(path, "BBingyan"+"/web/templates/"+filename)
	return templatePath, nil
}

// GenerateEmailBody 从模板文件生成邮件内容
func GenerateEmailBody(filename string, data interface{}) (string, error) {
	// 打开模板文件
	filePath, err := GetTemplatePath(filename)
	if err != nil {
		logger.Log.Error(nil, err.Error())
		return "", err
	}
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		logger.Log.Error(nil, err.Error())
		return "", err
	}

	// 创建一个字符串缓冲区来存储生成的内容
	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return "", err
	}
	return body.String(), nil
}

// SendEmail 发送邮件
func SendEmail(email, subject, body string) error {
	var addr = config.Configs.Smtp.SmtpHost + ":" + config.Configs.Smtp.SmtpPort
	var auth = smtp.PlainAuth("", config.Configs.Smtp.SmtpUser, config.Configs.Smtp.SmtpPassword, config.Configs.Smtp.SmtpHost)
	var from = config.Configs.Smtp.SmtpUser
	var to = []string{email}
	var msg = GenerateHTMLMsg(email, config.Configs.Smtp.SmtpUser, config.Configs.Smtp.SmtpNickname, subject, body)
	return smtp.SendMail(addr, auth, from, to, msg)
}
