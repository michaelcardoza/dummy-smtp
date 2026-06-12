package main

import (
	"fmt"

	gomail "gopkg.in/mail.v2"
)

func main() {
	message := gomail.NewMessage()

	message.SetHeader("From", "alice@example.com")
	message.SetHeader("To", "bob@example.com")
	message.SetHeader("Subject", "Mail Example From GoLang")

	text := "Hi,\nHow are you?\nThis email was sent via GoLang's smtplib module."
	message.SetBody("text/plain", text)

	html := `
	<html>
		<body>
			<p>Hi,<br>
			How are you?<br>
			This email was sent via GoLang's smtplib module.
			</p>
		</body>
	</html>
	`
	message.AddAlternative("text/html", html)

	dialer := gomail.NewDialer("localhost", 1025, "", "")

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Sent - open http://localhost:8025")
	}
}
