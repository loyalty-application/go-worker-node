package testing

import (
	"github.com/loyalty-application/go-worker-node/collections"
	"github.com/loyalty-application/go-worker-node/models"
)

func AddCampaignsTest() {

	acceptedMccs := make([]int, 0)
	acceptedMccs = append(acceptedMccs, 6625)
	acceptedMccs = append(acceptedMccs, 6626)
	acceptedMccs = append(acceptedMccs, 3434)
	campaign1 := models.Campaign{
		CampaignId:   "1",
		CardType:     "scis_freedom",
		StartDate:    "1/1/2021",
		EndDate:      "15/4/2023",
		Merchant:     "Marquardt  Kassulke and Keeling",
		MinSpend:     50.0,
		BonusRates:   0.05,
		AcceptedMCCs: acceptedMccs,
	}

	acceptedMccs = make([]int, 0)
	acceptedMccs = append(acceptedMccs, 3000)
	acceptedMccs = append(acceptedMccs, 3001)
	campaign2 := models.Campaign{
		CampaignId:   "2",
		CardType:     "scis_platinummiles",
		StartDate:    "31/1/2021",
		EndDate:      "31/3/2023",
		MinSpend:     500.0,
		BonusRates:   0.1,
		AcceptedMCCs: acceptedMccs,
	}

	acceptedMccs = make([]int, 0)
	acceptedMccs = append(acceptedMccs, 5000)
	acceptedMccs = append(acceptedMccs, 5001)
	acceptedMccs = append(acceptedMccs, 5002)
	campaign3 := models.Campaign{
		CampaignId:   "3",
		CardType:     "scis_shopping",
		StartDate:    "31/1/2021",
		EndDate:      "31/5/2023",
		Merchant:     "Reilly",
		MinSpend:     1000.0,
		BonusRates:   0.07,
		AcceptedMCCs: acceptedMccs,
	}

	acceptedMccs = make([]int, 0)
	acceptedMccs = append(acceptedMccs, 5000)
	acceptedMccs = append(acceptedMccs, 5001)
	acceptedMccs = append(acceptedMccs, 5002)
	campaign4 := models.Campaign{
		CampaignId:   "4",
		CardType:     "scis_premiummiles",
		StartDate:    "31/1/2023",
		EndDate:      "31/5/2023",
		MinSpend:     200.0,
		BonusRates:   0.17,
		AcceptedMCCs: acceptedMccs,
	}

	acceptedMccs = make([]int, 0)
	acceptedMccs = append(acceptedMccs, 3000)
	acceptedMccs = append(acceptedMccs, 3001)
	campaign5 := models.Campaign{
		CampaignId:   "5",
		CardType:     "scis_platinummiles",
		StartDate:    "31/1/2022",
		EndDate:      "31/3/2022",
		MinSpend:     1000.0,
		BonusRates:   0.15,
		AcceptedMCCs: acceptedMccs,
	}

	acceptedMccs = make([]int, 0)
	acceptedMccs = append(acceptedMccs, 3000)
	acceptedMccs = append(acceptedMccs, 3001)
	campaign6 := models.Campaign{
		CampaignId:   "6",
		CardType:     "scis_freedom",
		StartDate:    "31/1/2021",
		EndDate:      "31/7/2023",
		Merchant:     "West",
		MinSpend:     5.50,
		BonusRates:   0.15,
		AcceptedMCCs: acceptedMccs,
	}

	acceptedMccs = make([]int, 0)
	acceptedMccs = append(acceptedMccs, 3434)
	acceptedMccs = append(acceptedMccs, 3001)
	campaign7 := models.Campaign{
		CampaignId:   "7",
		CardType:     "scis_premiummiles",
		StartDate:    "31/1/2021",
		EndDate:      "21/6/2023",
		MinSpend:     15.0,
		BonusRates:   0.15,
		AcceptedMCCs: acceptedMccs,
	}

	// Add campaigns into DB
	campaignList := make([]models.Campaign, 0)
	campaignList = append(campaignList, campaign1)
	campaignList = append(campaignList, campaign2)
	campaignList = append(campaignList, campaign3)
	campaignList = append(campaignList, campaign4)
	campaignList = append(campaignList, campaign5)
	campaignList = append(campaignList, campaign6)
	campaignList = append(campaignList, campaign7)

	temp := models.CampaignList{
		Campaigns: campaignList,
	}
	collections.CreateCampaign("1", temp)
}
