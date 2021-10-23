package payment

type PaymentType struct {
	Price           float64 `bson:"price"            json:"price"`
	Monthly         bool    `bson:"monthly"          json:"monthly"`
	Yearly          bool    `bson:"yearly"           json:"yearly"`
	Once            bool    `bson:"once"             json:"once"`
	Free            bool    `bson:"free"             json:"free"`
	SubscriptionUID string  `bson:"subscription_uid" json:"subscription_uid"`
}

func (p *PaymentType) ExportFields() (bool, bool, bool, bool, float64) {
	return p.Monthly, p.Yearly, p.Once, p.Free, p.Price
}
