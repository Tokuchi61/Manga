package entity

import "time"

// AccountState defines lifecycle status of a user account.
type AccountState string

const (
	AccountStateActive      AccountState = "active"
	AccountStateDeactivated AccountState = "deactivated"
	AccountStateBanned      AccountState = "banned"
)

// ProfileVisibility defines whether profile data is public or private.
type ProfileVisibility string

const (
	ProfileVisibilityPublic  ProfileVisibility = "public"
	ProfileVisibilityPrivate ProfileVisibility = "private"
)

// HistoryVisibilityPreference defines global visibility upper bound for history sharing.
type HistoryVisibilityPreference string

const (
	HistoryVisibilityPublic  HistoryVisibilityPreference = "public"
	HistoryVisibilityPrivate HistoryVisibilityPreference = "private"
)

// VIPAction defines lifecycle operations for VIP state.
type VIPAction string

const (
	VIPActionActivate   VIPAction = "activate"
	VIPActionFreeze     VIPAction = "freeze"
	VIPActionResume     VIPAction = "resume"
	VIPActionDeactivate VIPAction = "deactivate"
)

// UserAccount is the aggregate owner for stage-5 user account/profile/membership data.
type UserAccount struct {
	ID                          string
	CredentialID                string
	Username                    string
	DisplayName                 string
	Bio                         string
	AvatarURL                   string
	BannerURL                   string
	ProfileVisibility           ProfileVisibility
	HistoryVisibilityPreference HistoryVisibilityPreference
	AccountState                AccountState
	VIPActive                   bool
	VIPFrozen                   bool
	VIPStartedAt                *time.Time
	VIPEndsAt                   *time.Time
	VIPFrozenAt                 *time.Time
	VIPFreezeReason             string
	CreatedAt                   time.Time
	UpdatedAt                   time.Time
}

func (u UserAccount) IsActive() bool {
	return u.AccountState == AccountStateActive
}
