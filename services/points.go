package services

import (
	"strconv"
	"github.com/loyalty-application/go-worker-node/models"
)


// QUESTIONS:
// 1. How to determine online shopping? MCC like doesn't tell us much
// 2. How to determine exchange rate?
// 3. Does the conversion stack? 
//    E.g., 0.5% cashback on all spend*, 1.0% cashback for all spend > 500 SGD, 3% cashback for all spend > 2000 SGD



func GetExchangeRate(currency string) float64 {
	// TODO: Implement an exchange rate getter
	return 1.0
}


func In(arr []int, target int) bool {
	for _, item := range arr {
		if item == target { return true }
	}
	return false
}


func ConvertPoints(transaction *models.Transaction) {

	// MCC code categorization (SUBJECT TO CHANGES)
	var ONLINE_SHOPPING = []int{5999, 5964, 5691, 5311, 5411, 5399}
	// var HOTEL = []int{7011}
	var EXCLUDED_MCCS = []int{6051, 9399, 6540}

	// If MCC is to be excluded, do not convert to points
	mcc, _ := strconv.Atoi(transaction.MCC)
	if In(EXCLUDED_MCCS, mcc) {
		return
	}

	// Get amount spent
	amountSpent := transaction.Amount * GetExchangeRate(transaction.Currency)

	// Conversion: $$ -> POINTS
	if transaction.CardType == "scis_shopping" {
		pointsConversionRate := 4
		if In(ONLINE_SHOPPING, mcc) {  // Bonus for online shopping
			pointsConversionRate = 10
		}
		transaction.Points = amountSpent * float64(pointsConversionRate + 1)

	// Conversion: $$ -> CASHBACK
	} else if transaction.CardType == "scis_freedom" {
		cashBack := 0.005 * amountSpent
		if amountSpent > 500 {
			cashBack += 0.01 * amountSpent
		}
		if amountSpent > 2000 {
			cashBack += 0.03 * amountSpent
		}
		transaction.CashBack = cashBack

	// Conversion: $$ -> MILES
	} else {
		transaction.Miles = 0
	}
}