package services

import (
	"strconv"
	"github.com/loyalty-application/go-worker-node/models"
)

func getExchangeRate(currency string) float64 {
	// TODO: Implement an exchange rate getter
	if currency == "SGD" {
		return 1.0
	}
	if currency == "USD" {
		return 1.33
	}
	return 1.0
}


func in(arr []int, target int) bool {
	for _, item := range arr {
		if item == target { return true }
	}
	return false
}


func ConvertPoints(transaction *models.Transaction) {

	var EXCLUDED_MCCS = []int{6051, 9399, 6540}

	// If MCC is to be excluded, do not convert to points
	mcc, _ := strconv.Atoi(transaction.MCC)
	if in(EXCLUDED_MCCS, mcc) {
		return
	}

	// Get amount spent, if in foreign currency, convert to SGD
	amountSpent := transaction.Amount * getExchangeRate(transaction.Currency)

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
}


func calculatePoints(transaction *models.Transaction, amountSpent float64) {

	// MCC code categorization
	var ONLINE_SHOPPING = []int{5999, 5964, 5691, 5311, 5411, 5399, 5815, 5816, 5817, 5818}
	// var SHOPPING = []int{
	// 	4816, 5045, 5262, 5309, 5310, 5311, 5331, 5399, 5611, 5621, 5631, 5641, 5651, 5655, 5661, 
	// 	5691, 5699, 5732, 5733, 5734, 5735, 5912, 5942, 5944, 5945, 5946, 5947, 5948, 5949, 5964, 
	// 	5965, 5966, 5967, 5968, 5969, 5970, 5992, 5999, 5621, 5631}

	mcc, _ := strconv.Atoi(transaction.MCC)

	pointsConversionRate := 1  // Base conversion rate for all spend types

	// Bonus for online shopping
	if in(ONLINE_SHOPPING, mcc) {
		pointsConversionRate = 10

	// Bonus for shopping
	} else if mcc >= 5000 && mcc < 7299 {
		pointsConversionRate = 4
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

	var FOREIGN_HOTEL = []int{7011}
	mcc, _ := strconv.Atoi(transaction.MCC)

	// Get base rate for premium vs platinum
	milesConversionRate := 1.1
	if transaction.CardType == "scis_platinummiles" {
		milesConversionRate = 1.4
	}

	// If spend type is foreign
	if transaction.Currency != "SGD" {
		// Base conversion rate for foreign card spend
		milesConversionRate = 2.2
		if transaction.CardType == "scis_platinummiles" {
			milesConversionRate = 3

			// If it is a foreign hotel spend
			if in(FOREIGN_HOTEL, mcc) { milesConversionRate = 6 }
		}
		
	// For local hotel spends
	} else if mcc >= 3501 && mcc < 3831 {
		milesConversionRate = 3
	}
	transaction.Miles = amountSpent * milesConversionRate
}