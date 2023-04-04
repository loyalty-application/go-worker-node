package models

type CampaignList struct {
	Campaigns []Campaign `json:"campaigns" bson:",inline"`
}

type Campaign struct {
	// TODO: should have on merchantId
	UserId       string  `json:"user_id" bson:"user_id" example:"u00001"`
	CampaignId   string  `json:"campaign_id" bson:"campaign_id" example:"cmp00001"`
	Merchant     string  `json:"merchant" bson:"merchant" example:"7-11"`
	CardType     string  `json:"card_type" bson:"card_type" example:"super_miles_card"`
	Description  string  `json:"description" bson:"description"`
	StartDate    string  `json:"start_date" bson:"start_date" example:"2023-03-02T13:10:23Z"`
	EndDate      string  `json:"end_date" bson:"end_date" example:"2023-03-03T13:10:23Z"`
	MinSpend     float64 `json:"min_spend" bson:"min_spend" example:"100.0"`
	BonusRates   float64 `json:"bonus_rates" bson:"bonus_rates" example:"8.8"`
	AcceptedMCCs []int   `json:"accepted_mccs" bson:"accepted_mccs"`
}
