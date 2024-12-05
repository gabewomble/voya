package data

import (
	"reflect"
	"server/internal/validator"
)

type memberStatus struct {
	Pending   string
	Accepted  string
	Declined  string
	Removed   string
	Cancelled string
}

var MemberStatus = memberStatus{
	Pending:   "pending",
	Accepted:  "accepted",
	Declined:  "declined",
	Removed:   "removed",
	Cancelled: "cancelled",
}

func IsValidMemberStatus(status string) bool {
	v := reflect.ValueOf(MemberStatus)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).String() == status {
			return true
		}
	}
	return false
}

func ValidateMemberStatus(v *validator.Validator, status string, fieldName string) {
	v.CheckStrNotEmpty(status, fieldName)
	v.Check(IsValidMemberStatus(status), fieldName, "must be one of \"pending\", \"accepted\", \"declined\", \"removed\", or \"cancelled\" ")
}
