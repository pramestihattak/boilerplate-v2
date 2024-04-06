package status

import (
	pb "boilerplate-v2/gen/status"

	"google.golang.org/grpc/codes"
	grpc_status "google.golang.org/grpc/status"
)

func ResponseFromCode(code StatusCode) *pb.StatusResponse {
	var res pb.StatusResponse
	if len(code) == StatusCodeLen {
		if statusMap[code].StatusDesc != "" {
			res.StatusDesc = statusMap[code].StatusDesc
			res.StatusCode = string(code)
			return &res
		}
	}

	res.StatusDesc = statusMap[SystemErrCode_Generic].StatusDesc
	res.StatusCode = string(SystemErrCode_Generic)
	return &res
}

// ResponseFromCodeToErr accepts StatusCode and converts it to an grpc error message.
func ResponseFromCodeToErr(code StatusCode) error {
	var (
		grpcCode    codes.Code
		grpcMessage string
	)
	statRes := ResponseFromCode(code)
	grpcMessage = statRes.StatusDesc

	// Response func performs validation and ensures the status code is valid and adheres to the format
	// get the status code prefix 2/4/5
	switch string(statRes.StatusCode[2]) {
	case "4":
		// this is user error
		// return invalid argument
		grpcCode = codes.InvalidArgument
	default: // includes "5"
		grpcCode = codes.Internal
	}

	st := grpc_status.New(grpcCode, grpcMessage)
	if st2, err := st.WithDetails((statRes)); err == nil {
		st = st2 // if err != nil, this just skips adding the details on the response.
	}
	return st.Err()
}
