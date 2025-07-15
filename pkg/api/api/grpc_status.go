package api

import (
	"errors"
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pkgErr "github.com/bucketeer-io/bucketeer/pkg/error"
)

// ToDo: Once multilingual support has been moved to the frontend, delete localizedMessage.
func NewGRPCStatus(
	err error,
	errorDomain string,
	localizedMessage *errdetails.LocalizedMessage,
	anotherDetailData ...map[string]string,
) *status.Status {
	var reason string
	var messageKey string
	var metadatas []map[string]string
	var st *status.Status
	var invalidAugmentError *pkgErr.ErrorInvalidAugment

	if errors.As(err, &invalidAugmentError) {
		msg := fmt.Sprintf("%s: invalid augment", invalidAugmentError.Pkg)
		st = status.New(codes.InvalidArgument, msg)
		reason = "INVALID_AUGMENT"
		messageKey = fmt.Sprintf("%s.invalid_augment", invalidAugmentError.Pkg)
		for _, arg := range invalidAugmentError.InvalidArgs {
			metadatas = append(metadatas, map[string]string{
				"messageKey": messageKey,
				"field":      arg,
			})
		}
	} else {
		reason = "UNKNOWN"
		messageKey = "unknown"
	}
	// when adding multiple details
	for _, md := range anotherDetailData {
		for k, v := range md {
			metadatas = append(metadatas, map[string]string{
				"messageKey": messageKey,
				k:            v,
			})
		}
	}

	for _, md := range metadatas {
		st, err = st.WithDetails(&errdetails.ErrorInfo{
			Reason:   reason,
			Domain:   errorDomain,
			Metadata: md,
		})
		if err != nil {
			return status.New(codes.Internal, err.Error())
		}
	}
	if localizedMessage != nil {
		st, err = st.WithDetails(localizedMessage)
		if err != nil {
			return status.New(codes.Internal, err.Error())
		}
	}
	return st
}
