/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package allerror provides a set of error codes and error types used in the application.
package allerror

import (
	"errors"
	"strings"
)

const (
	// ErrorCodeInvalidParam is const
	errorCodeNoPermission = "no_permission"

	// ErrorCodeUserNotFound is const
	ErrorCodeUserNotFound = "user_not_found"

	// ErrorCodeUserNotFound is const
	ErrorCodeFileUploadFailed = "file_upload_failed"

	// ErrorRateLimitOver is const
	ErrorRateLimitOver = "rate_limit_over"

	// ErrorCodeModelNotFound is const
	ErrorCodeModelNotFound = "model_not_found"

	// ErrorCodeDatasetNotFound is const
	ErrorCodeDatasetNotFound = "dataset_not_found"

	// ErrorCodeSpaceNotFound is const
	ErrorCodeSpaceNotFound = "space_not_found"

	// ErrorCodeResourceDisabled is const
	ErrorCodeResourceDisabled = "resource_disabled"

	// ErrorCodeResourceNoApplicationFile is const
	ErrorCodeResourceNoApplicationFile = "resource_no_application_file"

	// ErrorCodeResourceAlreadyDisabled is const
	ErrorCodeResourceAlreadyDisabled = "resource_already_disabled"

	// ErrorCodeSpaceVariableNotFound space variable
	ErrorCodeSpaceVariableNotFound = "space_variable_not_found"

	// ErrorCodeSpaceSecretNotFound space secret
	ErrorCodeSpaceSecretNotFound = "space_secret_not_found"

	// ErrorCodeTokenNotFound is const
	ErrorCodeTokenNotFound = "token_not_found"

	// ErrorFailedToDeleteToken is const
	ErrorFailedToDeleteToken = "failed_to_delete_token"

	// ErrorCodeRepoNotFound is const
	ErrorCodeRepoNotFound = "repo_not_found"

	// ErrorCodeOrganizationNotFound is const
	ErrorCodeOrganizationNotFound = "organization_not_found"

	// ErrorCodeCountExceeded is const
	ErrorCodeCountExceeded = "count_exceeded"

	// ErrorCodeSpaceAppNotFound space app
	ErrorCodeSpaceAppNotFound = "space_app_not_found"

	// ErrorCodeSpaceAppCreateFailed is const
	ErrorCodeSpaceAppCreateFailed = "space_app_create_failed"

	// ErrorCodeSpaceCommitConflict is const
	ErrorCodeSpaceCommitConflict = "space_commit_conflict"

	// ErrorCodeSpaceAppUnmatchedStatus is const
	ErrorCodeSpaceAppUnmatchedStatus = "space_app_unmatched_status"

	// ErrorCodeSpaceAppRestartOverTime is const
	ErrorCodeSpaceAppRestartOverTime = "space_app_restart_over_time"

	// ErrorCodeSpaceAppPauseFailed pause space app
	ErrorCodeSpaceAppPauseFailed = "space_app_pause_failed"

	// ErrorCodeCompAccountException comp account exception when pause space app
	ErrorCodeCompAccountException = "comp_account_exception"

	// ErrorCodeSpaceAppResumeFailed pause space app
	ErrorCodeSpaceAppResumeFailed = "space_app_resume_failed"

	// ErrorCodeSpaceAppResumeOverTime is const
	ErrorCodeSpaceAppResumeOverTime = "space_app_resume_over_time"

	// ErrorCodeSpaceAppSleepFailed sleep space app
	ErrorCodeSpaceAppSleepFailed = "space_app_sleep_failed"

	// ErrorCodeSpaceAppWakeupFailed sleep space app
	ErrorCodeSpaceAppWakeupFailed = "space_app_wakeup_failed" // #nosec G101

	// ErrorCodeAccessTokenInvalid This error code is for restful api
	ErrorCodeAccessTokenInvalid = "access_token_invalid"

	// ErrorCodeCSRFTokenMissing is const
	ErrorCodeCSRFTokenMissing = "csrf_token_missing" // #nosec G101

	// ErrorCodeCSRFTokenInvalid is const
	ErrorCodeCSRFTokenInvalid = "csrf_token_invalid" // #nosec G101

	// ErrorCodeCSRFTokenNotFound is const
	ErrorCodeCSRFTokenNotFound = "csrf_token_not_found" // #nosec G101

	// ErrorCodeSessionInvalid is const
	ErrorCodeSessionInvalid = "session_invalid"

	// ErrorCodeSessionIdInvalid is const
	ErrorCodeSessionIdInvalid = "session_id_invalid"

	// ErrorCodeSessionIdMissing is const
	ErrorCodeSessionIdMissing = "session_id_missing"

	// ErrorCodeSessionNotFound is const
	ErrorCodeSessionNotFound = "session_not_found"

	// ErrorCodeBranchExist is const
	ErrorCodeBranchExist = "branch_exist"

	// ErrorCodeBranchNotExist failed to find branch
	ErrorCodeBranchNotExist = "branch_not_exist"

	// ErrorCodeBranchInavtive is const
	ErrorCodeBranchInavtive = "branch_inactive"

	// ErrorCodeBaseBranchNotFound is const
	ErrorCodeBaseBranchNotFound = "base_branch_not_found"

	// ErrorCodeOrgExistResource is const
	ErrorCodeOrgExistResource = "org_resource_exist"

	// ErrorCodeInvalidParam is const
	errorCodeInvalidParam = "invalid_param"

	// ErrorEmailError is const
	ErrorEmailError = "email_error"

	// ErrorEmailCodeError is const
	ErrorEmailCodeError = "email_verify_code_error"

	// ErrorEmailCodeInvalid is const
	ErrorEmailCodeInvalid = "email_verify_code_invalid"

	// ErrorCodeNeedBindEmail is const
	ErrorCodeNeedBindEmail = "user_no_email"

	// ErrorCodeUserDuplicateBind is const
	ErrorCodeUserDuplicateBind = "user_duplicate_bind"

	// ErrorVerifyEmailFailed is const
	ErrorVerifyEmailFailed = "verify_email_failed"

	// ErrorVerifyEmailGitError is const
	ErrorVerifyEmailGitError = "verify_email_git_error"

	// ErrorCodeEmailDuplicateBind is const
	ErrorCodeEmailDuplicateBind = "email_duplicate_bind"

	// ErrorCodeEmailVerifyFailed is const
	ErrorCodeEmailVerifyFailed = "email_verify_failed"

	// ErrorCodeEmailDuplicateSend is const
	ErrorCodeEmailDuplicateSend = "email_duplicate_send"

	// ErrorCodeDisAgreedPrivacy is const
	ErrorCodeDisAgreedPrivacy = "disagreed_privacy"

	// ErrorCodeExpired is const
	ErrorCodeExpired = "expired"

	// ErrorCodePrivilegeOrgIdMismatch is const
	ErrorCodePrivilegeOrgIdMismatch = "privilege_org_id_mismatch"

	// ErrorCodeNotInPrivilegeOrg is const
	ErrorCodeNotInPrivilegeOrg = "not_in_privilege_org"

	// ErrorCodeInsufficientQuota user has insufficient quota balance
	ErrorCodeInsufficientQuota = "insufficient_quota"

	// ErrorCodeNoUsedQuota user is not currently using any quota
	ErrorCodeNoUsedQuota = "no_used_quota"

	// ErrorCodeNoNpuPermission user has no npu permission
	ErrorCodeNoNpuPermission = "no_npu_permission"

	// ErrorCodeComputilityAccountFindError find computility account error
	ErrorCodeComputilityAccountFindError = "computility_account_find_error"

	// ErrorCodeComputilityOrgFindError find computility org error
	ErrorCodeComputilityOrgFindError = "computility_org_find_error"

	// ErrorCodeComputilityOrgUpdateError update computility org error
	ErrorCodeComputilityOrgUpdateError = "computility_org_update_error"

	// ErrorComputilityOrgQuotaLowerBoundError quota count lower bound error
	ErrorComputilityOrgQuotaLowerBoundError = "computility_org_quota_lower_bound_error"

	// ErrorComputilityOrgQuotaMultipleError quota count not a multiple of default quota
	ErrorComputilityOrgQuotaMultipleError = "computility_org_quota_multiple_error"

	// ErrorBaseCase is const
	ErrorBaseCase = "internal_error"

	// ErrorMsgPublishFailed is const
	ErrorMsgPublishFailed = "msg_publish_failed"

	// ErrorDuplicateCreating duplicate creating
	ErrorDuplicateCreating = "duplicate_creating"

	// ErrorFailedGetOwnerInfo failed to get owner info
	ErrorFailedGetOwnerInfo = "failed_to_get_owner_info"

	// ErrorFailGetPlatformUser failed to get platform user
	ErrorFailGetPlatformUser = "failed_to_get_platform_user"

	// ErrorFailedCreateOrg failed to create org
	ErrorFailedCreateOrg = "failed_to_create_org"

	// ErrorFailedCreateToOrg failed to create to org
	ErrorFailedCreateToOrg = "failed_to_create_to_org"

	// ErrorFailSaveOrgMember failed to save org member
	ErrorFailSaveOrgMember = "failed_to_save_org_member"

	// ErrorSystemError system error
	ErrorSystemError = "system_error"

	// ErrorOrgNameRequesterAllEmpty is const
	ErrorOrgNameRequesterAllEmpty = "org_name_requester_all_empty"

	// ErrorOverOrgnameInviteeInviter is const
	ErrorOverOrgnameInviteeInviter = "over_orgname_invitee_inviter"

	// ErrorMemberInvitationParamAllEmpty is const
	ErrorMemberInvitationParamAllEmpty = "member_invitation_param_all_empty"

	// ErrorNoInvitationFound no invitation
	ErrorNoInvitationFound = "no_invitation_found"

	// ErrorFailedToDeleteUser failed to delete user in db
	ErrorFailedToDeleteUser = "failed_to_delete_user"

	// ErrorFailedToGetUserInfo failed to get user info
	ErrorFailedToGetUserInfo = "failed_to_get_user_info"

	// ErrorFailedToRevokePrivacy failed to revoke privacy
	ErrorFailedToRevokePrivacy = "failed_to_revoke_privacy"

	// ErrorFailedToAgreePrivacy failed to agree privacy
	ErrorFailedToAgreePrivacy = "failed_to_agree_privacy"

	// ErrorNameAlreadyBeenTaken name %s is already been taken
	ErrorNameAlreadyBeenTaken = "name_is_already_been_taken"

	// ErrorAccountCannotDeleteTheOrg account %s can't delete the org
	ErrorAccountCannotDeleteTheOrg = "account_can_not_delete_the_org"

	// ErrorInviteBadRequest
	ErrorInviteBadRequest = "user_has_been_invited"
	// email
	ErrorUserBadEmail = "user_email_address_error"

	// ErrorFailedToGetOrg failed to get org when get org by user, %w
	ErrorFailedToGetOrg = "failed_to_get_org"

	// ErrorMissingName missing name when creating token
	ErrorMissingName = "missing_name"

	// ErrorMissingAccount missing account when creating token
	ErrorMissingAccount = "missing_account"

	// ErrorFailedToRemoveMember failed to remove member
	ErrorFailedToRemoveMember = "failed_to_remove_member"

	// ErrorUserAlreadyInOrg the user is already a member of the org
	ErrorUserAlreadyInOrg = "the_user_is_already_a_member_of_the_org"

	// ErrorOrgNotAllowRequestMember org not allow request member
	ErrorOrgNotAllowRequestMember = "org_not_allow_request_member"

	// ErrorInvalidActorName invalid actor name
	ErrorInvalidActorName = "invalid_actor_name"

	// ErrorOrgFullnameIsEmpty org fullname is empty
	ErrorOrgFullnameIsEmpty = "org_fullname_is_empty"

	// ErrorInvalidAccount invalid account
	ErrorInvalidAccount = "invalid_account"

	// ErrorInvalidOrg invalid org
	ErrorInvalidOrg = "invalid_org"

	// ErrorInvalidActor invalid actor
	ErrorInvalidActor = "invalid_actor"

	// ErrorInvalidUser invalid user
	ErrorInvalidUser = "invalid_user"

	// ErrorInvalidRequester invalid requester
	ErrorInvalidRequester = "invalid_requester"

	// ErrorFullnameCanNotBeEmpty fullname can't be empty
	ErrorFullnameCanNotBeEmpty = "fullname_can_not_be_empty"

	// ErrorFailedToUpdateUserInfo failed to update user info
	ErrorFailedToUpdateUserInfo = "failed_to_update_user_info"

	// ErrorFailedToUPdateGitUserInfo ailed to update git user info
	ErrorFailedToUPdateGitUserInfo = "failed_to_update_git_userinfo"

	// ErrorUsernameInvalid username invalid
	ErrorUsernameInvalid = "username_invalid"

	// ErrorFailedToGetPlatformUserInfo is const
	ErrorFailedToGetPlatformUserInfo = "failed_to_get_platform_user_info"

	// ErrorFailedToDeleteUserInGitServer is const
	ErrorFailedToDeleteUserInGitServer = "failed_to_delete_user_in_git_server"

	// ErrorUserAlreadyRequestedToBeDelete is const
	ErrorUserAlreadyRequestedToBeDelete = "user_already_requested_to_be_delete"

	// ErrorFailedToCreateToken failed to create token
	ErrorFailedToCreateToken = "failed_to_create_token"

	// ErrorFailedToEcryptToken failed to ecrypt token
	ErrorFailedToEcryptToken = "failed_to_ecrypt_token"

	// ErrorInputParamIsEmpty input param is empty
	ErrorInputParamIsEmpty = "input_param_is_empty"

	// ErrorDeleteTokenParamIsEmpty delete token param is empty
	ErrorDeleteTokenParamIsEmpty = "delete_token_param_is_empty"

	// ErrorFailedToSaveOrg failed to save org
	ErrorFailedToSaveOrg = "failed_to_save_org"

	// ErrorNothingChanged nothing changed
	ErrorNothingChanged = "nothing_changed"

	// ErrorFailedToGetOrgInfo failed to get org info
	ErrorFailedToGetOrgInfo = "failed_to_get_org_info"

	// ErrorFailedToGetMemberInfo failed to get member info
	ErrorFailedToGetMemberInfo = "failed_to_get_member_info"

	// ErrorFailedToSaveMemberForAddingMember is const
	ErrorFailedToSaveMemberForAddingMember = "failed_to_save_member_for_adding_member"

	// ErrorOnlyOwnerCanNotBeRemoved the only owner can not be removed
	ErrorOnlyOwnerCanNotBeRemoved = "the_only_owner_can_not_be_removed"

	// ErrorFailedToValidateCmd failed to validate cmd
	ErrorFailedToValidateCmd = "failed_to_validate_cmd"

	// ErrorFailedToDeleteGitMember failed to delete git member
	ErrorFailedToDeleteGitMember = "failed_to_delete_git_member"

	// ErrorFailedToDeleteMember failed to delete member
	ErrorFailedToDeleteMember = "failed_to_delete_member"

	// ErrorFailedToChangeOwnerOfOrg failed to change owner of org
	ErrorFailedToChangeOwnerOfOrg = "failed_to_change_owner_of_org"

	// ErrorFailedToGetMembersByOrgName failed to get members by org name: %s, %s
	ErrorFailedToGetMembersByOrgName = "failed_to_get_members_by_org_name"

	// ErrorUserAccountIsAlreadyAMemberOfOrgAccount is const
	ErrorUserAccountIsAlreadyAMemberOfOrgAccount = "user_account_is_already_a_member_of_the_org_account"

	// ErrorFailedToAddMemberToOrg failed to add member:%s to org:%s
	ErrorFailedToAddMemberToOrg = "failed_to_add_member"

	// ErrorInvalidStatus invalid status %s
	ErrorInvalidStatus = "invalid_status"

	// ErrorUsernameIsAlreadyTaken user name %s is already taken
	ErrorUsernameIsAlreadyTaken = "user_name_is_already_taken"

	// ErrorFailedToCreatePlatformUser failed to create platform user: %s
	ErrorFailedToCreatePlatformUser = "failed_to_create_platform_user"

	// ErrorFailToSaveUserInDb failed to save user in db: %s
	ErrorFailToSaveUserInDb = "failed_to_save_user"

	// ErrorCodeConcurrentUpdating failed to save org in db: %s
	ErrorCodeConcurrentUpdating = "concurrent_updating"

	// ErrorFailToRetrieveActivityData failed to retrieve activities: %s
	ErrorFailToRetrieveActivityData = "failed_to_retrieve_activity"

	ErrorCodeFailToCreateIssue = "failed_to_create_issue"
	ErrorCodeFailToUpdateIssue = "failed_to_update_issue"
	ErrorCodeIssueNotFound     = "issue_not_found"
	ErrorCodeIssueClosed       = "issue_closed"
	ErrorCodeIssueIsOpen       = "issue_is_open"

	ErrorCodeFailToCreateComment = "failed_to_create_comment"
	ErrorCodeFailToUpdateComment = "failed_to_update_comment"
	ErrorCodeFailToDeleteComment = "failed_to_delete_comment"
	ErrorCodeCommentNotFound     = "comment_not_found"
	ErrorCodeFailToCheckMentions = "fail_to_check_mentions"
	ErrorCodePicExceedMaxCount   = "pic_exceed_max_count"
	ErrorCodePicUrlNotAllowed    = "pic_url_not_allowed"
	ErrorCodePicSizeTooLarge     = "pic_size_too_large"

	// ErrorImageBlock failed to pass image moderation
	ErrorImageBlock = "error_image_block"

	// ErrorImageFailed failed to moderation image
	ErrorImageFailed = "error_image_failed"

	// ErrorTextFailed failed to moderation text
	ErrorTextFailed = "error_text_failed"

	// ErrorTextBlock failed to pass text moderation
	ErrorTextBlock = "error_text_block"

	// ErrorImageUnsupported image pixel not supported for moderation
	ErrorImageUnsupported = "error_image_unsupported"

	ErrorCodeDiscussionDisabled = "discussion_is_disabled"
	ErrorCodeDiscussionEnabled  = "discussion_is_enabled"

	// ErrorCodeModelCiUnmatchedStatus is const
	ErrorCodeModelCiUnmatchedStatus = "model_ci_unmatched_status"

	ErrorCodeModelPrivate = "model_is_private"

	ErrorCodeModelCiIsRunning = "model_ci_is_running"

	ErrorCodeModelCiNotFound = "model_ci_not_found"

	ErrorCodeModelCiNoTestCase = "model_ci_no_test_case"

	// error code for contest
	ErrorContestApplyGetContestError = "contest_apply_get_contest_error"
	ErrorContestOverError            = "contest_over_error"
	ErrorContestApplyStatusError     = "contest_apply_status_error"
	ErrorContestWorkNotFoundError    = "contest_work_not_found_error"
	ErrorContestWorkDeleteError      = "contest_work_delete_error"
	ErrorContestSubmitWorkError      = "contest_submit_work_error"
)

