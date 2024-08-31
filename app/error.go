package app

import "errors"

type ErrorExceedMaxRelatedResourceNum struct {
	error
}

type ErrorPrivateRepo struct {
	error
}

type ErrorExccedMaxTrainingRecord struct {
	error
}

type ErrorOnlyOneRunningTraining struct {
	error
}

type ErrorUnavailableRepoFile struct {
	error
}

type ErrorUpdateLFSFile struct {
	error
}

type ErrorPreviewLFSFile struct {
	error
}

type errorUnavailableTraining struct {
	error
}

func IsErrorUnavailableTraining(err error) bool {
	return errors.As(err, &errorUnavailableTraining{})
}

type ErrorDuplicateTrainingName struct {
	error
}

const (
	ErrorCodeSystem = "system"

	// It exceed the max times for a day.
	ErrorCodeAIQuestionExceedMaxTimes           = "aiquestion_exceed_max_times"
	ErrorCodeAIQuestionSubmissionExpiry         = "aiquestion_submission_expiry"
	ErrorCodeAIQuestionSubmissionUnmatchedTimes = "aiquestion_submission_unmatched_times"

	ErrorRepoFileTooManyFilesToDelete = "repofile_too_many_files_to_delete"

	ErrorCompetitionDuplicateSubmission = "competition_duplicate_submission"

	ErrorBigModelSensitiveInfo = "bigmodel_sensitive_info"
	ErrorBigModelRecourseBusy  = "bigmodel_resource_busy"

	ErrorTrainNoLog        = "train_no_log"
	ErrorTrainNoOutput     = "train_no_output"
	ErrorTrainNotFound     = "train_not_found"
	ErrorTrainExccedMaxNum = "train_excced_max_num" // excced max training num for a user

	ErrorWuKongInvalidId        = "wukong_invalid_id"
	ErrorWuKongInvalidOwner     = "wukong_invalid_owner"
	ErrorWuKongInvalidPath      = "wukong_invalid_path"
	ErrorWuKongNoAuthorization  = "wukong_no_authorization"
	ErrorWuKongInvalidLink      = "wukong_invalid_link"
	ErrorWuKongDuplicateLike    = "wukong_duplicate_like"
	ErrorWuKongExccedMaxLikeNum = "wukong_excced_max_like_num"

	ErrorFinetuneExpiry           = "finetune_expiry"
	ErrorFinetuneNotFound         = "finetune_not_found"
	ErrorFinetuneExccedMaxNum     = "finetune_excced_max_num"
	ErrorFinetuneNoPermission     = "finetune_no_permission"
	ErrorFinetuneRunningJobExists = "finetune_running_job_exists"

	ErrorAICCFinetuneNoLog    = "aicc_finetune_no_log"
	ErrorAICCFinetuneNotFound = "aicc_finetune_not_found"
)
