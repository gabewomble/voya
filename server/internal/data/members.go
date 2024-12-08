package data

import (
	"server/internal/repository"
	"server/internal/validator"

	"github.com/google/uuid"
)

type ValidateUpdateMemberStatusParams struct {
	Fieldname    string
	MemberStatus repository.MemberStatusEnum
	OwnerID      uuid.UUID
	TargetUserID uuid.UUID
	UserID       uuid.UUID
	Validator    *validator.Validator
}

func ValidateUpdateMemberStatus(params ValidateUpdateMemberStatusParams) {
	if params.UserID != params.TargetUserID {
		switch params.MemberStatus {
		case repository.MemberStatusEnumAccepted:
		case repository.MemberStatusEnumDeclined:
			params.Validator.AddError(params.Fieldname, "user can only accept or decline their own invite")
		}
	} else {
		switch params.MemberStatus {
		case repository.MemberStatusEnumCancelled:
			params.Validator.AddError(params.Fieldname, "user cannot cancel their own invite")
		}
	}

	if params.MemberStatus == repository.MemberStatusEnumOwner && params.UserID != params.OwnerID {
		params.Validator.AddError(params.Fieldname, "user cannot set owner")
	}
}