// errorImpl
type errorImpl struct {
	code     string
	msg      string
	innerErr error // error info for diagnostic
}

// Error returns the error message.
//
// This function returns the error message of the errorImpl struct.
//
// No parameters.
// Returns a string representing the error message.
func (e errorImpl) Error() string {
	return e.msg
}

// ErrorCode returns the error code.
//
// This function returns the error code of the errorImpl struct.
// The error code is a string representing the type of the error, it could be used for error handling and diagnostic.
//
// No parameters.
// Returns a string representing the error code.
func (e errorImpl) ErrorCode() string {
	return e.code
}

// InnerError returns the inner error.
type InnerError interface {
	InnerError() error
}

// InnerErr returns the inner error.
func (e errorImpl) InnerError() error {
	return e.innerErr
}

// InnerErr returns the inner error.
func InnerErr(err error) error {
	var v InnerError
	if ok := errors.As(err, &v); ok {
		return v.InnerError()
	}

	return err
}

// New creates a new error with the specified code and message.
//
// This function creates a new errorImpl struct with the specified code, message and error info
// for diagnostic. If the message is empty, the function will replace all "_" in the code with
// " " as the message.
//
// Parameters:
//
//	code: a string representing the type of the error
//	msg: a string representing the error message, which is returned to client or end user
//	err: error info for diagnostic, which is used for diagnostic by developers
//
// Returns an errorImpl struct.
func New(code string, msg string, err error) errorImpl {
	v := errorImpl{
		code:     code,
		innerErr: err,
	}

	if msg == "" {
		v.msg = strings.ReplaceAll(code, "_", " ")
	} else {
		v.msg = msg
	}

	return v
}

