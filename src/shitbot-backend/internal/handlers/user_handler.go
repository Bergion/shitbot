package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"shitbot/internal/auth"
	"shitbot/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const InitialCoinNumber = 100

type Response struct {
	Success bool        `json:"success"`
	Message interface{} `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type UserHandler struct {
	collection *mongo.Collection
}

func NewUserHandler(client *mongo.Client) *UserHandler {
	collection := client.Database("botDB").Collection("users")
	return &UserHandler{collection: collection}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response := &Response{}
	w.Header().Set("Content-Type", "application/json")
	tgUser, ok := auth.FromContext(ctx)
	if !ok {
		response.Message = "Forbidden"
		responseBody, _ := json.Marshal(response)
		w.Write(responseBody)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	refCode := r.URL.Query().Get("ref_code")
	var referredBy models.User
	if refCode != "" {
		err := h.collection.FindOne(ctx, bson.M{"referralCode": refCode}).Decode(&referredBy)
		if err != nil {
			log.Println(err)
			response.Message = "Invalid referral code"
			responseBody, _ := json.Marshal(response)
			w.Write(responseBody)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	user := &models.User{
		Id:              primitive.NewObjectID(),
		TelegramUserId:  tgUser.ID,
		Username:        tgUser.Username,
		IsPremium:       tgUser.IsPremium,
		AllowsWriteToPm: tgUser.AllowsWriteToPm,
		ReferralCode:    primitive.NewObjectID().Hex(),
		ReferredBy:      referredBy.Id,
		Wallet:          nil,
		Coins:           InitialCoinNumber,
	}

	if _, err := h.collection.InsertOne(ctx, user); err != nil {
		log.Println(err)
		response.Message = "Internal server"
		responseBody, _ := json.Marshal(response)
		w.Write(responseBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Success = true
	response.Data = user
	responseBody, _ := json.Marshal(response)
	w.Write(responseBody)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response := &Response{}
	w.Header().Set("Content-Type", "application/json")
	tgUser, ok := auth.FromContext(ctx)
	if !ok {
		response.Message = "Forbidden"
		responseBody, _ := json.Marshal(response)
		w.Write(responseBody)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var user models.User
	err := h.collection.FindOne(ctx, bson.M{"telegramUserId": tgUser.ID}).Decode(&user)
	if err != nil {
		log.Println(err)
		response.Message = "User not found"
		responseBody, _ := json.Marshal(response)
		w.Write(responseBody)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response.Success = true
	response.Data = user
	responseBody, _ := json.Marshal(response)
	w.Write(responseBody)
}
