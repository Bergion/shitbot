package handlers

import (
	"fmt"
	"log"
	"net/http"
	"shitbot/internal/auth"
	"shitbot/internal/models"
	response "shitbot/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	collection *mongo.Collection
}

func NewUserHandler(client *mongo.Client) *UserHandler {
	collection := client.Database("botDB").Collection("users")
	return &UserHandler{collection: collection}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tgUser, ok := auth.FromContext(ctx)
	if !ok {
		response.Forbidden(w, "Forbidden")
		return
	}

	user := models.NewUserFromTelegram(*tgUser)
	var writeModels []mongo.WriteModel

	refCode := r.URL.Query().Get("ref_code")
	if refCode != "" {
		var referredBy models.User
		err := h.collection.FindOne(ctx, bson.M{"referralCode": refCode}).Decode(&referredBy)
		if err != nil {
			log.Println(err)
		} else {
			referredBy.AddReferral(user)

			updateOperation := mongo.NewUpdateOneModel()
			updateOperation.SetFilter(bson.M{"id": referredBy.Id})
			updateOperation.SetUpdate(bson.M{"referrals": referredBy.Referrals})

			writeModels = append(writeModels, updateOperation)
		}
	}

	insertOperation := mongo.NewInsertOneModel()
	insertOperation.SetDocument(*user)
	writeModels = append(writeModels, insertOperation)
	fmt.Println(writeModels)
	if _, err := h.collection.BulkWrite(ctx, writeModels); err != nil {
		log.Println(err)
		response.InternalServerError(w, "Internal server error")
		return
	}

	response.Ok(w, user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tgUser, ok := auth.FromContext(ctx)
	if !ok {
		response.Forbidden(w, "Forbidden")
		return
	}

	var user models.User
	err := h.collection.FindOne(ctx, bson.M{"account.telegramUserId": tgUser.ID}).Decode(&user)
	if err != nil {
		log.Println(err)
		response.NotFound(w, "User not found")
		return
	}

	response.Ok(w, user)
}
