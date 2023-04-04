package models

type CardList struct {
	Cards []Card `json:"cards" bson:",inline"`
}

type Card struct {
	CardId       string  `json:"card_id" bson:"card_id" example:"4111222233334444"`
	UserId       string  `json:"user_id" bson:"user_id" example:"12345"`
	CardPan      string  `json:"card_pan" bson:"card_pan" example:"xyz"`
	CardType     string  `json:"card_type" bson:"card_type" example:"super_miles_card"`
	ValueType    string  `json:"value_type"  bson:"value_type" example:"miles"`
	ShortCardPan string  `json:"short_card_pan" bson:"short_card_pan" example:"1234"`
	Value        float64 `json:"value" bson:"value" example:"100.0"`
}

type CardUpdateRequest struct {
	Value float64 `json:"value" bson:"value" example:"100.0"`
}
