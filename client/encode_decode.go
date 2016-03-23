package client

import (
	"github.com/majest/user-service/pb"
	"github.com/majest/user-service/server"
	"golang.org/x/net/context"
)

func EncodeFindOneRequest(ctx context.Context, request interface{}) (interface{}, error) {
	req := request.(*server.UserFindOneRequest)
	return &pb.UserFindOneRequest{&pb.UserSearch{
		Id:        req.UserSearch.Id,
		FirstName: req.UserSearch.FirstName,
		LastName:  req.UserSearch.LastName,
		Email:     req.UserSearch.Email,
		PostCode:  req.UserSearch.PostCode,
	}}, nil
}

func DecodeFindOneResponse(ctx context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UserResponse)

	if resp.User != nil {
		return &server.UserResponse{&server.User{
			Id:        resp.User.Id,
			FirstName: resp.User.FirstName,
			LastName:  resp.User.LastName,
			Email:     resp.User.Email,
			Address:   resp.User.Address,
			Street:    resp.User.Street,
			PostCode:  resp.User.PostCode,
			City:      resp.User.City,
			Country:   resp.User.Country,
			Phone:     resp.User.Phone,
		}}, nil
	}

	return &server.UserResponse{}, nil
}
