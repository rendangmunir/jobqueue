package main

import (
	"context"
	"jobqueue/config"
	"jobqueue/delivery/graphql"
	_dataloader "jobqueue/delivery/graphql/dataloader"
	"jobqueue/delivery/graphql/mutation"
	"jobqueue/delivery/graphql/query"
	"jobqueue/delivery/graphql/schema"
	"jobqueue/entity"
	"jobqueue/pkg/handler"
	"jobqueue/pkg/server"
	"jobqueue/queue"
	inmemrepo "jobqueue/repository/inmem"
	"jobqueue/service"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"

	"github.com/labstack/echo"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	setupLogger()
	logger := logrus.New()
	logger.SetReportCaller(true)
	e := server.New(config.Data.Server)
	e.Echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${remote_ip} ${time_rfc3339_nano} \"${method} ${path}\" ${status} ${bytes_out} \"${referer}\" \"${user_agent}\"\n",
	}))
	e.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.OPTIONS},
	}))

	//graphql schema
	opts := make([]_graphql.SchemaOpt, 0)
	opts = append(opts, _graphql.SubscribeResolverTimeout(10*time.Second))

	//initialize in mem database
	inMemDb := make(map[string]*entity.Job)

	//set job repository
	jobRepository := inmemrepo.
		NewJobRepository().
		SetInMemConnection(inMemDb).
		Build()
	dataloader := _dataloader.
		New().
		SetJobRepository(jobRepository).
		SetBatchFunction().
		Build()
		
	// Set up job queue
	jobQueue := queue.NewJobQueue(150)

	//set job service
	jobService := service.NewJobService().
		SetJobRepository(jobRepository).
		SetJobQueue(jobQueue).
		Build()
	
	jobHandler := queue.CreateJobHandler(jobService)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	jobQueue.StartWorkers(ctx, 5, jobHandler)

	jobMutation := mutation.NewJobMutation(jobService, dataloader)
	jobQuery := query.NewJobQuery(jobService, dataloader)

	rootResolver := graphql.
		New().
		SetJobMutation(jobMutation).
		SetJobQuery(jobQuery).
		Build()

	graphqlSchema := _graphql.MustParseSchema(schema.String(), rootResolver, opts...)
	e.Echo.POST("/graphql",
		handler.GraphQLHandler(&relay.Handler{Schema: graphqlSchema}),
		dataloader.EchoMiddelware,
	)
	e.Echo.GET("/graphql",
		handler.GraphQLHandler(&relay.Handler{Schema: graphqlSchema}),
		dataloader.EchoMiddelware,
	)
	e.Echo.GET("/graphiql", handler.GraphiQLHandler)
	// Handle graceful shutdown
	go func() {
		if err := e.Start(); err != nil && err != http.ErrServerClosed  {
			e.Echo.Logger.Fatal("shutting down the server")
		}
	}()

	// Listen for SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Stop server
	if err := e.Echo.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}

	// Wait for workers
	jobQueue.Wait()
	e.Echo.Logger.Fatal(e.Start())
}

func setupLogger() {
	configLogger := zap.NewDevelopmentConfig()
	configLogger.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	configLogger.DisableStacktrace = true
	logger, _ := configLogger.Build()
	zap.ReplaceGlobals(logger)
}
