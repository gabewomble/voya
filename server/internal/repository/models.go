// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"server/internal/dbtypes"
)

type MemberStatusEnum string

const (
	MemberStatusEnumOwner     MemberStatusEnum = "owner"
	MemberStatusEnumPending   MemberStatusEnum = "pending"
	MemberStatusEnumAccepted  MemberStatusEnum = "accepted"
	MemberStatusEnumDeclined  MemberStatusEnum = "declined"
	MemberStatusEnumRemoved   MemberStatusEnum = "removed"
	MemberStatusEnumCancelled MemberStatusEnum = "cancelled"
)

func (e *MemberStatusEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MemberStatusEnum(s)
	case string:
		*e = MemberStatusEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for MemberStatusEnum: %T", src)
	}
	return nil
}

type NullMemberStatusEnum struct {
	MemberStatusEnum MemberStatusEnum `json:"member_status_enum"`
	Valid            bool             `json:"valid"` // Valid is true if MemberStatusEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullMemberStatusEnum) Scan(value interface{}) error {
	if value == nil {
		ns.MemberStatusEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.MemberStatusEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullMemberStatusEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.MemberStatusEnum), nil
}

type NotificationType string

const (
	NotificationTypeTripCancelled         NotificationType = "trip_cancelled"
	NotificationTypeTripDateChange        NotificationType = "trip_date_change"
	NotificationTypeTripInvitePending     NotificationType = "trip_invite_pending"
	NotificationTypeTripInviteAccepted    NotificationType = "trip_invite_accepted"
	NotificationTypeTripInviteCancelled   NotificationType = "trip_invite_cancelled"
	NotificationTypeTripInviteDeclined    NotificationType = "trip_invite_declined"
	NotificationTypeTripMemberLeft        NotificationType = "trip_member_left"
	NotificationTypeTripMemberRemoved     NotificationType = "trip_member_removed"
	NotificationTypeTripOwnershipTransfer NotificationType = "trip_ownership_transfer"
)

func (e *NotificationType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = NotificationType(s)
	case string:
		*e = NotificationType(s)
	default:
		return fmt.Errorf("unsupported scan type for NotificationType: %T", src)
	}
	return nil
}

type NullNotificationType struct {
	NotificationType NotificationType `json:"notification_type"`
	Valid            bool             `json:"valid"` // Valid is true if NotificationType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullNotificationType) Scan(value interface{}) error {
	if value == nil {
		ns.NotificationType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.NotificationType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullNotificationType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.NotificationType), nil
}

type Notification struct {
	ID        uuid.UUID                    `json:"id"`
	UserID    uuid.UUID                    `json:"user_id"`
	TripID    uuid.UUID                    `json:"trip_id"`
	Type      NotificationType             `json:"type"`
	Message   string                       `json:"message"`
	CreatedAt *time.Time                   `json:"created_at"`
	ReadAt    *time.Time                   `json:"read_at"`
	Metadata  dbtypes.NotificationMetadata `json:"metadata"`
}

type Token struct {
	Hash         []byte     `json:"hash"`
	UserID       uuid.UUID  `json:"user_id"`
	Expiry       *time.Time `json:"expiry"`
	Scope        string     `json:"scope"`
	RefreshToken []byte     `json:"refresh_token"`
}

type Trip struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
	StartDate   pgtype.Date `json:"start_date"`
	EndDate     pgtype.Date `json:"end_date"`
	CreatedAt   *time.Time  `json:"created_at"`
	UpdatedAt   *time.Time  `json:"updated_at"`
}

type TripMember struct {
	TripID       uuid.UUID        `json:"trip_id"`
	UserID       uuid.UUID        `json:"user_id"`
	InvitedBy    uuid.UUID        `json:"invited_by"`
	MemberStatus MemberStatusEnum `json:"member_status"`
	UpdatedBy    uuid.UUID        `json:"updated_by"`
	UpdatedAt    *time.Time       `json:"updated_at"`
}

type User struct {
	ID           uuid.UUID  `json:"id"`
	CreatedAt    *time.Time `json:"created_at"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash []byte     `json:"password_hash"`
	Activated    bool       `json:"activated"`
	Version      int32      `json:"version"`
	Username     string     `json:"username"`
}
