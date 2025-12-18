package tool

import (
	"context"
	"errors"
	"fmt"
	"go-MQ/common"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 模拟你的 EnvVars 配置
type EmailConfig struct {
	Host   string
	Port   int
	User   string
	Pass   string
	From   string
	Secure bool // 如果是 465 端口通常需要 SSL
}

const (
	EmailCodeExpiration = 10 * time.Minute
	EmailSendLimit      = 60 * time.Second // 限制 60 秒
)

var config = EmailConfig{
	Host:   viper.GetString("email.host"),
	Port:   viper.GetInt("email.port"),
	User:   viper.GetString("email.user"),
	Pass:   viper.GetString("email.pass"),
	From:   viper.GetString("email.from"),
	Secure: viper.GetBool("email.secure"),
}

// SendVerificationEmail 发送验证码邮件
func SendVerificationEmail(to string) (bool, string, error) {
	e := email.NewEmail()

	code := generateVerificationCode()

	// 设置发件人、收件人、主题
	e.From = config.From
	e.To = []string{to}
	e.Subject = "MQ - 邮箱验证"

	// 设置 HTML 内容 (复制自你的 Node.js 代码)
	e.HTML = []byte(fmt.Sprintf(`
		<div style="padding: 20px; font-family: Arial, sans-serif;">
			<p>您的验证码是：</p>
			<div style="background: #f5f5f5; padding: 15px; border-radius: 5px; font-size: 24px; font-weight: bold; text-align: center; letter-spacing: 5px;">
				%s
			</div>
			<p style="color: #666; margin-top: 20px;">验证码有效期为 10 分钟，请尽快使用。</p>
			<p style="color: #999; font-size: 12px; margin-top: 30px;">如果这不是您的操作，请忽略此邮件。</p>
		</div>
	`, code))

	// SMTP 服务器地址
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	auth := smtp.PlainAuth("", config.User, config.Pass, config.Host)

	// 发送邮件
	err := e.Send(addr, auth)
	if err != nil {
		logrus.Error("邮件发送失败: ", err)
		return false, "", err
	}

	logrus.Info("邮件发送成功")

	return true, code, nil
}

// GenerateVerificationCode 生成 6 位随机验证码
func generateVerificationCode() string {
	// 初始化随机种子
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := r.Intn(900000) + 100000
	return fmt.Sprintf("%d", code)
}

func SaveEmailCode(ctx context.Context, useremail string, code string) error {
	rdb := common.GetRedisClient()
	userKey := useremail

	// 1. 检查是否存在发送限制锁
	limitKey := "emaillimit:" + userKey
	exists, err := rdb.Exists(ctx, limitKey).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		// 如果键存在，说明 60 秒还没过
		return errors.New("发送过于频繁，请 60 秒后再试")
	}

	// 2. 使用管道（Pipeline）或事务确保两个键同时写入
	// 一个存验证码，一个存限制锁
	pipe := rdb.Pipeline()

	// 写入验证码（10分钟有效）
	pipe.Set(ctx, "emailcode:"+userKey, code, EmailCodeExpiration)

	// 写入限制锁（60秒有效），值可以随便填
	pipe.Set(ctx, limitKey, "1", EmailSendLimit)

	_, err = pipe.Exec(ctx)
	if err != nil {
		logrus.Error("redis set email code and limit failed: ", err)
		return err
	}

	return nil
}

func IsTrueEmailCode(ctx context.Context, useremail string, code string) bool {
	storedCode, err := common.GetRedisClient().Get(ctx, "emailcode:"+useremail).Result()
	if err != nil {
		logrus.Error("redis get email code failed: ", err)
		return false
	}
	return storedCode == code
}
