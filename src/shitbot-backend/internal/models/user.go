package models

import (
	"shitbot/internal/auth"
	"shitbot/internal/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           interface{} `bson:"id,omitempty" json:"id,omitempty"`
	Account      Account     `bson:"accountSettings" json:"accountSettings"`
	ReferralCode string      `bson:"referralCode" json:"referralCode"`
	ReferredBy   interface{} `bson:"referredBy" json:"referredBy"`
	Wallet       interface{} `bson:"wallet" json:"wallet"`
	Referrals    []Referral  `bson:"referrals" json:"referrals"`
	Coins        int64       `bson:"coins" json:"coins"`
	TapLevel     int         `bson:"tapLevel" json:"tapLevel"`
	Stat         Stat        `bson:"stat" json:"stat"`
}

type Account struct {
	TelegramUserId  int64  `bson:"telegramUserId" json:"telegramUserId"`
	Username        string `bson:"username" json:"username"`
	IsPremium       bool   `bson:"isPremium" json:"isPremium"`
	AllowsWriteToPm bool   `bson:"allowsWriteToPm" json:"allowsWriteToPm"`
}

type Referral struct {
	Username  string `bson:"username" json:"username"`
	IsPremium bool   `bson:"isPremium" json:"isPremium"`
	Coins     int    `bson:"coins" json:"coins"`
}

type Stat struct {
	Earned int `bson:"earned" json:"earned"`
	Taps   int `bson:"taps" json:"taps"`
}

func NewUserFromTelegram(telegramUser auth.TelegramUser) *User {
	account := Account{
		TelegramUserId:  telegramUser.ID,
		Username:        telegramUser.Username,
		IsPremium:       telegramUser.IsPremium,
		AllowsWriteToPm: telegramUser.AllowsWriteToPm,
	}
	return &User{
		Id:           primitive.NewObjectID(),
		Account:      account,
		ReferralCode: primitive.NewObjectID().Hex(),
		Wallet:       nil,
		Coins:        constants.InitialCoinNumber,
		Referrals:    make([]Referral, 0),
		TapLevel:     1,
		Stat:         Stat{},
	}
}

func (u *User) AddReferral(referredUser *User) {
	referral := Referral{
		Username:  referredUser.Account.Username,
		IsPremium: referredUser.Account.IsPremium,
	}
	referredUser.ReferredBy = u.Id
	u.Referrals = append(u.Referrals, referral)
}
