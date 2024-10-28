package management

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/dezh-tech/immortal/management/docs"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	Echo     *echo.Echo
	DataBase *mongo.Client
	Config   *Config
}

func NewServer(cfg Config) (*Server, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.DatabaseConfig.ConnectionTimeout)*time.Millisecond)
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.DatabaseConfig.URI).
		SetServerAPIOptions(serverAPI).
		SetConnectTimeout(time.Duration(cfg.DatabaseConfig.ConnectionTimeout) * time.Millisecond).
		SetBSONOptions(&options.BSONOptions{
			UseJSONStructTags: true,
			NilSliceAsEmpty:   true,
		})

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Server{
		Echo:     echo.New(),
		Config:   &cfg,
		DataBase: client,
	}, nil
}

//	@title			Immortal Management API
//	@version		1.0
//	@description	Immortal Management API provide some useful api to manage the relay and other services related to Immortal

//	@contact.name	Dezh Technologies
//	@contact.url	https://dezh.tech
//	@contact.email	hi@Dezh.tech

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

// @BasePath	/
func (server *Server) Start() error {

	url := net.JoinHostPort(server.Config.Bind, //nolint
		strconv.Itoa(int(server.Config.Port)))

	docs.SwaggerInfo.Host = url

	return server.Echo.Start(url)
}
