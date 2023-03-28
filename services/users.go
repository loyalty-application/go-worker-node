package services

import (
	"errors"
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

func GetCardFromRecord(userRecord models.UserRecord) (result models.Card, err error) {
	cardId := userRecord.CardId
	
	// Check if card with card id already exists
	if card, err := collections.RetrieveSpecificCard(cardId); err != mongo.ErrNoDocuments {
		log.Println("Card already exists")
		return card, err
	}

	// Validate card type
	isValidCardType := false
	validCardType := [4]string{"scis_platinummiles", "scis_premiummiles", "scis_shopping", "scis_freedom"}
	for _, elem := range validCardType {
		if elem == userRecord.CardType {
			isValidCardType = true
		}
	}

	if (!isValidCardType) {
		log.Println("Invalid Card Type")
		err = errors.New("Invalid Card Type")
		return result, err
	}

	// Create card if card doesn't exist
	result = models.Card{
		UserEmail: userRecord.Email,
		CardId: userRecord.CardId,
		CardPan: userRecord.CardPan,
		CardType: userRecord.CardType,
		ValueType: ProcessCardType(userRecord.CardType),
		Value: 0,
	}

	// Store into database
	_, err = collections.CreateCard(result)
	if err != nil {
		log.Println(err.Error())
	}

	return result, err
}

// return type of card
func ProcessCardType(str string) string {

	switch str {
	case "scis_freedom":
		return "Points"
	case "scis_premiummiles", "scis_platinummiles":
		return "Miles"
	case "scis_shopping":
		return "Cashback"
	default:
		return "Error"
	}

}

func contains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}