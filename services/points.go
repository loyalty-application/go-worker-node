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

	var EXCLUDED_MCCS = []int{6051, 9399, 6540}

	// If MCC is to be excluded, do not convert to points
	mcc, _ := strconv.Atoi(transaction.MCC)
	if In(EXCLUDED_MCCS, mcc) {
		return
	}

	// Get amount spent, if in foreign currency, convert to SGD
	amountSpent := transaction.Amount * GetExchangeRate(transaction.Currency)

	// Conversion: $$ -> POINTS
	if transaction.CardType == "scis_shopping" {
		calculatePoints(transaction, amountSpent)
		
	// Conversion: $$ -> CASHBACK
	} else if transaction.CardType == "scis_freedom" {
		calculateCashBack(transaction, amountSpent)

	// Conversion: $$ -> MILES
	} else {
		calculateMiles(transaction, amountSpent)
	}

	// INSERT: Apply campaign benefits
}


func calculatePoints(transaction *models.Transaction, amountSpent float64) {

	// MCC code categorization
	var ONLINE_SHOPPING = []int{5999, 5964, 5691, 5311, 5411, 5399, 5815, 5816, 5817, 5818}
	var SHOPPING = []int{
		4816, 5045, 5262, 5309, 5310, 5311, 5331, 5399, 5611, 5621, 5631, 5641, 5651, 5655, 5661, 
		5691, 5699, 5732, 5733, 5734, 5735, 5912, 5942, 5944, 5945, 5946, 5947, 5948, 5949, 5964, 
		5965, 5966, 5967, 5968, 5969, 5970, 5992, 5999, 5621, 5631}

	mcc, _ := strconv.Atoi(transaction.MCC)

	pointsConversionRate := 0.01  // Base conversion rate for all spend types

	// Bonus for online shopping
	if In(ONLINE_SHOPPING, mcc) {
		pointsConversionRate = 0.1

	// Bonus for shopping
	} else if In(SHOPPING, mcc) {
		pointsConversionRate = 0.04
	}
	
	transaction.Points = amountSpent * float64(pointsConversionRate)
}

func calculateCashBack(transaction *models.Transaction, amountSpent float64) {
	cashBack := 0.05 * amountSpent
	if amountSpent > 2000 {
		cashBack = 0.03 * amountSpent
	} else if amountSpent > 500 {
		cashBack = 0.01 * amountSpent
	}
	transaction.CashBack = cashBack
}

func calculateMiles(transaction *models.Transaction, amountSpent float64) {

	var HOTEL = []int{7011}
	mcc, _ := strconv.Atoi(transaction.MCC)

	// Get base rate for premium vs platinum
	milesConversionRate := 0.011
	if transaction.CardType == "scis_platinummiles" {
		milesConversionRate = 0.014
	}

	// If spend type is foreign
	if transaction.Currency != "SGD" {
		// Base conversion rate for foreign card spend
		milesConversionRate = 0.022
		if transaction.CardType == "scis_platinummiles" {
			milesConversionRate = 0.03

			// If it is a foreign hotel spend
			if In(HOTEL, mcc) { milesConversionRate = 0.06 }
		}
		
	// For local hotel spends
	} else if In(HOTEL, mcc) {
		milesConversionRate = 0.03
	}
	transaction.Miles = amountSpent * milesConversionRate
}