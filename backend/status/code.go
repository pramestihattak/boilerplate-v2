package status

import (
	pb "boilerplate-v2/gen/status"
)

type StatusCode string

const (
	EmptyCode     StatusCode = ""
	StatusCodeLen int        = 6
)

const (
	Success_Generic StatusCode = "SC2000"

	UserErrCode_Generic                       StatusCode = "SC4000"
	UserErrCode_AccountExist                  StatusCode = "SC4001"
	UserErrCode_VerificationMissingParameters StatusCode = "SC4002"
	UserErrCode_AccountNotFound               StatusCode = "SC4003"

	SystemErrCode_Generic            StatusCode = "SC5000"
	SystemErrCode_FailedReadMetadata StatusCode = "SC5002"
	SystemErrCode_FailedSanitize     StatusCode = "SC5004"
	SystemErrCode_FailedToRegister   StatusCode = "SC5005"
	SystemErrCode_FailedToVerify     StatusCode = "SC5006"
)

var statusMap = map[StatusCode]pb.StatusResponse{
	// 2XXX - success
	Success_Generic: {StatusDesc: "completed"},

	// 4XXX - user errors
	UserErrCode_Generic:                       {StatusDesc: "generic user error"},
	UserErrCode_AccountExist:                  {StatusDesc: "account with this email is exist"},
	UserErrCode_VerificationMissingParameters: {StatusDesc: "missing parameters for verification"},
	UserErrCode_AccountNotFound:               {StatusDesc: "account not found"},

	// 5XXX - system errors
	SystemErrCode_Generic:            {StatusDesc: "generic system error"},
	SystemErrCode_FailedReadMetadata: {StatusDesc: "could not read metadata from context"},
	SystemErrCode_FailedSanitize:     {StatusDesc: "failed to sanitize parameters"},
	SystemErrCode_FailedToRegister:   {StatusDesc: "failed to registering user"},
	SystemErrCode_FailedToVerify:     {StatusDesc: "failed to verify user"},
}
