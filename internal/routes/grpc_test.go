package routes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	gomock "go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"soft.structx.io/dino/logging"
	pb "soft.structx.io/dino/pb/routes/v1"
	"soft.structx.io/dino/setup"
)

type RouteServerSuite struct {
	suite.Suite

	svc         Service
	routeServer pb.RouteServiceServer

	testApp *fxtest.App
}

func (suite *RouteServerSuite) SetupSuite() {
	tb := suite.T()

	suite.svc = nil

	type params struct {
		fx.In

		Logger *zap.Logger
	}

	var opts = fx.Options(
		setup.Module,
		logging.Module,
		fx.Provide(func(p params) pb.RouteServiceServer {
			return newRouteServer(p.Logger, suite.svc)
		}),
		fx.Populate(&suite.routeServer),
	)
	fxtest.New(tb, opts).RequireStart()
}

func (suite *RouteServerSuite) TestCreateRoute() {
	tb := suite.T()
	ctx := tb.Context()

	assert := suite.Assert()

	ctlr := gomock.NewController(tb)
	mockSvc := NewMockService(ctlr)
	mockSvc.EXPECT().Create(
		gomock.Any(),
		gomock.AssignableToTypeOf(RouteCreate{}),
	).Return(Route{
		ID:                  "xyz",
		Tunnel:              "test-tunnel",
		Hostname:            "hello.world.local",
		DestinationProtocol: "http",
		DestinationIP:       "127.0.0.1",
		DestinationPort:     8080,
		Enabled:             true,
		CreatedAt:           time.Now(),
		UpdatedAt:           nil,
	}, nil).Times(1)

	tt := []struct {
		expected error
		req      *pb.CreateRouteRequest
	}{
		{
			expected: nil,
			req: &pb.CreateRouteRequest{
				Create: &pb.RouteCreate{
					Tunnel:       "test-tunnel",
					Hostname:     "hello.world.local",
					DestProtocol: "http",
					DestAddr:     "127.0.0.1",
					DestPort:     8080,
				},
			},
		},
	}

	for _, tc := range tt {
		resp, actual := suite.routeServer.CreateRoute(ctx, tc.req)
		assert.Equal(tc.expected, actual)

		if actual == nil {
			assert.NotNil(resp.Route)

			assert.NotEmpty(resp.Route.Uid)

			assert.Equal(tc.req.Create.Tunnel, resp.Route.Tunnel)
		}
	}
}

func (suite *RouteServerSuite) TeardownSuite() {
	suite.testApp.RequireStop()
}

func TestRouteServerSuite(t *testing.T) {
	suite.Run(t, new(RouteServerSuite))
}
