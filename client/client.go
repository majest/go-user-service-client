package client

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	srv "github.com/majest/go-microservice/server"
	"github.com/majest/user-service/pb"
	"github.com/majest/user-service/server"
	"golang.org/x/net/context"
)

type client struct {
	Context         context.Context
	Logger          log.Logger
	FindOneEndpoint endpoint.Endpoint
}

// CreateClient returns a client to our Demo service
func CreateClient(ctx context.Context, logger log.Logger, config *srv.ClientConfig) (server.UserService, error) {

	// create and return the client
	return client{
		Context: ctx,
		Logger:  logger,
		FindOneEndpoint: srv.MakeEndpoint(
			"FindOne",
			logger,
			EncodeFindOneRequest,
			DecodeFindOneResponse,
			pb.UserResponse{},
			// span.Sum,
			// trace.Sum,
			config,
		),
	}, nil
}

func (c client) FindOne(u *server.UserFindOneRequest) (*server.UserResponse, error) {
	response, err := c.FindOneEndpoint(c.Context, u)
	if err != nil {
		c.Logger.Log("findOne", err.Error())
	}

	return response.(*server.UserResponse), err
}
