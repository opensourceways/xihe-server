package app

const (
	ErrorCodeSytem = "system"

	// It exceed the max times for a day.
	ErrorCodeAIQuestionExceedMaxTimes           = "aiquestion_exceed_max_times"
	ErrorCodeAIQuestionSubmissionExpiry         = "aiquestion_submission_expiry"
	ErrorCodeAIQuestionSubmissionUnmatchedTimes = "aiquestion_submission_unmatched_times"

	ErrorBigModelSensitiveInfo     = "bigmodel_sensitive_info"
	ErrorBigModelRecourseBusy      = "bigmodel_resource_busy"
	ErrorBigModelConcurrentRequest = "bigmodel_concurrent_request"

	ErrorWuKongNoPicture        = "bigmodel_no_wukong_picture"
	ErrorWuKongInvalidId        = "wukong_invalid_id"
	ErrorWuKongInvalidOwner     = "wukong_invalid_owner"
	ErrorWuKongInvalidPath      = "wukong_invalid_path"
	ErrorWuKongNoAuthorization  = "wukong_no_authorization"
	ErrorWuKongInvalidLink      = "wukong_invalid_link"
	ErrorWuKongDuplicateLike    = "wukong_duplicate_like"
	ErrorWuKongExccedMaxLikeNum = "wukong_excced_max_like_num"
)
