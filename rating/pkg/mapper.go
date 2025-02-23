package ratingmodel

import (
	"github.com/karaMuha/go-movie/pb"
	"github.com/karaMuha/go-movie/pkg/dtos"
)

func RespErrToProto(respErr *dtos.RespErr) *pb.ResponseStatus {
	return &pb.ResponseStatus{
		StatusCode: int32(respErr.StatusCode),
		Message:    respErr.StatusMessage,
	}
}
