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

	createdResponse := responses.NewCREATED(nil)
	createdResponse.Send(w)
}

func (h *RelationshipHandler) GetFriendListByEmailHandler(w http.ResponseWriter, r *http.Request) {
	var emailReq friend.EmailRequest
	err := json.NewDecoder(r.Body).Decode(&emailReq)
	if err != nil {
		responses.NewBadRequestError("Invalid request payload: unable to decode JSON").Send(w)
		return
	}

	if emailReq.Email == "" {
		responses.NewBadRequestError("Email cannot be empty").Send(w)
		return
	}

	friends, err := h.relationshipService.GetFriendListByEmail(context.Background(), emailReq.Email)
	if err != nil {
		responses.NewInternalServerError(err.Error()).Send(w)
		return
	}

	friendList := friend.FriendList{
		Friends: friends,
		Count:   len(friends),
	}

	okResponse := responses.NewOK(friendList)
	okResponse.Send(w)
}

func (h *RelationshipHandler) GetCommonListHandler(w http.ResponseWriter, r *http.Request) {
	var commonFriendsReq friend.CommonFriendListReq
	err := json.NewDecoder(r.Body).Decode(&commonFriendsReq)
	if err != nil || len(commonFriendsReq.Friends) != 2 {
		responses.NewBadRequestError("Invalid request payload").Send(w)
		return
	}

	commonList, err := h.relationshipService.GetCommonList(context.Background(), &commonFriendsReq)
	if err != nil {
		responses.NewInternalServerError(err.Error()).Send(w)
		return
	}

	friendList := friend.FriendList{
		Friends: commonList,
		Count:   len(commonList),
	}

	okResponse := responses.NewOK(friendList)
	okResponse.Send(w)
}
