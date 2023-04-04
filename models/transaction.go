package models

type TransactionList struct {
	Transactions []Transaction `json:"transactions" bson:",inline"`
}

type Transaction struct {
	Id              string    `json:"id" bson:"id" example:"1"`
	TransactionId   string    `json:"transaction_id" bson:"transaction_id" example:"txn00001"`
	Merchant        string    `json:"merchant" bson:"merchant" example:"7-11"`
	MCC             string    `json:"mcc" bson:"mcc" example:"5311"`
	Currency        string    `json:"currency" bson:"currency" example:"USD"`
	Amount          float64   `json:"amount" bson:"amount" example:"20.10"`
	TransactionDate string    `json:"transaction_date" bson:"transaction_date" example:"yyyy-mm-dd hh:mm:ss"`
	CardId          string    `json:"card_id" bson:"card_id" example:"4111222233334444"`
	CardPan         string    `json:"card_pan" bson:"card_pan" example:"xyz"`
	CardType        string    `json:"card_type" bson:"card_type" example:"super_miles_card"`
	Points          float64   `json:"points" bson:"points" example:"4.5"`
	Miles           float64   `json:"miles" bson:"miles" example:"10.1"`
	CashBack        float64   `json:"cash_back" bson:"cash_back" example:"30.45"`
	CampaignApplied Campaign `json:"campaign_applied" bson:"campaign_applied" example:""`
}
