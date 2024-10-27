package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	relationshipCtrl "github.com/koeylp/friends-management/cmd/internal/controller/relationship"
	"github.com/koeylp/friends-management/cmd/internal/handler/rest/response"
	"github.com/koeylp/friends-management/cmd/internal/model/dto/relationship/block"
	"github.com/koeylp/friends-management/cmd/internal/model/dto/relationship/friend"
	"github.com/koeylp/friends-management/cmd/internal/model/dto/relationship/subscription"
	utils "github.com/koeylp/friends-management/cmd/internal/pkg/error_util"
)

// RelationshipHandler handles HTTP requests for relationship-related operations.
type RelationshipHandler struct {
	relationshipCtrl relationshipCtrl.RelationshipController
}

// NewRelationshipHandler initializes a new RelationshipHandler with the provided service.
func NewRelationshipHandler(relationshipCtrl relationshipCtrl.RelationshipController) *RelationshipHandler {
	return &RelationshipHandler{relationshipCtrl: relationshipCtrl}
}

// CreateFriendHandler handles the creation of a friendship relationship.
func (h *RelationshipHandler) CreateFriendHandler(w http.ResponseWriter, r *http.Request) {
	var createFriendReq friend.CreateFriend
	err := json.NewDecoder(r.Body).Decode(&createFriendReq)
	if err != nil || len(createFriendReq.Friends) != 2 || createFriendReq.Friends[0] == createFriendReq.Friends[1] {
		response.NewBadRequestError("Invalid request payload").Send(w)
		return
	}
	if err := friend.ValidateCreateFriendRequest(&createFriendReq); err != nil {
		response.NewBadRequestError(err.Error()).Send(w)
		return
	}

	err = h.relationshipCtrl.CreateFriend(context.Background(), &createFriendReq)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	createdResponse := response.NewCREATED(nil)
	createdResponse.Send(w)
}

// GetFriendListByEmailHandler handles retrieving a friend list by user email.
func (h *RelationshipHandler) GetFriendListByEmailHandler(w http.ResponseWriter, r *http.Request) {
	var emailReq friend.EmailRequest
	err := json.NewDecoder(r.Body).Decode(&emailReq)
	if err != nil {
		response.NewBadRequestError("Invalid request payload: unable to decode JSON").Send(w)
		return
	}

	if err := friend.ValidateEmailRequest(&emailReq); err != nil {
		response.NewBadRequestError(err.Error()).Send(w)
		return
	}

	friends, err := h.relationshipCtrl.GetFriendListByEmail(context.Background(), emailReq.Email)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	okResponse := response.NewOK(map[string]interface{}{"friends": friends})
	okResponse.Send(w)
}

// GetCommonListHandler handles retrieving a common friend list for two users.
func (h *RelationshipHandler) GetCommonListHandler(w http.ResponseWriter, r *http.Request) {
	var commonFriendsReq friend.CommonFriendListReq
	err := json.NewDecoder(r.Body).Decode(&commonFriendsReq)
	if err != nil || len(commonFriendsReq.Friends) != 2 || commonFriendsReq.Friends[0] == commonFriendsReq.Friends[1] {
		response.NewBadRequestError("Invalid request payload").Send(w)
		return
	}
	if err := friend.ValidateCommonFriendListRequest(&commonFriendsReq); err != nil {
		response.NewBadRequestError(err.Error()).Send(w)
		return
	}

	commonList, err := h.relationshipCtrl.GetCommonList(context.Background(), &commonFriendsReq)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	okResponse := response.NewOK(map[string]interface{}{"friends": commonList})
	okResponse.Send(w)
}

// SubscribeHandler handles the subscription of updates between users.
func (h *RelationshipHandler) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	var subcribeReq subscription.SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&subcribeReq); err != nil || subcribeReq.Requestor == subcribeReq.Target {
		response.NewBadRequestError("Invalid request payload").Send(w)
		return
	}
	if err := subscription.ValidateSubscribeRequest(&subcribeReq); err != nil {
		response.NewBadRequestError(err.Error()).Send(w)
		return
	}

	err := h.relationshipCtrl.Subscribe(context.Background(), &subcribeReq)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	createdResponse := response.NewCREATED(nil)
	createdResponse.Send(w)
}

// BlockUpdatesHandler handles blocking updates between users.
func (h *RelationshipHandler) BlockUpdatesHandler(w http.ResponseWriter, r *http.Request) {
	var blockReq block.BlockRequest
	if err := json.NewDecoder(r.Body).Decode(&blockReq); err != nil || blockReq.Requestor == blockReq.Target {
		response.NewBadRequestError("Invalid request payload").Send(w)
		return
	}

	if err := block.ValidateBlockRequest(&blockReq); err != nil {
		response.NewBadRequestError(err.Error()).Send(w)
		return
	}

	err := h.relationshipCtrl.BlockUpdates(context.Background(), &blockReq)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	createdResponse := response.NewCREATED(nil)
	createdResponse.Send(w)
}

// GetUpdatableEmailAddressesHandler retrieves emails that can receive updates.
func (h *RelationshipHandler) GetUpdatableEmailAddressesHandler(w http.ResponseWriter, r *http.Request) {
	var recipientsReq subscription.RecipientRequest
	if err := json.NewDecoder(r.Body).Decode(&recipientsReq); err != nil {
		response.NewBadRequestError("Invalid request payload").Send(w)
		return
	}

	if err := subscription.ValidateRecipientRequest(&recipientsReq); err != nil {
		response.NewBadRequestError(err.Error()).Send(w)
		return
	}

	recipients, err := h.relationshipCtrl.GetUpdatableEmailAddresses(context.Background(), &recipientsReq)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	okResponse := response.NewOK(map[string]interface{}{"recipients": recipients})
	okResponse.Send(w)
}
