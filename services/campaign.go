package services

import (
	"strconv"

	"github.com/loyalty-application/go-worker-node/models"
	"github.com/loyalty-application/go-worker-node/collections"
)


// Return the campaign that is applicable to the user, else return nil
func ApplyApplicableCampaign(transaction *models.Transaction)  {

	amountSpent := transaction.Amount * getExchangeRate(transaction.Currency)

	// Retrieve all ACTIVE campaigns
	allCampaigns, _ := collections.RetrieveActiveCampaigns(transaction.TransactionDate)

	var bonus float64 = 0

	// Iterate through campaigns and find the best applicable one
	for _, campaign := range allCampaigns {

		// If campaign is applicable, apply bonus, return
		if isApplicable(campaign, transaction) {
			campaignBonus := campaign.BonusRates * amountSpent
			if campaignBonus > bonus {
				bonus = campaignBonus
			}
		}
	}

	// Apply the bonus to the transaction
	if transaction.CardType == "scis_shopping" {
		transaction.Points += bonus
	} else if transaction.CardType == "scis_freedom" {
		transaction.CashBack += bonus
	} else {
		transaction.Miles += bonus
	}
}


// Takes in a campaign and a transaction, returns true if the campaign is
// applicable to the transaction, else false
func isApplicable(campaign models.Campaign, transaction *models.Transaction) bool {
	// Check for matching card type
	if campaign.CardType != transaction.CardType {
		return false
	}

	// Check for min. spending
	if transaction.Amount < campaign.MinSpend {
		return false
	}

	// Check for applicable spend type (mcc)
	mcc, _ := strconv.Atoi(transaction.MCC)
	if !in(campaign.AcceptedMCCs, mcc) {
		return false
	}

	return true
}