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
	UserErrCode_LoginWrongPassword            StatusCode = "SC4004"
	UserErrCode_MissingBearerToken            StatusCode = "SC4005"

	SystemErrCode_Generic             StatusCode = "SC5000"
	SystemErrCode_FailedReadMetadata  StatusCode = "SC5002"
	SystemErrCode_FailedSanitize      StatusCode = "SC5004"
	SystemErrCode_FailedToRegister    StatusCode = "SC5005"
	SystemErrCode_FailedToVerify      StatusCode = "SC5006"
	SystemErrCode_FailedToLogin       StatusCode = "SC5007"
	SystemErrCode_FailedToGetAuthData StatusCode = "SC5008"
	SystemErrCode_FailedToSetHeader   StatusCode = "SC5009"
)

var statusMap = map[StatusCode]pb.StatusResponse{
	// 2XXX - success
	Success_Generic: {StatusDesc: "completed"},

	// 4XXX - user errors
	UserErrCode_Generic:                       {StatusDesc: "generic user error"},
	UserErrCode_AccountExist:                  {StatusDesc: "account with this email is exist"},
	UserErrCode_VerificationMissingParameters: {StatusDesc: "missing parameters for verification"},
	UserErrCode_AccountNotFound:               {StatusDesc: "account not found"},
	UserErrCode_LoginWrongPassword:            {StatusDesc: "wrong password"},
	UserErrCode_MissingBearerToken:            {StatusDesc: "missing bearer token"},

	// 5XXX - system errors
	SystemErrCode_Generic:             {StatusDesc: "generic system error"},
	SystemErrCode_FailedReadMetadata:  {StatusDesc: "could not read metadata from context"},
	SystemErrCode_FailedSanitize:      {StatusDesc: "failed to sanitize parameters"},
	SystemErrCode_FailedToRegister:    {StatusDesc: "failed to registering user"},
	SystemErrCode_FailedToVerify:      {StatusDesc: "failed to verify user"},
	SystemErrCode_FailedToLogin:       {StatusDesc: "failed to login"},
	SystemErrCode_FailedToGetAuthData: {StatusDesc: "failed to get auth data"},
	SystemErrCode_FailedToSetHeader:   {StatusDesc: "failed to set header"},
}
