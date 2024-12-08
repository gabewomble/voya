package data

import (
	"server/internal/repository"
	"server/internal/validator"
)

type ValidateUpdateMemberStatusParams struct {
	Fieldname     string
	MemberStatus  repository.MemberStatusEnum
	TargetMember  repository.GetTripMemberRow
	CurrentMember repository.GetTripMemberRow
	Validator     *validator.Validator
}

func ValidateUpdateMemberStatus(params ValidateUpdateMemberStatusParams) {
	if params.CurrentMember.MemberStatus == repository.MemberStatusEnumOwner && params.TargetMember.ID == params.CurrentMember.ID {
		params.Validator.AddError(params.Fieldname, "owner can't modify their own status")
	}

	if params.MemberStatus == params.TargetMember.MemberStatus {
		params.Validator.AddError(params.Fieldname, "member status is already set to this value")
	}

	if params.TargetMember.MemberStatus == repository.MemberStatusEnumOwner && params.CurrentMember.MemberStatus != repository.MemberStatusEnumOwner {
		params.Validator.AddError(params.Fieldname, "user can only modify owner if they are an owner")
	}

	switch params.MemberStatus {
	case repository.MemberStatusEnumAccepted:
		validateUpdateStatusAcceptOrDecline(params)
	case repository.MemberStatusEnumDeclined:
		validateUpdateStatusAcceptOrDecline(params)
	case repository.MemberStatusEnumCancelled:
		validateUpdateStatusCancelled(params)
	case repository.MemberStatusEnumPending:
		validateUpdateStatusPending(params)
	case repository.MemberStatusEnumOwner:
		validateUpdateStatusOwner(params)
	case repository.MemberStatusEnumRemoved:
		validateUpdateStatusRemoved(params)
	}
}

func validateUpdateStatusOwner(params ValidateUpdateMemberStatusParams) {
	if params.TargetMember.MemberStatus != repository.MemberStatusEnumAccepted {
		params.Validator.AddError(params.Fieldname, "user must be member of trip to be owner")
	}
	if params.CurrentMember.MemberStatus != repository.MemberStatusEnumOwner {
		params.Validator.AddError(params.Fieldname, "user can only set owner if they are an owner")
	}
}

func validateUpdateStatusPending(params ValidateUpdateMemberStatusParams) {
	if params.TargetMember.ID == params.CurrentMember.ID {
		params.Validator.AddError(params.Fieldname, "user can't invite themselves")
	}
	switch params.TargetMember.MemberStatus {
	case repository.MemberStatusEnumAccepted:
	case repository.MemberStatusEnumOwner:
		params.Validator.AddError(params.Fieldname, "user is already a member of the trip")
	}
}

func validateUpdateStatusCancelled(params ValidateUpdateMemberStatusParams) {
	if params.TargetMember.MemberStatus != repository.MemberStatusEnumPending {
		params.Validator.AddError(params.Fieldname, "user can only cancel pending invites")
	}
	if params.CurrentMember.ID == params.TargetMember.ID {
		params.Validator.AddError(params.Fieldname, "user can't cancel their own invite")
	}
}

func validateUpdateStatusAcceptOrDecline(params ValidateUpdateMemberStatusParams) {
	if params.TargetMember.MemberStatus != repository.MemberStatusEnumPending {
		params.Validator.AddError(params.Fieldname, "user can only accept or decline pending invites")
	}
	if params.CurrentMember.ID != params.TargetMember.ID {
		params.Validator.AddError(params.Fieldname, "user can only accept or decline their own invite")
	}
}

func validateUpdateStatusRemoved(params ValidateUpdateMemberStatusParams) {
	if (params.CurrentMember.MemberStatus != repository.MemberStatusEnumAccepted) && (params.CurrentMember.MemberStatus != repository.MemberStatusEnumOwner) {
		params.Validator.AddError(params.Fieldname, "user can only remove members if they are an owner or accepted member")
	}
}
