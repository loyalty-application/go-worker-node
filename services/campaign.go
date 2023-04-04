package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/loyalty-application/go-worker-node/models"
)

// TODO:
// Return the campaign that is applicable to the user, else return nil
func ApplyApplicableCampaign(transaction *models.Transaction, allCampaigns []models.Campaign) (result models.Campaign, hasCampaign bool) {

	amountSpent := transaction.Amount * getExchangeRate(transaction.Currency)

	var bonus float64 = 0

	hasCampaign = false
	// Iterate through campaigns and find the best applicable one
	for _, campaign := range allCampaigns {

		// If campaign is applicable, apply bonus, return
		if isApplicable(campaign, transaction) {
			campaignBonus := campaign.BonusRates * amountSpent
			if campaignBonus > bonus {
				bonus = campaignBonus

				// PROBLEMATIC: Using pointer leads to incorrect attachment
				transaction.CampaignApplied = &campaign
				fmt.Println("Assigned to campaign:", campaign)

				hasCampaign = true
				result = campaign
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

	return result, hasCampaign
}


// Takes in a campaign and a transaction, returns true if the campaign is
// applicable to the transaction, else false
func isApplicable(campaign models.Campaign, transaction *models.Transaction) bool {
	// DEBUG
	fmt.Println("Campaign CT:", campaign.CardType)
	fmt.Println("Trans CT   :", transaction.CardType)

	// Check for matching card type
	if campaign.CardType != transaction.CardType {
		return false
	}

	// Check for min. spending
	if transaction.Amount < campaign.MinSpend {
		return false
	}

	// Check if within campaign date
	if !inTimeSpan(campaign.StartDate, campaign.EndDate, transaction.DateTime) {
		return false
	}

	// Check for applicable spend type (mcc)
	mcc, _ := strconv.Atoi(transaction.MCC)
	if !in(campaign.AcceptedMCCs, mcc) {
		return false
	}

	fmt.Println("Applied", campaign.CampaignId, " to", transaction.Id)
	return true
}


func inTimeSpan(start, end, check time.Time) bool {
    return check.After(start) && check.Before(end)
}