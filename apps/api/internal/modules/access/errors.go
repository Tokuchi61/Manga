package access

import "errors"

var (
	ErrValidation              = errors.New("access_validation_failed")
	ErrRoleAlreadyExists       = errors.New("access_role_already_exists")
	ErrRoleNotFound            = errors.New("access_role_not_found")
	ErrPermissionAlreadyExists = errors.New("access_permission_already_exists")
	ErrPermissionNotFound      = errors.New("access_permission_not_found")
	ErrPolicyConflict          = errors.New("access_policy_conflict")
	ErrPolicyNotFound          = errors.New("access_policy_not_found")
	ErrAuthorizationDenied     = errors.New("access_authorization_denied")
)
