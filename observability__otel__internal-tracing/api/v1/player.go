package v1

import (
	"context"
	"net/http"
	"thanhldt060802/common/tracer"
	"thanhldt060802/model"
	"thanhldt060802/service"
	"time"

	authMdw "thanhldt060802/middleware/auth"

	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type apiPlayer struct {
	tracer        trace.Tracer
	playerService service.IPlayerService
}

func RegisterAPIExample(api hureg.APIGen, playerService service.IPlayerService) {
	handler := &apiPlayer{
		tracer:        otel.Tracer("api/v1/player.go"),
		playerService: playerService,
	}

	apiGroup := api.AddBasePath("/player")

	hureg.Register(
		apiGroup,
		huma.Operation{
			OperationID: "player-get-by-id",
			Method:      http.MethodGet,
			Path:        "/{player_uuid}",
			Security:    authMdw.DefaultAuthSecurity,
			Description: "Get player by id.",
			Middlewares: huma.Middlewares{authMdw.NewAuthMiddleware(api)},
		},
		handler.GetById,
	)
}

type GetPlayerByIdResponse struct {
	Body struct {
		Data *model.Player `json:"data" doc:"Player data"`
	}
}

func (handler *apiPlayer) GetById(ctx context.Context, req *struct {
	PlayerUuid string `path:"player_uuid" format:"uuid" doc:"Player uuid"`
}) (res *GetPlayerByIdResponse, err error) {
	ctx, span := tracer.StartSpanInternal(ctx)
	defer span.End()

	span.AddEvent("Call to playerService.GetById")
	time.Sleep(1 * time.Second)

	player, err := handler.playerService.GetById(ctx, req.PlayerUuid)
	if err != nil {
		span.Err = err
		return
	}

	res = &GetPlayerByIdResponse{}
	res.Body.Data = player
	return
}