// notfoudError
type notfoudError struct {
	errorImpl
}

// NotFound is a marker method for a not found error.
func (e notfoudError) NotFound() {}

// NewNotFound creates a new not found error with the specified code and message.
func NewNotFound(code string, msg string, err error) notfoudError {
	return notfoudError{errorImpl: New(code, msg, err)}
}

// IsNotFound checks if the given error is a not found error.
func IsNotFound(err error) (notfoudError, bool) {
	if err == nil {
		return notfoudError{}, false
	}
	var v notfoudError
	ok := errors.As(err, &v)

	return v, ok
}

// IsUserDuplicateBind checks if the given error is a user duplicate bind error.
func IsUserDuplicateBind(err error) bool {
	if err == nil {
		return false
	}

	var e errorImpl
	if ok := errors.As(err, &e); ok {
		if e.ErrorCode() == ErrorCodeUserDuplicateBind {
			return true
		}
	}

	return false
}

// noPermissionError
type noPermissionError struct {
	errorImpl
}

// NoPermission is a marker method for a "no permission" error.
func (e noPermissionError) NoPermission() {}

// NewNoPermission creates a new "no permission" error with the specified message.
func NewNoPermission(msg string, err error) noPermissionError {
	return noPermissionError{errorImpl: New(errorCodeNoPermission, msg, err)}
}

