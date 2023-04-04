package testing

import (
	"time"

	"github.com/loyalty-application/go-worker-node/models"
	"github.com/loyalty-application/go-worker-node/collections"
)

func AddCampaignsTest() {
	
	acceptedMccs := make([]int, 0)
	acceptedMccs = append(acceptedMccs, 6625)
	acceptedMccs = append(acceptedMccs, 6626)
	acceptedMccs = append(acceptedMccs, 3434)
	startDate, _ := time.Parse("2/1/2006", "1/1/2021")
	endDate, _ := time.Parse("2/1/2006", "15/4/2023")
	campaign1 := models.Campaign {
		CampaignId: "1",
		CardType: "scis_freedom",
		StartDate: startDate,
		EndDate: endDate,
		MinSpend: 50.0,
		BonusRates: 0.05,
		AcceptedMCCs: acceptedMccs,
	}

	acceptedMccs = make([]int, 0)
	acceptedMccs = append(acceptedMccs, 3000)
	acceptedMccs = append(acceptedMccs, 3001)
	startDate, _ = time.Parse("2/1/2006", "31/1/2021")
	endDate, _ = time.Parse("2/1/2006", "31/3/2023")
	campaign2 := models.Campaign {
		CampaignId: "2",
		CardType: "scis_platinummiles",
		StartDate: startDate,
		EndDate: endDate,
		MinSpend: 500.0,
		BonusRates: 0.1,
		AcceptedMCCs: acceptedMccs,
	}

	acceptedMccs = make([]int, 0)
	acceptedMccs = append(acceptedMccs, 5000)
	acceptedMccs = append(acceptedMccs, 5001)
	acceptedMccs = append(acceptedMccs, 5002)
	startDate, _ = time.Parse("2/1/2006", "31/1/2021")
	endDate, _ = time.Parse("2/1/2006", "31/5/2023")
	campaign3 := models.Campaign {
		CampaignId: "3",
		CardType: "scis_shopping",
		StartDate: startDate,
		EndDate: endDate,
		MinSpend: 10.0,
		BonusRates: 0.07,
		AcceptedMCCs: acceptedMccs,
	}

	acceptedMccs = make([]int, 0)
	acceptedMccs = append(acceptedMccs, 5000)
	acceptedMccs = append(acceptedMccs, 5001)
	acceptedMccs = append(acceptedMccs, 5002)
	startDate, _ = time.Parse("2/1/2006", "31/1/2023")
	endDate, _ = time.Parse("2/1/2006", "31/5/2023")
	campaign4 := models.Campaign {
		CampaignId: "4",
		CardType: "scis_premiummiles",
		StartDate: startDate,
		EndDate: endDate,
		MinSpend: 200.0,
		BonusRates: 0.17,
		AcceptedMCCs: acceptedMccs,
	}

	acceptedMccs = make([]int, 0)
	acceptedMccs = append(acceptedMccs, 3000)
	acceptedMccs = append(acceptedMccs, 3001)
	startDate, _ = time.Parse("2/1/2006", "31/1/2022")
	endDate, _ = time.Parse("2/1/2006", "31/3/2022")
	campaign5 := models.Campaign {
		CampaignId: "5",
		CardType: "scis_platinummiles",
		StartDate: startDate,
		EndDate: endDate,
		MinSpend: 1000.0,
		BonusRates: 0.15,
		AcceptedMCCs: acceptedMccs,
	}

	acceptedMccs = make([]int, 0)
	acceptedMccs = append(acceptedMccs, 3000)
	acceptedMccs = append(acceptedMccs, 3001)
	startDate, _ = time.Parse("2/1/2006", "31/1/2021")
	endDate, _ = time.Parse("2/1/2006", "31/7/2023")
	campaign6 := models.Campaign {
		CampaignId: "6",
		CardType: "scis_premiummiles",
		StartDate: startDate,
		EndDate: endDate,
		MinSpend: 19.0,
		BonusRates: 0.15,
		AcceptedMCCs: acceptedMccs,
	}

	acceptedMccs = make([]int, 0)
	acceptedMccs = append(acceptedMccs, 3434)
	acceptedMccs = append(acceptedMccs, 3001)
	startDate, _ = time.Parse("2/1/2006", "31/1/2021")
	endDate, _ = time.Parse("2/1/2006", "31/6/2023")
	campaign7 := models.Campaign {
		CampaignId: "7",
		CardType: "scis_premiummiles",
		StartDate: startDate,
		EndDate: endDate,
		MinSpend: 15.0,
		BonusRates: 0.15,
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

	temp := models.CampaignList {
		Campaigns: campaignList,
	}
	collections.CreateCampaign("1", temp)
}
