package models

type User struct {
	Id              interface{} `bson:"id,omitempty" json:"id,omitempty"`
	TelegramUserId  int64       `bson:"telegramUserId" json:"telegramUserId"`
	Username        string      `bson:"username" json:"username"`
	IsPremium       bool        `bson:"isPremium" json:"isPremium"`
	ReferralCode    string      `bson:"referralCode" json:"referralCode"`
	ReferredBy      interface{} `bson:"referredBy" json:"referredBy"`
	Wallet          interface{} `bson:"wallet" json:"wallet"`
	AllowsWriteToPm bool        `bson:"allowsWriteToPm" json:"allowsWriteToPm"`

	Coins int `bson:"coins" json:"coins"`
}
