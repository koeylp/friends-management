package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/koeylp/friends-management/internal/dto/relationship/friend"
	"github.com/koeylp/friends-management/internal/responses"
	"github.com/koeylp/friends-management/internal/services"
)

type RelationshipHandler struct {
	relationshipService services.RelationshipService
}

func NewRelationshipHandler(relationshipService services.RelationshipService) *RelationshipHandler {
	return &RelationshipHandler{relationshipService: relationshipService}
}

func (h *RelationshipHandler) CreateFriendHandler(w http.ResponseWriter, r *http.Request) {
	var createFriendReq friend.CreateFriend
	err := json.NewDecoder(r.Body).Decode(&createFriendReq)
	if err != nil || len(createFriendReq.Friends) != 2 {
		responses.NewBadRequestError("Invalid request payload").Send(w)
		return
	}

	err = h.relationshipService.CreateFriend(context.Background(), &createFriendReq)
	if err != nil {
		responses.NewBadRequestError(err.Error()).Send(w)
		return
	}

	createdResponse := responses.NewCREATED("Create Friend Connection Successfully!", nil)
	createdResponse.Send(w)
}
