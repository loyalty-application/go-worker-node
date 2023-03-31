package services

import (
	"errors"
	"log"

	"github.com/loyalty-application/go-worker-node/models"
)



func GetCardFromRecord(userRecord models.UserRecord) (result models.Card, err error) {
	// cardId := userRecord.CardId
	
	// // Check if card with card id already exists
	// if card, err := collections.RetrieveSpecificCard(cardId); err != mongo.ErrNoDocuments {
	// 	log.Println("Card already exists")
	// 	return card, err
	// }

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
		CardId: userRecord.CardId,
		UserId: userRecord.Id,
		CardPan: userRecord.CardPan,
		CardType: userRecord.CardType,
		ValueType: ProcessCardType(userRecord.CardType),
		Value: 0,
	}

	// // Store into database
	// _, err = collections.CreateCard(result)
	// if err != nil {
	// 	log.Println(err.Error())
	// }

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

func Contains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}