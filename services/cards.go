package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"

	"github.com/loyalty-application/go-worker-node/collections"
	"github.com/loyalty-application/go-worker-node/models"
)

func UpdateCardValues(cardIdList []string) {

	// key = cardId, value = points / miles / cashback
	cardMap := make(map[string]float64)
	for _, cardId := range cardIdList {
		temp, _ := collections.RetrieveCardValuesFromTransaction(cardId)
		cardMap[cardId] = temp
	}

	collections.UpdateCardValues(cardMap)
}

func GetCardFromRecord(userRecord models.UserRecord) (result models.Card, err error) {

	// Validate card type
	isValidCardType := false
	validCardType := [4]string{"scis_platinummiles", "scis_premiummiles", "scis_shopping", "scis_freedom"}
	for _, elem := range validCardType {
		if elem == userRecord.CardType {
			isValidCardType = true
		}
	}

	if !isValidCardType {
		log.Println("Invalid Card Type")
		err = errors.New("Invalid Card Type")
		return result, err
	}

	// Create card if card doesn't exist
	result = models.Card{
		CardId:    userRecord.CardId,
		UserId:    userRecord.Id,
		CardPan:   userRecord.CardPan,
		ShortCardPan: userRecord.CardPan[len(userRecord.CardPan)-4:],
		CardType:  userRecord.CardType,
		ValueType: ProcessCardType(userRecord.CardType),
		Value:     0,
	}

	hash := sha256.Sum256([]byte(userRecord.CardPan))
	hashBytes := hash[:]
	hashString := hex.EncodeToString(hashBytes)

	result.CardPan = hashString

	// // Store into database
	// _, err = collections.CreateCard(result)
	// if err != nil {
	// 	log.Println(err.Error())
	// }

	return result, err
}

func HashCardPan(cardPan string) (hashedCardPan string) {
	hash := sha256.Sum256([]byte(cardPan))
	hashBytes := hash[:]
	hashedCardPan = hex.EncodeToString(hashBytes)

	return hashedCardPan
}

// return type of card
func ProcessCardType(str string) string {

	switch str {
	case "scis_freedom":
		return "Cashback"
	case "scis_premiummiles", "scis_platinummiles":
		return "Miles"
	case "scis_shopping":
		return "Points"
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