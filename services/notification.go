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
		m.SetHeader("To", "ojh809@gmail.com")
		m.SetHeader("Subject", "Campaign Notification")
		m.SetBody("text/plain", messageList[i])

		if err := gomail.Send(s, m); err != nil {
			log.Printf("Failed to send email to %s: %v", email, err)
		} else {
			log.Printf("Email sent to %s", email)
		}
	}
}