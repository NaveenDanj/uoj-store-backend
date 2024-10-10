package service

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/mail.v2"
)

func SendEmail(to string, subject string, html string) error {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}

	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	host := os.Getenv("EMAIL_HOST")

	m := mail.NewMessage()

	m.SetHeader("From", username)

	m.SetHeader("To", to)

	m.SetHeader("Subject", subject)

	m.SetBody("text/html", html)

	d := mail.NewDialer(host, 587, username, password)

	err = d.DialAndSend(m)
	if err != nil {
		return err
	}

	return nil
}

func ProcessOTPEmail(otp string, name string) string {
	return `<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OTP Verification - UOJ-Store</title>
    <style>
        body {
            font-family: 'Arial', sans-serif;
            background-color: #f4f4f4;
            margin: 0;
            padding: 0;
        }

        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            border-radius: 8px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            padding: 20px;
        }

        .header {
            background-color: #4a90e2;
            color: white;
            text-align: center;
            padding: 20px 0;
            border-radius: 8px 8px 0 0;
        }

        .header h1 {
            margin: 0;
            font-size: 24px;
        }

        .content {
            padding: 20px;
            text-align: center;
        }

        .content p {
            font-size: 18px;
            color: #555555;
            margin-bottom: 20px;
        }

        .otp-box {
            display: inline-block;
            background-color: #f1f1f1;
            padding: 15px;
            border-radius: 8px;
            font-size: 32px;
            font-weight: bold;
            letter-spacing: 10px;
            color: #333333;
        }

        .button {
            display: inline-block;
            background-color: #4a90e2;
            color: white;
            text-decoration: none;
            padding: 12px 20px;
            border-radius: 8px;
            font-size: 16px;
            margin-top: 20px;
        }

        .footer {
            margin-top: 30px;
            text-align: center;
            color: #999999;
            font-size: 14px;
        }

        .footer p {
            margin: 5px 0;
        }
    </style>
</head>

<body>
    <div class="container">
        <div class="header">
            <h1>OTP Verification - UOJ-Store</h1>
        </div>
        <div class="content">
            <p>Dear ` + name + `,</p>
            <p>To verify your account on <strong>UOJ-Store</strong>, please use the following OTP:</p>
            <div class="otp-box">` + otp + `</div>
            <p>This OTP is valid for 10 minutes. Please do not share it with anyone.</p>
            <a href="https://your-app-link.com" class="button">Access UOJ-Store</a>
        </div>
        <div class="footer">
            <p>Thank you for using UOJ-Store!</p>
            <p>If you didn’t request this OTP, please contact our support immediately.</p>
        </div>
    </div>
</body>

</html>
`
}
