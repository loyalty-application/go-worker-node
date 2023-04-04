package services

import (
	"crypto/tls"
	"log"
	"os"
	"strconv"

	"github.com/loyalty-application/go-worker-node/collections"
	"github.com/loyalty-application/go-worker-node/models"
	"gopkg.in/gomail.v2"
)

func SendNotification(notificationList []models.Notification) {
	
	// Load env variables
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	senderEmail := os.Getenv("SENDER_EMAIL")
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Println("Port is not a number")
		return
	}

	log.Println("Server =", smtpServer, " Port =", smtpPort, " Username =", smtpUsername, " Password =", smtpPassword)

	// Get email from cardId
	emailList := make([]string, 0)
	messageList := make([]string, 0)
	for _, notification := range notificationList {
		email, err := collections.RetrieveEmailFromCard(notification.CardId)
		if err == nil {
			emailList = append(emailList, email)
			messageList = append(messageList, notification.Message)
		}
	}
	log.Println("Email List =", emailList)
	log.Println("Message List =", messageList)
	
	d := gomail.NewDialer(smtpServer, smtpPort, smtpUsername, smtpPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	s, err := d.Dial()
	if err != nil {
		log.Println("Dial Error:", err.Error())
		return
	}
	defer s.Close()

	// Iterate over the list of recipients and message bodies, and send each email separately
	for i, email := range emailList {
		m := gomail.NewMessage()
		m.SetHeader("From", senderEmail)
		m.SetHeader("To", email)
		m.SetHeader("Subject", "Campaign Notification")
		m.SetBody("text/plain", messageList[i])

		if err := gomail.Send(s, m); err != nil {
			log.Printf("Failed to send email to %s: %v", email, err)
		} else {
			log.Printf("Email sent to %s", email)
		}
	}
	// // Connect to smtp server
	// auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)

	// // Create a TLS config object with insecureSkipVerify set to true (if your SMTP server uses a self-signed certificate)
	// tlsConfig := &tls.Config{InsecureSkipVerify: true}

	// conn, err := smtp.Dial(smtpServer + ":" + string(smtpPort))
	// if err != nil {
	// 	log.Println("Error connecting to SMTP server:", err.Error())
	// 	return
	// }
	// defer conn.Close()

	// // Use STARTTLS command to encrypt the connection
	// if err := conn.StartTLS(tlsConfig); err != nil {
	// 	log.Println("StartTLS Failed:", err.Error())
	// 	return
	// }

	// // Authenticate and send emails
	// if err := conn.Auth(auth); err != nil {
	// 	log.Println("Authentication Failed:", err.Error())
	// 	return
	// }

	// // Iterate over the list of recipients and message bodies, and send each email separately
	// for i, recipient := range emailList {

	// 	// Set Sender
	// 	if err := conn.Mail(senderEmail); err != nil {
	// 		log.Println("Sender doesn't exist:", senderEmail)
	// 		return
	// 	}

	// 	// Set Recipient
	// 	if err := conn.Rcpt("loyaltyapplication@outlook.com"); err != nil {
	// 		log.Println("Recipient doesn't exist:", recipient)
	// 		continue
	// 	}
	// 	data, err := conn.Data()
	// 	if err != nil {
	// 		log.Println("Data Error:", err.Error())
	// 		continue
	// 	}
	// 	defer data.Close()

	// 	msg := "From: " + senderEmail + "\n" +
	// 		"To: " + "loyaltyapplication@outlook.com" + "\n" +
	// 		"Subject: Campaign Notification\n" +
	// 		"\n" + messageList[i] + "\n"

	// 	_, err = data.Write([]byte(msg))
	// 	if err != nil {
	// 		log.Println("Write Error:", err.Error())
	// 		return
	// 	}

	// 	log.Printf("Email sent to %s", recipient)
	// }
}