package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/koeylp/friends-management/internal/dto/relationship/block"
	"github.com/koeylp/friends-management/internal/dto/relationship/friend"
	"github.com/koeylp/friends-management/internal/dto/relationship/subscription"
	"github.com/koeylp/friends-management/internal/responses"
	"github.com/koeylp/friends-management/internal/services"
	"github.com/koeylp/friends-management/utils"
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
	if err != nil || len(createFriendReq.Friends) != 2 || createFriendReq.Friends[0] == createFriendReq.Friends[1] {
		responses.NewBadRequestError("Invalid request payload").Send(w)
		return
	}

	err = h.relationshipService.CreateFriend(context.Background(), &createFriendReq)
	if err != nil {
		utils.HandleError(w, err)
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

	okResponse := responses.NewOK(map[string]interface{}{"friends": friends})
	okResponse.Send(w)
}

func (h *RelationshipHandler) GetCommonListHandler(w http.ResponseWriter, r *http.Request) {
	var commonFriendsReq friend.CommonFriendListReq
	err := json.NewDecoder(r.Body).Decode(&commonFriendsReq)
	if err != nil || len(commonFriendsReq.Friends) != 2 || commonFriendsReq.Friends[0] == commonFriendsReq.Friends[1] {
		responses.NewBadRequestError("Invalid request payload").Send(w)
		return
	}
	if err := friend.ValidateCommonFriendListRequest(&commonFriendsReq); err != nil {
		responses.NewBadRequestError(err.Error()).Send(w)
		return
	}

	commonList, err := h.relationshipService.GetCommonList(context.Background(), &commonFriendsReq)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	okResponse := responses.NewOK(map[string]interface{}{"friends": commonList})
	okResponse.Send(w)
}

func (h *RelationshipHandler) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	var subcribeReq subscription.SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&subcribeReq); err != nil || subcribeReq.Requestor == subcribeReq.Target {
		responses.NewBadRequestError("Invalid request payload").Send(w)
		return
	}
	if err := subscription.ValidateSubscribeRequest(&subcribeReq); err != nil {
		responses.NewBadRequestError(err.Error()).Send(w)
		return
	}

	err := h.relationshipService.Subscribe(context.Background(), &subcribeReq)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	createdResponse := responses.NewCREATED(nil)
	createdResponse.Send(w)
}

func (h *RelationshipHandler) BlockUpdatesHandler(w http.ResponseWriter, r *http.Request) {
	var blockReq block.BlockRequest
	if err := json.NewDecoder(r.Body).Decode(&blockReq); err != nil || blockReq.Requestor == blockReq.Target {
		responses.NewBadRequestError("Invalid request payload").Send(w)
		return
	}

	if err := block.ValidateBlockRequest(&blockReq); err != nil {
		responses.NewBadRequestError(err.Error()).Send(w)
		return
	}

	err := h.relationshipService.BlockUpdates(context.Background(), &blockReq)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	createdResponse := responses.NewCREATED(nil)
	createdResponse.Send(w)
}

func (h *RelationshipHandler) GetUpdatableEmailAddressesHandler(w http.ResponseWriter, r *http.Request) {
	var recipientsReq subscription.RecipientRequest
	if err := json.NewDecoder(r.Body).Decode(&recipientsReq); err != nil {
		responses.NewBadRequestError("Invalid request payload").Send(w)
		return
	}

	if err := subscription.ValidateRecipientRequest(&recipientsReq); err != nil {
		responses.NewBadRequestError(err.Error()).Send(w)
		return
	}

	recipients, err := h.relationshipService.GetUpdatableEmailAddresses(context.Background(), &recipientsReq)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	okResponse := responses.NewOK(map[string]interface{}{"recipients": recipients})
	okResponse.Send(w)
}