// IsNoPermission checks if the given error is a "no permission" error.
func IsNoPermission(err error) bool {
	if err == nil {
		return false
	}

	ok := errors.As(err, &noPermissionError{})

	return ok
}

// NewInvalidParam creates a new error with the specified invalid parameter message.
func NewInvalidParam(msg string, err error) errorImpl {
	return New(errorCodeInvalidParam, msg, err)
}

// NewCountExceeded creates a new error with the specified count exceeded message.
func NewCountExceeded(msg string, err error) errorImpl {
	return New(ErrorCodeCountExceeded, msg, err)
}

// limitRateError
type limitRateError struct {
	errorImpl
}

// OverLimit is a marker method for over limit rate error.
func (l limitRateError) OverLimit() {}

// NewOverLimit creates a new over limit error with the specified code and message.
func NewOverLimit(code string, msg string, err error) limitRateError {
	return limitRateError{errorImpl: New(code, msg, err)}
}

// NewExpired checks if the given error is a over limit error.
func NewExpired(msg string, err error) errorImpl {
	return New(ErrorCodeExpired, msg, err)
}

// NewCommonRespError creates a new error with the common resp error.
func NewCommonRespError(msg string, err error) errorImpl {
	return New(ErrorBaseCase, msg, err)
}

// resourceDisabledError
type resourceDisabledError struct {
	errorImpl
}

// NewResourceDisabled creates a new resource disable error with the specified code and message.
func NewResourceDisabled(code string, msg string, err error) resourceDisabledError {
	return resourceDisabledError{errorImpl: New(code, msg, err)}
}

// resourcePrivateError
type resourcePrivateError struct {
	errorImpl
}

// NewResourcePrivate creates a new resource private error with the specified code and message.
func NewResourcePrivate(code string, msg string, err error) resourcePrivateError {
	return resourcePrivateError{errorImpl: New(code, msg, err)}
}

// modelCiIsRunningError
type modelCiIsRunningError struct {
	errorImpl
}

// NewModelCiIsRunningError creates a new model ci is running error with the specified code and message.
func NewModelCiIsRunningError(code string, msg string, err error) modelCiIsRunningError {
	return modelCiIsRunningError{errorImpl: New(code, msg, err)}
}

// modelCiNoTestCaseError
type modelCiNoTestCaseError struct {
	errorImpl
}

// NewModelCiNoTestCaseError creates a new model ci no test case with the specified code and message.
func NewModelCiNoTestCaseError(code string, msg string, err error) modelCiNoTestCaseError {
	return modelCiNoTestCaseError{errorImpl: New(code, msg, err)}
}
