package services

import (
	"log"
	"time"

	"github.com/loyalty-application/go-worker-node/collections"
	"github.com/loyalty-application/go-worker-node/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserFromRecord(userRecord models.UserRecord) (result models.User, err error) {
	userEmail := userRecord.Email 

	// Retrieve user from db, if no user with email, create new user
	result, err = collections.RetrieveSpecificUser(userEmail)
	if err == mongo.ErrNoDocuments {
		result = models.User{
			UserID: &userRecord.Id,
			FirstName: &userRecord.FirstName,
			LastName: &userRecord.LastName,
			Email: &userRecord.Email,
			Phone: &userRecord.Phone,
		}
	}

	// Add new card id if not already inside
	if !contains(result.Card, userRecord.CardId) {
		result.Card = append(result.Card, userRecord.CardId)
	}

	// Convert CreatedAt and UpdatedAt to time.Time
	createdAt, err := time.Parse("2006-01-02 15:04:05.000000", userRecord.CreatedAt)
	if err != nil {
		log.Println("CreatedAt Parsing Error", err.Error())
	}
	updatedAt, err := time.Parse("2006-01-02 15:04:05.000000", userRecord.UpdatedAt)
	if err != nil {
		log.Println("UpdatedAt Parsing Error", err.Error())
	}

	result.CreatedAt = createdAt
	result.UpdatedAt = updatedAt

	// Store into database
	_, err = collections.CreateUser(result)
	if err != nil {
		log.Println(err.Error())
	}

	return result, err
}

func contains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}