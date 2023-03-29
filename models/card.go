package models

type CardList struct {
	Cards []Card `json:"cards" bson:",inline"`
}

type Card struct {
	CardId    string  `json:"card_id" bson:"_id" example:"4111222233334444"`
	UserEmail string  `json:"user_email" bson:"user_email" example:"sudo@gabriel.dev"`
	CardPan   string  `json:"card_pan" bson:"card_pan" example:"xyz"`
	CardType  string  `json:"card_type" bson:"card_type" example:"super_miles_card"`
	ValueType string  `json:"value_type"  bson:"value_type" example:"miles"`
	Value     float64 `json:"value" bson:"value" example:"100.0"`
}

type CardUpdateRequest struct {
	Value float64 `json:"value" bson:"value" example:"100.0"`
}
