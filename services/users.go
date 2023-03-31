package services

import (
	"log"
	"time"

	"github.com/loyalty-application/go-worker-node/models"
)

func GetUserFromRecord(userRecord models.UserRecord) (result models.User, err error) {

	// create new user
	result = models.User{
		UserID: &userRecord.Id,
		FirstName: &userRecord.FirstName,
		LastName: &userRecord.LastName,
		Email: &userRecord.Email,
		Phone: &userRecord.Phone,
	}

	// Convert CreatedAt and UpdatedAt to time.Time
	createdAt, err := time.Parse("2006-01-02 15:04:05.999999999", userRecord.CreatedAt)
	if err != nil {
		log.Println("CreatedAt Parsing Error", err.Error())
	}
	updatedAt, err := time.Parse("2006-01-02 15:04:05.999999999", userRecord.UpdatedAt)
	if err != nil {
		log.Println("UpdatedAt Parsing Error", err.Error())
	}

	result.CreatedAt = createdAt
	result.UpdatedAt = updatedAt

	// // Store into database
	// _, err = collections.CreateUser(result)
	// if err != nil {
	// 	log.Println(err.Error())
	// }

	return result, err
}