This file is a merged representation of the entire codebase, combined into a single document by Repomix.
The content has been processed where security check has been disabled.

# File Summary

## Purpose
This file contains a packed representation of the entire repository's contents.
It is designed to be easily consumable by AI systems for analysis, code review,
or other automated processes.

## File Format
The content is organized as follows:
1. This summary section
2. Repository information
3. Directory structure
4. Repository files (if enabled)
5. Multiple file entries, each consisting of:
  a. A header with the file path (## File: path/to/file)
  b. The full contents of the file in a code block

## Usage Guidelines
- This file should be treated as read-only. Any changes should be made to the
  original repository files, not this packed version.
- When processing this file, use the file path to distinguish
  between different files in the repository.
- Be aware that this file may contain sensitive information. Handle it with
  the same level of security as you would the original repository.

## Notes
- Some files may have been excluded based on .gitignore rules and Repomix's configuration
- Binary files are not included in this packed representation. Please refer to the Repository Structure section for a complete list of file paths, including binary files
- Files matching patterns in .gitignore are excluded
- Files matching default ignore patterns are excluded
- Security check has been disabled - content may contain sensitive information
- Files are sorted by Git change count (files with more changes are at the bottom)

# Directory Structure
```
api/
  v1/
    errors.go
    user.go
    v1.go
cmd/
  migration/
    wire/
      wire_gen.go
      wire.go
    main.go
  server/
    wire/
      wire_gen.go
      wire.go
    main.go
  task/
    wire/
      wire_gen.go
      wire.go
    main.go
config/
  local.yml
  prod.yml
deploy/
  build/
    Dockerfile
  docker-compose/
    docker-compose.yml
docs/
  docs.go
  swagger.json
  swagger.yaml
internal/
  handler/
    handler.go
    user.go
  job/
    job.go
    user.go
  middleware/
    cors.go
    jwt.go
    log.go
    sign.go
  model/
    user.go
  repository/
    repository.go
    user.go
  router/
    router.go
    user.go
  server/
    http.go
    job.go
    migration.go
    task.go
  service/
    service.go
    user.go
  task/
    task.go
    user.go
pkg/
  app/
    app.go
  config/
    config.go
  jwt/
    jwt.go
  log/
    log.go
  server/
    grpc/
      grpc.go
    http/
      http.go
    server.go
  sid/
    convert.go
    sid.go
  zapgorm2/
    zapgorm2.go
scripts/
  README.md
test/
  mocks/
    repository/
      repository.go
      user.go
    service/
      user.go
  server/
    handler/
      main_test.go
      user_test.go
    repository/
      user_test.go
    service/
      user_test.go
web/
  index.html
.dockerignore
.gitattributes
.gitignore
go.mod
LICENSE
Makefile
README_zh.md
README.md
```

# Files

## File: api/v1/errors.go
```go
package v1

var (
	// common errors
	ErrSuccess             = newError(0, "ok")
	ErrBadRequest          = newError(400, "Bad Request")
	ErrUnauthorized        = newError(401, "Unauthorized")
	ErrNotFound            = newError(404, "Not Found")
	ErrInternalServerError = newError(500, "Internal Server Error")

	// more biz errors
	ErrEmailAlreadyUse = newError(1001, "The email is already in use.")
)
```

## File: api/v1/user.go
```go
package v1

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}
type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
}
type LoginResponse struct {
	Response
	Data LoginResponseData
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" example:"alan"`
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
}
type GetProfileResponseData struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname" example:"alan"`
}
type GetProfileResponse struct {
	Response
	Data GetProfileResponseData
}
```

## File: api/v1/v1.go
```go
package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HandleSuccess(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	resp := Response{Code: errorCodeMap[ErrSuccess], Message: ErrSuccess.Error(), Data: data}
	if _, ok := errorCodeMap[ErrSuccess]; !ok {
		resp = Response{Code: 0, Message: "", Data: data}
	}
	ctx.JSON(http.StatusOK, resp)
}

func HandleError(ctx *gin.Context, httpCode int, err error, data interface{}) {
	if data == nil {
		data = map[string]string{}
	}
	resp := Response{Code: errorCodeMap[err], Message: err.Error(), Data: data}
	if _, ok := errorCodeMap[err]; !ok {
		resp = Response{Code: 500, Message: "unknown error", Data: data}
	}
	ctx.JSON(httpCode, resp)
}

type Error struct {
	Code    int
	Message string
}

var errorCodeMap = map[error]int{}

func newError(code int, msg string) error {
	err := errors.New(msg)
	errorCodeMap[err] = code
	return err
}
func (e Error) Error() string {
	return e.Message
}
```

## File: cmd/migration/wire/wire_gen.go
```go
// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/go-nunu/nunu-layout-advanced/internal/repository"
	"github.com/go-nunu/nunu-layout-advanced/internal/server"
	"github.com/go-nunu/nunu-layout-advanced/pkg/app"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Injectors from wire.go:

func NewWire(viperViper *viper.Viper, logger *log.Logger) (*app.App, func(), error) {
	db := repository.NewDB(viperViper, logger)
	migrateServer := server.NewMigrateServer(db, logger)
	appApp := newApp(migrateServer)
	return appApp, func() {
	}, nil
}

// wire.go:

var repositorySet = wire.NewSet(repository.NewDB, repository.NewRepository, repository.NewUserRepository)

var serverSet = wire.NewSet(server.NewMigrateServer)

// build App
func newApp(
	migrateServer *server.MigrateServer,
) *app.App {
	return app.NewApp(app.WithServer(migrateServer), app.WithName("demo-migrate"))
}
```

## File: cmd/migration/wire/wire.go
```go
//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/go-nunu/nunu-layout-advanced/internal/repository"
	"github.com/go-nunu/nunu-layout-advanced/internal/server"
	"github.com/go-nunu/nunu-layout-advanced/pkg/app"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	repository.NewRepository,
	repository.NewUserRepository,
)
var serverSet = wire.NewSet(
	server.NewMigrateServer,
)

// build App
func newApp(
	migrateServer *server.MigrateServer,
) *app.App {
	return app.NewApp(
		app.WithServer(migrateServer),
		app.WithName("demo-migrate"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serverSet,
		newApp,
	))
}
```

## File: cmd/migration/main.go
```go
package main

import (
	"context"
	"flag"
	"github.com/go-nunu/nunu-layout-advanced/cmd/migration/wire"
	"github.com/go-nunu/nunu-layout-advanced/pkg/config"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
)

func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	logger := log.NewLog(conf)

	app, cleanup, err := wire.NewWire(conf, logger)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}
```

## File: cmd/server/wire/wire_gen.go
```go
// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/go-nunu/nunu-layout-advanced/internal/handler"
	"github.com/go-nunu/nunu-layout-advanced/internal/job"
	"github.com/go-nunu/nunu-layout-advanced/internal/repository"
	"github.com/go-nunu/nunu-layout-advanced/internal/router"
	"github.com/go-nunu/nunu-layout-advanced/internal/server"
	"github.com/go-nunu/nunu-layout-advanced/internal/service"
	"github.com/go-nunu/nunu-layout-advanced/pkg/app"
	"github.com/go-nunu/nunu-layout-advanced/pkg/jwt"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"github.com/go-nunu/nunu-layout-advanced/pkg/server/http"
	"github.com/go-nunu/nunu-layout-advanced/pkg/sid"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Injectors from wire.go:

func NewWire(viperViper *viper.Viper, logger *log.Logger) (*app.App, func(), error) {
	jwtJWT := jwt.NewJwt(viperViper)
	handlerHandler := handler.NewHandler(logger)
	db := repository.NewDB(viperViper, logger)
	repositoryRepository := repository.NewRepository(logger, db)
	transaction := repository.NewTransaction(repositoryRepository)
	sidSid := sid.NewSid()
	serviceService := service.NewService(transaction, logger, sidSid, jwtJWT)
	userRepository := repository.NewUserRepository(repositoryRepository)
	userService := service.NewUserService(serviceService, userRepository)
	userHandler := handler.NewUserHandler(handlerHandler, userService)
	routerDeps := router.RouterDeps{
		Logger:      logger,
		Config:      viperViper,
		JWT:         jwtJWT,
		UserHandler: userHandler,
	}
	httpServer := server.NewHTTPServer(routerDeps)
	jobJob := job.NewJob(transaction, logger, sidSid)
	userJob := job.NewUserJob(jobJob, userRepository)
	jobServer := server.NewJobServer(logger, userJob)
	appApp := newApp(httpServer, jobServer)
	return appApp, func() {
	}, nil
}

// wire.go:

var repositorySet = wire.NewSet(repository.NewDB, repository.NewRepository, repository.NewTransaction, repository.NewUserRepository)

var serviceSet = wire.NewSet(service.NewService, service.NewUserService)

var handlerSet = wire.NewSet(handler.NewHandler, handler.NewUserHandler)

var jobSet = wire.NewSet(job.NewJob, job.NewUserJob)

var serverSet = wire.NewSet(server.NewHTTPServer, server.NewJobServer)

// build App
func newApp(
	httpServer *http.Server,
	jobServer *server.JobServer,

) *app.App {
	return app.NewApp(app.WithServer(httpServer, jobServer), app.WithName("demo-server"))
}
```

## File: cmd/server/wire/wire.go
```go
//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/go-nunu/nunu-layout-advanced/internal/handler"
	"github.com/go-nunu/nunu-layout-advanced/internal/job"
	"github.com/go-nunu/nunu-layout-advanced/internal/repository"
	"github.com/go-nunu/nunu-layout-advanced/internal/router"
	"github.com/go-nunu/nunu-layout-advanced/internal/server"
	"github.com/go-nunu/nunu-layout-advanced/internal/service"
	"github.com/go-nunu/nunu-layout-advanced/pkg/app"
	"github.com/go-nunu/nunu-layout-advanced/pkg/jwt"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"github.com/go-nunu/nunu-layout-advanced/pkg/server/http"
	"github.com/go-nunu/nunu-layout-advanced/pkg/sid"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	//repository.NewMongo,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
)

var jobSet = wire.NewSet(
	job.NewJob,
	job.NewUserJob,
)
var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJobServer,
)

// build App
func newApp(
	httpServer *http.Server,
	jobServer *server.JobServer,
	// task *server.Task,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, jobServer),
		app.WithName("demo-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		jobSet,
		serverSet,
		wire.Struct(new(router.RouterDeps), "*"),
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))
}
```

## File: cmd/server/main.go
```go
package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/go-nunu/nunu-layout-advanced/cmd/server/wire"
	"github.com/go-nunu/nunu-layout-advanced/pkg/config"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"go.uber.org/zap"
)

// @title           Nunu Example API
// @version         1.0.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8000
// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	logger := log.NewLog(conf)

	app, cleanup, err := wire.NewWire(conf, logger)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	logger.Info("server start", zap.String("host", fmt.Sprintf("http://%s:%d", conf.GetString("http.host"), conf.GetInt("http.port"))))
	logger.Info("docs addr", zap.String("addr", fmt.Sprintf("http://%s:%d/swagger/index.html", conf.GetString("http.host"), conf.GetInt("http.port"))))
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}
```

## File: cmd/task/wire/wire_gen.go
```go
// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/go-nunu/nunu-layout-advanced/internal/repository"
	"github.com/go-nunu/nunu-layout-advanced/internal/server"
	"github.com/go-nunu/nunu-layout-advanced/internal/task"
	"github.com/go-nunu/nunu-layout-advanced/pkg/app"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"github.com/go-nunu/nunu-layout-advanced/pkg/sid"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Injectors from wire.go:

func NewWire(viperViper *viper.Viper, logger *log.Logger) (*app.App, func(), error) {
	db := repository.NewDB(viperViper, logger)
	repositoryRepository := repository.NewRepository(logger, db)
	transaction := repository.NewTransaction(repositoryRepository)
	sidSid := sid.NewSid()
	taskTask := task.NewTask(transaction, logger, sidSid)
	userRepository := repository.NewUserRepository(repositoryRepository)
	userTask := task.NewUserTask(taskTask, userRepository)
	taskServer := server.NewTaskServer(logger, userTask)
	appApp := newApp(taskServer)
	return appApp, func() {
	}, nil
}

// wire.go:

var repositorySet = wire.NewSet(repository.NewDB, repository.NewRepository, repository.NewTransaction, repository.NewUserRepository)

var taskSet = wire.NewSet(task.NewTask, task.NewUserTask)

var serverSet = wire.NewSet(server.NewTaskServer)

// build App
func newApp(task2 *server.TaskServer,
) *app.App {
	return app.NewApp(app.WithServer(task2), app.WithName("demo-task"))
}
```

## File: cmd/task/wire/wire.go
```go
//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/go-nunu/nunu-layout-advanced/internal/repository"
	"github.com/go-nunu/nunu-layout-advanced/internal/server"
	"github.com/go-nunu/nunu-layout-advanced/internal/task"
	"github.com/go-nunu/nunu-layout-advanced/pkg/app"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"github.com/go-nunu/nunu-layout-advanced/pkg/sid"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
)

var taskSet = wire.NewSet(
	task.NewTask,
	task.NewUserTask,
)
var serverSet = wire.NewSet(
	server.NewTaskServer,
)

// build App
func newApp(
	task *server.TaskServer,
) *app.App {
	return app.NewApp(
		app.WithServer(task),
		app.WithName("demo-task"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		taskSet,
		serverSet,
		newApp,
		sid.NewSid,
	))
}
```

## File: cmd/task/main.go
```go
package main

import (
	"context"
	"flag"
	"github.com/go-nunu/nunu-layout-advanced/cmd/task/wire"
	"github.com/go-nunu/nunu-layout-advanced/pkg/config"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
)

func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	logger := log.NewLog(conf)
	logger.Info("start task")
	app, cleanup, err := wire.NewWire(conf, logger)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}

}
```

## File: config/local.yml
```yaml
env: local
http:
  #  host: 0.0.0.0
  host: 127.0.0.1
  port: 8000
security:
  api_sign:
    app_key: 123456
    app_security: 123456
  jwt:
    key: QQYnRFerJTSEcrfB89fw8prOaObmrch8
data:
  db:
    user:
      driver: sqlite
      dsn: storage/nunu-test.db?_busy_timeout=5000
  #    user:
  #      driver: mysql
  #      dsn: root:123456@tcp(127.0.0.1:3380)/user?charset=utf8mb4&parseTime=True&loc=Local
  #    user:
  #      driver: postgres
  #      dsn: host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai
  redis:
    addr: 127.0.0.1:6350
    password: ""
    db: 0
    read_timeout: 0.2s
    write_timeout: 0.2s
  mongo:
    uri: mongodb://root:123456@localhost:27017
log:
  log_level: debug
  mode: both               #  file or console or both
  encoding: console           # json or console
  log_file_name: "./storage/logs/server.log"
  max_backups: 30
  max_age: 7
  max_size: 1024
  compress: true
```

## File: config/prod.yml
```yaml
env: prod
http:
  host: 0.0.0.0
  #  host: 127.0.0.1
  port: 8000
security:
  api_sign:
    app_key: 123456
    app_security: 123456
  jwt:
    key: QQYnRFerJTSEcrfB89fw8prOaObmrch8
data:
  db:
    user:
      driver: sqlite
      dsn: storage/nunu-test.db?_busy_timeout=5000
  #    user:
  #      driver: mysql
  #      dsn: root:123456@tcp(127.0.0.1:3380)/user?charset=utf8mb4&parseTime=True&loc=Local
  #    user:
  #      driver: postgres
  #      dsn: host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai
  redis:
    addr: 127.0.0.1:6350
    password: ""
    db: 0
    read_timeout: 0.2s
    write_timeout: 0.2s
  mongo:
    uri: mongodb://root:123456@localhost:27017
log:
  log_level: info
  mode: file               #  file or console or both
  encoding: json           # json or console
  log_file_name: "./storage/logs/server.log"
  max_backups: 30
  max_age: 7
  max_size: 1024
  compress: true
```

## File: deploy/build/Dockerfile
```
ARG REGISTRY=docker.io
FROM ${REGISTRY}/golang:1.19-alpine AS builder
RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

ARG APP_RELATIVE_PATH

COPY .. /data/app
WORKDIR /data/app

RUN rm -rf /data/app/bin/
RUN export GOPROXY=https://goproxy.cn,direct && go mod tidy && go build -ldflags="-s -w" -o ./bin/server ${APP_RELATIVE_PATH}
RUN mv config /data/app/bin/


FROM ${REGISTRY}/alpine:3.16
RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories


RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata


ARG APP_ENV
ENV APP_ENV=${APP_ENV}

WORKDIR /data/app
COPY --from=builder /data/app/bin /data/app

EXPOSE 8000
ENTRYPOINT [ "./server" ]

#docker build -t  1.1.1.1:5000/demo-api:v1 --build-arg APP_CONF=config/prod.yml --build-arg  APP_RELATIVE_PATH=./cmd/server/...  .
#docker run -it --rm --entrypoint=ash 1.1.1.1:5000/demo-api:v1
```

## File: deploy/docker-compose/docker-compose.yml
```yaml
version: '3'

services:
  user-db:
    image: mysql:8.0.31-debian
    hostname: user-db
    container_name: user-db
    ports:
      - 3380:3306
    environment:
      - MYSQL_ROOT_PASSWORD=123456
      - MYSQL_ROOT_HOST=%
      - MYSQL_DATABASE=user
    # volumes:
    #   - ./data/mysql/user:/var/lib/mysql
    #   - ./conf/mysql/conf.d:/etc/mysql/conf.d
  cache-redis:
    image: redis:6-alpine
    hostname: cache-redis
    # volumes:
    #   - ./data/redis/cache/:/data
    #   - ./conf/redis/cache/redis.conf:/etc/redis/redis.conf
    ports:
      - 6350:6379
    command: ["redis-server","/etc/redis/redis.conf"]
```

## File: docs/docs.go
```go
// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "账号登录",
                "parameters": [
                    {
                        "description": "params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.LoginResponse"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "目前只支持邮箱登录",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "用户注册",
                "parameters": [
                    {
                        "description": "params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.Response"
                        }
                    }
                }
            }
        },
        "/user": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "获取用户信息",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.GetProfileResponse"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "修改用户信息",
                "parameters": [
                    {
                        "description": "params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.UpdateProfileRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "v1.GetProfileResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "$ref": "#/definitions/v1.GetProfileResponseData"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "v1.GetProfileResponseData": {
            "type": "object",
            "properties": {
                "nickname": {
                    "type": "string",
                    "example": "alan"
                },
                "userId": {
                    "type": "string"
                }
            }
        },
        "v1.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "1234@gmail.com"
                },
                "password": {
                    "type": "string",
                    "example": "123456"
                }
            }
        },
        "v1.LoginResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "$ref": "#/definitions/v1.LoginResponseData"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "v1.LoginResponseData": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                }
            }
        },
        "v1.RegisterRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "1234@gmail.com"
                },
                "password": {
                    "type": "string",
                    "example": "123456"
                }
            }
        },
        "v1.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "v1.UpdateProfileRequest": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "1234@gmail.com"
                },
                "nickname": {
                    "type": "string",
                    "example": "alan"
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "localhost:8000",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Nunu Example API",
	Description:      "This is a sample server celler server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
```

## File: docs/swagger.json
```json
{
    "swagger": "2.0",
    "info": {
        "description": "This is a sample server celler server.",
        "title": "Nunu Example API",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "1.0.0"
    },
    "host": "localhost:8000",
    "paths": {
        "/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "账号登录",
                "parameters": [
                    {
                        "description": "params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.LoginResponse"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "目前只支持邮箱登录",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "用户注册",
                "parameters": [
                    {
                        "description": "params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.Response"
                        }
                    }
                }
            }
        },
        "/user": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "获取用户信息",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.GetProfileResponse"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "修改用户信息",
                "parameters": [
                    {
                        "description": "params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.UpdateProfileRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "v1.GetProfileResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "$ref": "#/definitions/v1.GetProfileResponseData"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "v1.GetProfileResponseData": {
            "type": "object",
            "properties": {
                "nickname": {
                    "type": "string",
                    "example": "alan"
                },
                "userId": {
                    "type": "string"
                }
            }
        },
        "v1.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "1234@gmail.com"
                },
                "password": {
                    "type": "string",
                    "example": "123456"
                }
            }
        },
        "v1.LoginResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "$ref": "#/definitions/v1.LoginResponseData"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "v1.LoginResponseData": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                }
            }
        },
        "v1.RegisterRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "1234@gmail.com"
                },
                "password": {
                    "type": "string",
                    "example": "123456"
                }
            }
        },
        "v1.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "v1.UpdateProfileRequest": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "1234@gmail.com"
                },
                "nickname": {
                    "type": "string",
                    "example": "alan"
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}
```

## File: docs/swagger.yaml
```yaml
definitions:
  v1.GetProfileResponse:
    properties:
      code:
        type: integer
      data:
        $ref: '#/definitions/v1.GetProfileResponseData'
      message:
        type: string
    type: object
  v1.GetProfileResponseData:
    properties:
      nickname:
        example: alan
        type: string
      userId:
        type: string
    type: object
  v1.LoginRequest:
    properties:
      email:
        example: 1234@gmail.com
        type: string
      password:
        example: "123456"
        type: string
    required:
    - email
    - password
    type: object
  v1.LoginResponse:
    properties:
      code:
        type: integer
      data:
        $ref: '#/definitions/v1.LoginResponseData'
      message:
        type: string
    type: object
  v1.LoginResponseData:
    properties:
      accessToken:
        type: string
    type: object
  v1.RegisterRequest:
    properties:
      email:
        example: 1234@gmail.com
        type: string
      password:
        example: "123456"
        type: string
    required:
    - email
    - password
    type: object
  v1.Response:
    properties:
      code:
        type: integer
      data: {}
      message:
        type: string
    type: object
  v1.UpdateProfileRequest:
    properties:
      email:
        example: 1234@gmail.com
        type: string
      nickname:
        example: alan
        type: string
    required:
    - email
    type: object
host: localhost:8000
info:
  contact:
    email: support@swagger.io
    name: API Support
    url: http://www.swagger.io/support
  description: This is a sample server celler server.
  license:
    name: Apache 2.0
    url: http://www.apache.org/licenses/LICENSE-2.0.html
  termsOfService: http://swagger.io/terms/
  title: Nunu Example API
  version: 1.0.0
paths:
  /login:
    post:
      consumes:
      - application/json
      parameters:
      - description: params
        in: body
        name: request
        required: true
        schema:
          $ref: '#/definitions/v1.LoginRequest'
      produces:
      - application/json
      responses:
        "200":
          description: OK
          schema:
            $ref: '#/definitions/v1.LoginResponse'
      summary: 账号登录
      tags:
      - 用户模块
  /register:
    post:
      consumes:
      - application/json
      description: 目前只支持邮箱登录
      parameters:
      - description: params
        in: body
        name: request
        required: true
        schema:
          $ref: '#/definitions/v1.RegisterRequest'
      produces:
      - application/json
      responses:
        "200":
          description: OK
          schema:
            $ref: '#/definitions/v1.Response'
      summary: 用户注册
      tags:
      - 用户模块
  /user:
    get:
      consumes:
      - application/json
      produces:
      - application/json
      responses:
        "200":
          description: OK
          schema:
            $ref: '#/definitions/v1.GetProfileResponse'
      security:
      - Bearer: []
      summary: 获取用户信息
      tags:
      - 用户模块
    put:
      consumes:
      - application/json
      parameters:
      - description: params
        in: body
        name: request
        required: true
        schema:
          $ref: '#/definitions/v1.UpdateProfileRequest'
      produces:
      - application/json
      responses:
        "200":
          description: OK
          schema:
            $ref: '#/definitions/v1.Response'
      security:
      - Bearer: []
      summary: 修改用户信息
      tags:
      - 用户模块
securityDefinitions:
  Bearer:
    in: header
    name: Authorization
    type: apiKey
swagger: "2.0"
```

## File: internal/handler/handler.go
```go
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-nunu/nunu-layout-advanced/pkg/jwt"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
)

type Handler struct {
	logger *log.Logger
}

func NewHandler(
	logger *log.Logger,
) *Handler {
	return &Handler{
		logger: logger,
	}
}
func GetUserIdFromCtx(ctx *gin.Context) string {
	v, exists := ctx.Get("claims")
	if !exists {
		return ""
	}
	return v.(*jwt.MyCustomClaims).UserId
}
```

## File: internal/handler/user.go
```go
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-nunu/nunu-layout-advanced/api/v1"
	"github.com/go-nunu/nunu-layout-advanced/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type UserHandler struct {
	*Handler
	userService service.UserService
}

func NewUserHandler(handler *Handler, userService service.UserService) *UserHandler {
	return &UserHandler{
		Handler:     handler,
		userService: userService,
	}
}

// Register godoc
// @Summary 用户注册
// @Schemes
// @Description 目前只支持邮箱登录
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.RegisterRequest true "params"
// @Success 200 {object} v1.Response
// @Router /register [post]
func (h *UserHandler) Register(ctx *gin.Context) {
	req := new(v1.RegisterRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.Register(ctx, req); err != nil {
		h.logger.WithContext(ctx).Error("userService.Register error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// Login godoc
// @Summary 账号登录
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body v1.LoginRequest true "params"
// @Success 200 {object} v1.LoginResponse
// @Router /login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var req v1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	token, err := h.userService.Login(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}
	v1.HandleSuccess(ctx, v1.LoginResponseData{
		AccessToken: token,
	})
}

// GetProfile godoc
// @Summary 获取用户信息
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.GetProfileResponse
// @Router /user [get]
func (h *UserHandler) GetProfile(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	user, err := h.userService.GetProfile(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, user)
}

// UpdateProfile godoc
// @Summary 修改用户信息
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UpdateProfileRequest true "params"
// @Success 200 {object} v1.Response
// @Router /user [put]
func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)

	var req v1.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.UpdateProfile(ctx, userId, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}
```

## File: internal/job/job.go
```go
package job

import (
	"github.com/go-nunu/nunu-layout-advanced/internal/repository"
	"github.com/go-nunu/nunu-layout-advanced/pkg/jwt"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"github.com/go-nunu/nunu-layout-advanced/pkg/sid"
)

type Job struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewJob(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
) *Job {
	return &Job{
		logger: logger,
		sid:    sid,
		tm:     tm,
	}
}
```

## File: internal/job/user.go
```go
package job

import (
	"context"
	"github.com/go-nunu/nunu-layout-advanced/internal/repository"
	"time"
)

type UserJob interface {
	KafkaConsumer(ctx context.Context) error
}

func NewUserJob(
	job *Job,
	userRepo repository.UserRepository,
) UserJob {
	return &userJob{
		userRepo: userRepo,
		Job:      job,
	}
}

type userJob struct {
	userRepo repository.UserRepository
	*Job
}

func (t userJob) KafkaConsumer(ctx context.Context) error {
	// do something
	for {
		t.logger.Info("KafkaConsumer")
		time.Sleep(time.Second * 5)
	}
}
```

## File: internal/middleware/cors.go
```go
package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.Header("Access-Control-Allow-Methods", c.GetHeader("Access-Control-Request-Method"))
			c.Header("Access-Control-Allow-Headers", c.GetHeader("Access-Control-Request-Headers"))
			c.Header("Access-Control-Max-Age", "7200")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
```

## File: internal/middleware/jwt.go
```go
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-nunu/nunu-layout-advanced/api/v1"
	"github.com/go-nunu/nunu-layout-advanced/pkg/jwt"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"go.uber.org/zap"
	"net/http"
)

func StrictAuth(j *jwt.JWT, logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			logger.WithContext(ctx).Warn("No token", zap.Any("data", map[string]interface{}{
				"url":    ctx.Request.URL,
				"params": ctx.Params,
			}))
			v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}

		claims, err := j.ParseToken(tokenString)
		if err != nil {
			logger.WithContext(ctx).Error("token error", zap.Any("data", map[string]interface{}{
				"url":    ctx.Request.URL,
				"params": ctx.Params,
			}), zap.Error(err))
			v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		recoveryLoggerFunc(ctx, logger)
		ctx.Next()
	}
}

func NoStrictAuth(j *jwt.JWT, logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.Request.Header.Get("Authorization")
		if tokenString == "" {
			tokenString, _ = ctx.Cookie("accessToken")
		}
		if tokenString == "" {
			tokenString = ctx.Query("accessToken")
		}
		if tokenString == "" {
			ctx.Next()
			return
		}

		claims, err := j.ParseToken(tokenString)
		if err != nil {
			ctx.Next()
			return
		}

		ctx.Set("claims", claims)
		recoveryLoggerFunc(ctx, logger)
		ctx.Next()
	}
}

func recoveryLoggerFunc(ctx *gin.Context, logger *log.Logger) {
	if userInfo, ok := ctx.MustGet("claims").(*jwt.MyCustomClaims); ok {
		logger.WithValue(ctx, zap.String("UserId", userInfo.UserId))
	}
}
```

## File: internal/middleware/log.go
```go
package middleware

import (
	"bytes"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/random"
	"github.com/gin-gonic/gin"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"go.uber.org/zap"
	"io"
	"time"
)

func RequestLogMiddleware(logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// The configuration is initialized once per request
		uuid, err := random.UUIdV4()
		if err != nil {
			return
		}
		trace := cryptor.Md5String(uuid)
		logger.WithValue(ctx, zap.String("trace", trace))
		logger.WithValue(ctx, zap.String("request_method", ctx.Request.Method))
		logger.WithValue(ctx, zap.Any("request_headers", ctx.Request.Header))
		logger.WithValue(ctx, zap.String("request_url", ctx.Request.URL.String()))
		if ctx.Request.Body != nil {
			bodyBytes, _ := ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 关键点
			logger.WithValue(ctx, zap.String("request_params", string(bodyBytes)))
		}
		logger.WithContext(ctx).Info("Request")
		ctx.Next()
	}
}
func ResponseLogMiddleware(logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw
		startTime := time.Now()
		ctx.Next()
		duration := time.Since(startTime).String()
		logger.WithContext(ctx).Info("Response", zap.Any("response_body", blw.body.String()), zap.Any("time", duration))
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
```

## File: internal/middleware/sign.go
```go
package middleware

import (
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/gin-gonic/gin"
	v1 "github.com/go-nunu/nunu-layout-advanced/api/v1"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"github.com/spf13/viper"
	"net/http"
	"sort"
	"strings"
)

func SignMiddleware(logger *log.Logger, conf *viper.Viper) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requiredHeaders := []string{"Timestamp", "Nonce", "Sign", "App-Version"}

		for _, header := range requiredHeaders {
			value, ok := ctx.Request.Header[header]
			if !ok || len(value) == 0 {
				v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
				ctx.Abort()
				return
			}
		}

		data := map[string]string{
			"AppKey":     conf.GetString("security.api_sign.app_key"),
			"Timestamp":  ctx.Request.Header.Get("Timestamp"),
			"Nonce":      ctx.Request.Header.Get("Nonce"),
			"AppVersion": ctx.Request.Header.Get("App-Version"),
		}

		var keys []string
		for k := range data {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool { return strings.ToLower(keys[i]) < strings.ToLower(keys[j]) })

		var str string
		for _, k := range keys {
			str += k + data[k]
		}
		str += conf.GetString("security.api_sign.app_security")

		if ctx.Request.Header.Get("Sign") != strings.ToUpper(cryptor.Md5String(str)) {
			v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
```

## File: internal/model/user.go
```go
package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id        uint   `gorm:"primarykey"`
	UserId    string `gorm:"unique;not null"`
	Nickname  string `gorm:"not null"`
	Password  string `gorm:"not null"`
	Email     string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) TableName() string {
	return "users"
}
```

## File: internal/repository/repository.go
```go
package repository

import (
	"context"
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"github.com/go-nunu/nunu-layout-advanced/pkg/zapgorm2"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

const ctxTxKey = "TxKey"

type Repository struct {
	db *gorm.DB
	//rdb    *redis.Client
	//mongo  *mongo.Client
	logger *log.Logger
}

func NewRepository(
	logger *log.Logger,
	db *gorm.DB,
	// rdb *redis.Client,
	//
	//	mongo *mongo.Client,
) *Repository {
	return &Repository{
		db: db,
		//rdb:    rdb,
		//mongo:  mongo,
		logger: logger,
	}
}

type Transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewTransaction(r *Repository) Transaction {
	return r
}

// DB return tx
// If you need to create a Transaction, you must call DB(ctx) and Transaction(ctx,fn)
func (r *Repository) DB(ctx context.Context) *gorm.DB {
	v := ctx.Value(ctxTxKey)
	if v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx
		}
	}
	return r.db.WithContext(ctx)
}

func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, ctxTxKey, tx)
		return fn(ctx)
	})
}

func NewDB(conf *viper.Viper, l *log.Logger) *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)

	logger := zapgorm2.New(l.Logger)
	driver := conf.GetString("data.db.user.driver")
	dsn := conf.GetString("data.db.user.dsn")

	// GORM doc: https://gorm.io/docs/connecting_to_the_database.html
	switch driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger,
		})
	case "postgres":
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		panic("unknown db driver")
	}
	if err != nil {
		panic(err)
	}
	db = db.Debug()

	// Connection Pool config
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}
func NewRedis(conf *viper.Viper) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.GetString("data.redis.addr"),
		Password: conf.GetString("data.redis.password"),
		DB:       conf.GetInt("data.redis.db"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("redis error: %s", err.Error()))
	}

	return rdb
}
func NewMongo(conf *viper.Viper) (*mongo.Client, func(), error) {
	// https://www.mongodb.com/zh-cn/docs/drivers/go/current/
	uri := conf.GetString("data.mongo.uri")
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(fmt.Sprintf("mongo client error: %s", err.Error()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("mongo ping error: %s", err.Error()))
	}

	return client, func() {
		err = client.Disconnect(ctx)
		if err != nil {
			panic(fmt.Sprintf("mongo disconnect error: %s", err.Error()))
		}
	}, err
}
```

## File: internal/repository/user.go
```go
package repository

import (
	"context"
	"errors"
	v1 "github.com/go-nunu/nunu-layout-advanced/api/v1"
	"github.com/go-nunu/nunu-layout-advanced/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

func NewUserRepository(
	r *Repository,
) UserRepository {
	return &userRepository{
		Repository: r,
	}
}

type userRepository struct {
	*Repository
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.DB(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.DB(ctx).Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, userId string) (*model.User, error) {
	var user model.User
	if err := r.DB(ctx).Where("user_id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.DB(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
```

## File: internal/router/router.go
```go
package router

import (
	"github.com/go-nunu/nunu-layout-advanced/internal/handler"
	"github.com/go-nunu/nunu-layout-advanced/pkg/jwt"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"github.com/spf13/viper"
)

type RouterDeps struct {
	Logger      *log.Logger
	Config      *viper.Viper
	JWT         *jwt.JWT
	UserHandler *handler.UserHandler
}
```

## File: internal/router/user.go
```go
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-nunu/nunu-layout-advanced/internal/middleware"
)

func InitUserRouter(
	deps RouterDeps,
	r *gin.RouterGroup,
) {
	// No route group has permission
	noAuthRouter := r.Group("/")
	{
		noAuthRouter.POST("/register", deps.UserHandler.Register)
		noAuthRouter.POST("/login", deps.UserHandler.Login)
	}
	// Non-strict permission routing group
	noStrictAuthRouter := r.Group("/").Use(middleware.NoStrictAuth(deps.JWT, deps.Logger))
	{
		noStrictAuthRouter.GET("/user", deps.UserHandler.GetProfile)
	}

	// Strict permission routing group
	strictAuthRouter := r.Group("/").Use(middleware.StrictAuth(deps.JWT, deps.Logger))
	{
		strictAuthRouter.PUT("/user", deps.UserHandler.UpdateProfile)
	}
}
```

## File: internal/server/http.go
```go
package server

import (
	"github.com/gin-gonic/gin"
	apiV1 "github.com/go-nunu/nunu-layout-advanced/api/v1"
	"github.com/go-nunu/nunu-layout-advanced/docs"
	"github.com/go-nunu/nunu-layout-advanced/internal/middleware"
	"github.com/go-nunu/nunu-layout-advanced/internal/router"
	"github.com/go-nunu/nunu-layout-advanced/pkg/server/http"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewHTTPServer(
	deps router.RouterDeps,
) *http.Server {
	if deps.Config.GetString("env") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	s := http.NewServer(
		gin.Default(),
		deps.Logger,
		http.WithServerHost(deps.Config.GetString("http.host")),
		http.WithServerPort(deps.Config.GetInt("http.port")),
	)

	// swagger doc
	docs.SwaggerInfo.BasePath = "/v1"
	s.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		//ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", deps.Config.GetInt("app.http.port"))),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
	))

	s.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(deps.Logger),
		middleware.RequestLogMiddleware(deps.Logger),
		//middleware.SignMiddleware(log),
	)
	s.GET("/", func(ctx *gin.Context) {
		deps.Logger.WithContext(ctx).Info("hello")
		apiV1.HandleSuccess(ctx, map[string]interface{}{
			":)": "Thank you for using nunu!",
		})
	})

	v1 := s.Group("/v1")
	router.InitUserRouter(deps, v1)

	return s
}
```

## File: internal/server/job.go
```go
package server

import (
	"context"
	"github.com/go-nunu/nunu-layout-advanced/internal/job"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
)

type JobServer struct {
	log     *log.Logger
	userJob job.UserJob
}

func NewJobServer(
	log *log.Logger,
	userJob job.UserJob,
) *JobServer {
	return &JobServer{
		log:     log,
		userJob: userJob,
	}
}

func (j *JobServer) Start(ctx context.Context) error {
	// Tips: If you want job to start as a separate process, just refer to the task implementation and adjust the code accordingly.

	// eg: kafka consumer
	err := j.userJob.KafkaConsumer(ctx)
	return err
}
func (j *JobServer) Stop(ctx context.Context) error {
	return nil
}
```

## File: internal/server/migration.go
```go
package server

import (
	"context"
	"github.com/go-nunu/nunu-layout-advanced/internal/model"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

type MigrateServer struct {
	db  *gorm.DB
	log *log.Logger
}

func NewMigrateServer(db *gorm.DB, log *log.Logger) *MigrateServer {
	return &MigrateServer{
		db:  db,
		log: log,
	}
}
func (m *MigrateServer) Start(ctx context.Context) error {
	if err := m.db.AutoMigrate(
		&model.User{},
	); err != nil {
		m.log.Error("user migrate error", zap.Error(err))
		return err
	}
	m.log.Info("AutoMigrate success")
	os.Exit(0)
	return nil
}
func (m *MigrateServer) Stop(ctx context.Context) error {
	m.log.Info("AutoMigrate stop")
	return nil
}
```

## File: internal/server/task.go
```go
package server

import (
	"context"
	"github.com/go-co-op/gocron"
	"github.com/go-nunu/nunu-layout-advanced/internal/task"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"go.uber.org/zap"
	"time"
)

type TaskServer struct {
	log       *log.Logger
	scheduler *gocron.Scheduler
	userTask  task.UserTask
}

func NewTaskServer(
	log *log.Logger,
	userTask task.UserTask,
) *TaskServer {
	return &TaskServer{
		log:      log,
		userTask: userTask,
	}
}
func (t *TaskServer) Start(ctx context.Context) error {
	gocron.SetPanicHandler(func(jobName string, recoverData interface{}) {
		t.log.Error("TaskServer Panic", zap.String("job", jobName), zap.Any("recover", recoverData))
	})

	// eg: crontab task
	t.scheduler = gocron.NewScheduler(time.UTC)
	// if you are in China, you will need to change the time zone as follows
	// t.scheduler = gocron.NewScheduler(time.FixedZone("PRC", 8*60*60))

	//_, err := t.scheduler.Every("3s").Do(func()
	_, err := t.scheduler.CronWithSeconds("0/3 * * * * *").Do(func() {
		err := t.userTask.CheckUser(ctx)
		if err != nil {
			t.log.Error("CheckUser error", zap.Error(err))
		}
	})
	if err != nil {
		t.log.Error("CheckUser error", zap.Error(err))
	}

	t.scheduler.StartBlocking()
	return nil
}
func (t *TaskServer) Stop(ctx context.Context) error {
	t.scheduler.Stop()
	t.log.Info("TaskServer stop...")
	return nil
}
```

## File: internal/service/service.go
```go
package service

import (
	"github.com/go-nunu/nunu-layout-advanced/internal/repository"
	"github.com/go-nunu/nunu-layout-advanced/pkg/jwt"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"github.com/go-nunu/nunu-layout-advanced/pkg/sid"
)

type Service struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewService(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
	jwt *jwt.JWT,
) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
		jwt:    jwt,
		tm:     tm,
	}
}
```

## File: internal/service/user.go
```go
package service

import (
	"context"
	v1 "github.com/go-nunu/nunu-layout-advanced/api/v1"
	"github.com/go-nunu/nunu-layout-advanced/internal/model"
	"github.com/go-nunu/nunu-layout-advanced/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService interface {
	Register(ctx context.Context, req *v1.RegisterRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (string, error)
	GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error)
	UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error
}

func NewUserService(
	service *Service,
	userRepo repository.UserRepository,
) UserService {
	return &userService{
		userRepo: userRepo,
		Service:  service,
	}
}

type userService struct {
	userRepo repository.UserRepository
	*Service
}

func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	// check username
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return v1.ErrInternalServerError
	}
	if err == nil && user != nil {
		return v1.ErrEmailAlreadyUse
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Generate user ID
	userId, err := s.sid.GenString()
	if err != nil {
		return err
	}
	user = &model.User{
		UserId:   userId,
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	// Transaction demo
	err = s.tm.Transaction(ctx, func(ctx context.Context) error {
		// Create a user
		if err = s.userRepo.Create(ctx, user); err != nil {
			return err
		}
		// TODO: other repo
		return nil
	})
	return err
}

func (s *userService) Login(ctx context.Context, req *v1.LoginRequest) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return "", v1.ErrUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", err
	}
	token, err := s.jwt.GenToken(user.UserId, time.Now().Add(time.Hour*24*90))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &v1.GetProfileResponseData{
		UserId:   user.UserId,
		Nickname: user.Nickname,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	user.Email = req.Email
	user.Nickname = req.Nickname

	if err = s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
```

## File: internal/task/task.go
```go
package task

import (
	"github.com/go-nunu/nunu-layout-advanced/internal/repository"
	"github.com/go-nunu/nunu-layout-advanced/pkg/jwt"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"github.com/go-nunu/nunu-layout-advanced/pkg/sid"
)

type Task struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewTask(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
) *Task {
	return &Task{
		logger: logger,
		sid:    sid,
		tm:     tm,
	}
}
```

## File: internal/task/user.go
```go
package task

import (
	"context"
	"github.com/go-nunu/nunu-layout-advanced/internal/repository"
)

type UserTask interface {
	CheckUser(ctx context.Context) error
}

func NewUserTask(
	task *Task,
	userRepo repository.UserRepository,
) UserTask {
	return &userTask{
		userRepo: userRepo,
		Task:     task,
	}
}

type userTask struct {
	userRepo repository.UserRepository
	*Task
}

func (t userTask) CheckUser(ctx context.Context) error {
	// do something
	t.logger.Info("CheckUser")
	return nil
}
```

## File: pkg/app/app.go
```go
package app

import (
	"context"
	"github.com/go-nunu/nunu-layout-advanced/pkg/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	name    string
	servers []server.Server
}

type Option func(a *App)

func NewApp(opts ...Option) *App {
	a := &App{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func WithServer(servers ...server.Server) Option {
	return func(a *App) {
		a.servers = servers
	}
}

func WithName(name string) Option {
	return func(a *App) {
		a.name = name
	}
}

func (a *App) Run(ctx context.Context) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for _, srv := range a.servers {
		go func(srv server.Server) {
			err := srv.Start(ctx)
			if err != nil {
				log.Printf("Server start err: %v", err)
			}
		}(srv)
	}

	select {
	case <-signals:
		// Received termination signal
		log.Println("Received termination signal")
	case <-ctx.Done():
		// Context canceled
		log.Println("Context canceled")
	}

	// Gracefully stop the servers
	for _, srv := range a.servers {
		err := srv.Stop(ctx)
		if err != nil {
			log.Printf("Server stop err: %v", err)
		}
	}

	return nil
}
```

## File: pkg/config/config.go
```go
package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func NewConfig(p string) *viper.Viper {
	envConf := os.Getenv("APP_CONF")
	if envConf == "" {
		envConf = p
	}
	fmt.Println("load conf file:", envConf)
	return getConfig(envConf)
}

func getConfig(path string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigFile(path)
	err := conf.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return conf
}
```

## File: pkg/jwt/jwt.go
```go
package jwt

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type JWT struct {
	key []byte
}

type MyCustomClaims struct {
	UserId string
	jwt.RegisteredClaims
}

func NewJwt(conf *viper.Viper) *JWT {
	return &JWT{key: []byte(conf.GetString("security.jwt.key"))}
}

func (j *JWT) GenToken(userId string, expiresAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "",
			Subject:   "",
			ID:        "",
			Audience:  []string{},
		},
	})

	// Sign and get the complete encoded token as a string using the key
	tokenString, err := token.SignedString(j.key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWT) ParseToken(tokenString string) (*MyCustomClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	if strings.TrimSpace(tokenString) == "" {
		return nil, errors.New("token is empty")
	}
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.key, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
```

## File: pkg/log/log.go
```go
package log

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

const ctxLoggerKey = "zapLogger"

type Logger struct {
	*zap.Logger
}

func NewLog(conf *viper.Viper) *Logger {
	// log address "out.log" User-defined
	lp := conf.GetString("log.log_file_name")
	lv := conf.GetString("log.log_level")
	var level zapcore.Level
	//debug<info<warn<error<fatal<panic
	switch lv {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	hook := lumberjack.Logger{
		Filename:   lp,                             // Log file path
		MaxSize:    conf.GetInt("log.max_size"),    // Maximum size unit for each log file: M
		MaxBackups: conf.GetInt("log.max_backups"), // The maximum number of backups that can be saved for log files
		MaxAge:     conf.GetInt("log.max_age"),     // Maximum number of days the file can be saved
		Compress:   conf.GetBool("log.compress"),   // Compression or not
	}

	var encoder zapcore.Encoder
	if conf.GetString("log.encoding") == "console" {
		encoder = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "Logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
			EncodeTime:     timeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		})
	} else {
		encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
	}
	// default(both) log to console and file
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // Print to console and file
		level,
	)
	mode := conf.GetString("log.mode")
	switch mode {
	case "console":
		core = zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
	case "file":
		core = zapcore.NewCore(
			encoder,
			zapcore.AddSync(&hook),
			level,
		)
	}
	if conf.GetString("env") != "prod" {
		return &Logger{zap.New(core, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
	}
	return &Logger{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//enc.AppendString(t.Format("2006-01-02 15:04:05"))
	enc.AppendString(t.Format("2006-01-02 15:04:05.000000000"))
}

// WithValue Adds a field to the specified context
func (l *Logger) WithValue(ctx context.Context, fields ...zapcore.Field) context.Context {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
		c.Request = c.Request.WithContext(context.WithValue(ctx, ctxLoggerKey, l.WithContext(ctx).With(fields...)))
		return c
	}
	return context.WithValue(ctx, ctxLoggerKey, l.WithContext(ctx).With(fields...))
}

// WithContext Returns a zap instance from the specified context
func (l *Logger) WithContext(ctx context.Context) *Logger {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
	}
	zl := ctx.Value(ctxLoggerKey)
	ctxLogger, ok := zl.(*zap.Logger)
	if ok {
		return &Logger{ctxLogger}
	}
	return l
}
```

## File: pkg/server/grpc/grpc.go
```go
package grpc

import (
	"context"
	"fmt"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"google.golang.org/grpc"
	"net"
	"time"
)

type Server struct {
	*grpc.Server
	host   string
	port   int
	logger *log.Logger
}

type Option func(s *Server)

func NewServer(logger *log.Logger, opts ...Option) *Server {
	s := &Server{
		Server: grpc.NewServer(),
		logger: logger,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
func WithServerHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}
func WithServerPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		s.logger.Sugar().Fatalf("Failed to listen: %v", err)
	}
	if err = s.Server.Serve(lis); err != nil {
		s.logger.Sugar().Fatalf("Failed to serve: %v", err)
	}
	return nil

}
func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	s.Server.GracefulStop()

	s.logger.Info("Server exiting")

	return nil
}
```

## File: pkg/server/http/http.go
```go
package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"net/http"
	"time"
)

type Server struct {
	*gin.Engine
	httpSrv *http.Server
	host    string
	port    int
	logger  *log.Logger
}
type Option func(s *Server)

func NewServer(engine *gin.Engine, logger *log.Logger, opts ...Option) *Server {
	s := &Server{
		Engine: engine,
		logger: logger,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
func WithServerHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}
func WithServerPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.httpSrv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
		Handler: s,
	}

	if err := s.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Sugar().Fatalf("listen: %s\n", err)
	}

	return nil
}
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Sugar().Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.httpSrv.Shutdown(ctx); err != nil {
		s.logger.Sugar().Fatal("Server forced to shutdown: ", err)
	}

	s.logger.Sugar().Info("Server exiting")
	return nil
}
```

## File: pkg/server/server.go
```go
package server

import (
	"context"
	"net/url"
)

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

// Endpointer is registry endpoint.
type Endpointer interface {
	Endpoint() (*url.URL, error)
}
```

## File: pkg/sid/convert.go
```go
package sid

const (
	base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func IntToBase62(n int) string {
	if n == 0 {
		return string(base62[0])
	}

	var result []byte
	for n > 0 {
		result = append(result, base62[n%62])
		n /= 62
	}

	// 反转字符串
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}
```

## File: pkg/sid/sid.go
```go
package sid

import (
	"github.com/sony/sonyflake"
)

type Sid struct {
	sf *sonyflake.Sonyflake
}

func NewSid() *Sid {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	if sf == nil {
		panic("sonyflake not created")
	}
	return &Sid{sf}
}
func (s Sid) GenString() (string, error) {
	id, err := s.sf.NextID()
	if err != nil {
		return "", err
	}
	return IntToBase62(int(id)), nil
}
func (s Sid) GenUint64() (uint64, error) {
	return s.sf.NextID()
}
```

## File: pkg/zapgorm2/zapgorm2.go
```go
package zapgorm2

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
)

const ctxLoggerKey = "zapLogger"

type Logger struct {
	ZapLogger                 *zap.Logger
	SlowThreshold             time.Duration
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	LogLevel                  gormlogger.LogLevel
}

func New(zapLogger *zap.Logger) gormlogger.Interface {
	return &Logger{
		ZapLogger:                 zapLogger,
		LogLevel:                  gormlogger.Warn,
		SlowThreshold:             100 * time.Millisecond,
		Colorful:                  false,
		IgnoreRecordNotFoundError: false,
		ParameterizedQueries:      false,
	}
}

func (l *Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.logger(ctx).Sugar().Infof(msg, data...)
	}
}

// Warn print warn messages
func (l Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.logger(ctx).Sugar().Warnf(msg, data...)
	}
}

// Error print error messages
func (l Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.logger(ctx).Sugar().Errorf(msg, data...)
	}
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	elapsedStr := fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)
	logger := l.logger(ctx)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!errors.Is(err, gormlogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			logger.Error("trace", zap.Error(err), zap.String("elapsed", elapsedStr), zap.Int64("rows", rows), zap.String("sql", sql))
		} else {
			logger.Error("trace", zap.Error(err), zap.String("elapsed", elapsedStr), zap.Int64("rows", rows), zap.String("sql", sql))
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			logger.Warn("trace", zap.String("slow", slowLog), zap.String("elapsed", elapsedStr), zap.Int64("rows", rows), zap.String("sql", sql))
		} else {
			logger.Warn("trace", zap.String("slow", slowLog), zap.String("elapsed", elapsedStr), zap.Int64("rows", rows), zap.String("sql", sql))
		}
	case l.LogLevel == gormlogger.Info:
		sql, rows := fc()
		if rows == -1 {
			logger.Info("trace", zap.String("elapsed", elapsedStr), zap.Int64("rows", rows), zap.String("sql", sql))
		} else {
			logger.Info("trace", zap.String("elapsed", elapsedStr), zap.Int64("rows", rows), zap.String("sql", sql))
		}
	}
}

var (
	gormPackage = filepath.Join("gorm.io", "gorm")
)

func (l Logger) logger(ctx context.Context) *zap.Logger {
	logger := l.ZapLogger
	if ctx != nil {
		if c, ok := ctx.(*gin.Context); ok {
			ctx = c.Request.Context()
		}
		zl := ctx.Value(ctxLoggerKey)
		ctxLogger, ok := zl.(*zap.Logger)
		if ok {
			logger = ctxLogger
		}
	}

	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		default:
			return logger.WithOptions(zap.AddCallerSkip(i - 1))
		}
	}
	return logger
}
```

## File: scripts/README.md
```markdown
# `/scripts`

Scripts to perform various build, install, analysis, etc operations.

These scripts keep the root level Makefile small and simple.

Examples:

* https://github.com/kubernetes/helm/tree/master/scripts
* https://github.com/cockroachdb/cockroach/tree/master/scripts
* https://github.com/hashicorp/terraform/tree/master/scripts
```

## File: test/mocks/repository/repository.go
```go
// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTransaction is a mock of Transaction interface.
type MockTransaction struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionMockRecorder
}

// MockTransactionMockRecorder is the mock recorder for MockTransaction.
type MockTransactionMockRecorder struct {
	mock *MockTransaction
}

// NewMockTransaction creates a new mock instance.
func NewMockTransaction(ctrl *gomock.Controller) *MockTransaction {
	mock := &MockTransaction{ctrl: ctrl}
	mock.recorder = &MockTransactionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransaction) EXPECT() *MockTransactionMockRecorder {
	return m.recorder
}

// Transaction mocks base method.
func (m *MockTransaction) Transaction(ctx context.Context, fn func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transaction", ctx, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// Transaction indicates an expected call of Transaction.
func (mr *MockTransactionMockRecorder) Transaction(ctx, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transaction", reflect.TypeOf((*MockTransaction)(nil).Transaction), ctx, fn)
}
```

## File: test/mocks/repository/user.go
```go
// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/user.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	model "github.com/go-nunu/nunu-layout-advanced/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockUserRepositoryMockRecorder) Create(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserRepository)(nil).Create), ctx, user)
}

// GetByEmail mocks base method.
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", ctx, email)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *MockUserRepositoryMockRecorder) GetByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockUserRepository)(nil).GetByEmail), ctx, email)
}

// GetByID mocks base method.
func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockUserRepositoryMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockUserRepository)(nil).GetByID), ctx, id)
}

// Update mocks base method.
func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUserRepositoryMockRecorder) Update(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserRepository)(nil).Update), ctx, user)
}
```

## File: test/mocks/service/user.go
```go
// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/user.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	v1 "github.com/go-nunu/nunu-layout-advanced/api/v1"
	gomock "github.com/golang/mock/gomock"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// GetProfile mocks base method.
func (m *MockUserService) GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfile", ctx, userId)
	ret0, _ := ret[0].(*v1.GetProfileResponseData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfile indicates an expected call of GetProfile.
func (mr *MockUserServiceMockRecorder) GetProfile(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfile", reflect.TypeOf((*MockUserService)(nil).GetProfile), ctx, userId)
}

// Login mocks base method.
func (m *MockUserService) Login(ctx context.Context, req *v1.LoginRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, req)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserServiceMockRecorder) Login(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserService)(nil).Login), ctx, req)
}

// Register mocks base method.
func (m *MockUserService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockUserServiceMockRecorder) Register(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockUserService)(nil).Register), ctx, req)
}

// UpdateProfile mocks base method.
func (m *MockUserService) UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfile", ctx, userId, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProfile indicates an expected call of UpdateProfile.
func (mr *MockUserServiceMockRecorder) UpdateProfile(ctx, userId, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfile", reflect.TypeOf((*MockUserService)(nil).UpdateProfile), ctx, userId, req)
}
```

## File: test/server/handler/main_test.go
```go
package handler

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-nunu/nunu-layout-advanced/internal/handler"
	"github.com/go-nunu/nunu-layout-advanced/internal/middleware"
	"github.com/go-nunu/nunu-layout-advanced/pkg/config"
	jwt2 "github.com/go-nunu/nunu-layout-advanced/pkg/jwt"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	userId = "xxx"
)
var logger *log.Logger
var hdl *handler.Handler
var jwt *jwt2.JWT
var router *gin.Engine

func TestMain(m *testing.M) {
	fmt.Println("begin")
	err := os.Setenv("APP_CONF", "../../../config/local.yml")
	if err != nil {
		fmt.Println("Setenv error", err)
	}
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	// modify log directory
	logPath := filepath.Join("../../../", conf.GetString("log.log_file_name"))
	conf.Set("log.log_file_name", logPath)

	logger = log.NewLog(conf)
	hdl = handler.NewHandler(logger)

	jwt = jwt2.NewJwt(conf)
	gin.SetMode(gin.TestMode)
	router = gin.Default()
	router.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
		//middleware.SignMiddleware(log),
	)

	code := m.Run()
	fmt.Println("test end")

	os.Exit(code)
}

func performRequest(r http.Handler, method, path string, body *bytes.Buffer) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	return resp
}

func genToken(t *testing.T) string {
	token, err := jwt.GenToken(userId, time.Now().Add(time.Hour*24*90))
	if err != nil {
		t.Error(err)
		return token
	}
	return token
}

func newHttpExcept(t *testing.T, router *gin.Engine) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(router),
			Jar:       httpexpect.NewCookieJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			// httpexpect.NewDebugPrinter(t, true),
		},
	})
}
```

## File: test/server/handler/user_test.go
```go
package handler

import (
	v1 "github.com/go-nunu/nunu-layout-advanced/api/v1"
	"github.com/go-nunu/nunu-layout-advanced/internal/handler"
	"github.com/go-nunu/nunu-layout-advanced/internal/middleware"
	"github.com/go-nunu/nunu-layout-advanced/test/mocks/service"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUserHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.RegisterRequest{
		Password: "123456",
		Email:    "xxx@gmail.com",
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().Register(gomock.Any(), &params).Return(nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService)
	router.POST("/register", userHandler.Register)

	e := newHttpExcept(t, router)
	obj := e.POST("/register").
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
}

func TestUserHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.LoginRequest{
		Email:    "xxx@gmail.com",
		Password: "123456",
	}

	tk := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOiJ4eHgiLCJleHAiOjE3MzgyMjA1MTQsIm5iZiI6MTczMDQ0NDUxNCwiaWF0IjoxNzMwNDQ0NTE0fQ.3D4YupmPBCkv16ESnYyWSV5Mxcdu0twzEUqx0K-UiWo"
	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().Login(gomock.Any(), &params).Return(tk, nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService)
	router.POST("/login", userHandler.Login)

	obj := newHttpExcept(t, router).POST("/login").
		WithHeader("Content-Type", "application/json").
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	obj.Value("data").Object().Value("accessToken").IsEqual(tk)
}

func TestUserHandler_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nickname := "xxxxx"
	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().GetProfile(gomock.Any(), userId).Return(&v1.GetProfileResponseData{
		UserId:   userId,
		Nickname: nickname,
	}, nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService)
	router.Use(middleware.NoStrictAuth(jwt, logger))
	router.GET("/user", userHandler.GetProfile)

	obj := newHttpExcept(t, router).GET("/user").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
	objData := obj.Value("data").Object()
	objData.Value("userId").IsEqual(userId)
	objData.Value("nickname").IsEqual(nickname)
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := v1.UpdateProfileRequest{
		Nickname: "alan",
		Email:    "alan@gmail.com",
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().UpdateProfile(gomock.Any(), userId, &params).Return(nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.PUT("/user", userHandler.UpdateProfile)

	obj := newHttpExcept(t, router).PUT("/user").
		WithHeader("Content-Type", "application/json").
		WithHeader("Authorization", "Bearer "+genToken(t)).
		WithJSON(params).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()
	obj.Value("code").IsEqual(0)
	obj.Value("message").IsEqual("ok")
}
```

## File: test/server/repository/user_test.go
```go
package repository

import (
	"context"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-nunu/nunu-layout-advanced/internal/model"
	"github.com/go-nunu/nunu-layout-advanced/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var logger *log.Logger

func setupRepository(t *testing.T) (repository.UserRepository, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm connection: %v", err)
	}

	//rdb, _ := redismock.NewClientMock()

	repo := repository.NewRepository(logger, db)
	userRepo := repository.NewUserRepository(repo)

	return userRepo, mock
}

func TestUserRepository_Create(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	user := &model.User{
		Id:        1,
		UserId:    "123",
		Nickname:  "Test",
		Password:  "password",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WithArgs(user.UserId, user.Nickname, user.Password, user.Email, user.CreatedAt, user.UpdatedAt, user.DeletedAt, user.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_Update(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	user := &model.User{
		Id:        1,
		UserId:    "123",
		Nickname:  "Test",
		Password:  "password",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := userRepo.Update(ctx, user)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetById(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	userId := "123"

	rows := sqlmock.NewRows([]string{"id", "user_id", "username", "nickname", "password", "email", "created_at", "updated_at"}).
		AddRow(1, "123", "test", "Test", "password", "test@example.com", time.Now(), time.Now())
	mock.ExpectQuery("SELECT \\* FROM `users`").WillReturnRows(rows)

	user, err := userRepo.GetByID(ctx, userId)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "123", user.UserId)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByUsername(t *testing.T) {
	userRepo, mock := setupRepository(t)

	ctx := context.Background()
	email := "test@example.com"

	rows := sqlmock.NewRows([]string{"id", "user_id", "username", "nickname", "password", "email", "created_at", "updated_at"}).
		AddRow(1, "123", "test", "Test", "password", "test@example.com", time.Now(), time.Now())
	mock.ExpectQuery("SELECT \\* FROM `users`").WillReturnRows(rows)

	user, err := userRepo.GetByEmail(ctx, email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)

	assert.NoError(t, mock.ExpectationsWereMet())
}
```

## File: test/server/service/user_test.go
```go
package service_test

import (
	"context"
	"errors"
	"flag"
	"fmt"
	v1 "github.com/go-nunu/nunu-layout-advanced/api/v1"
	"github.com/go-nunu/nunu-layout-advanced/pkg/jwt"
	"github.com/go-nunu/nunu-layout-advanced/test/mocks/repository"
	"os"
	"testing"

	"github.com/go-nunu/nunu-layout-advanced/internal/model"
	"github.com/go-nunu/nunu-layout-advanced/internal/service"
	"github.com/go-nunu/nunu-layout-advanced/pkg/config"
	"github.com/go-nunu/nunu-layout-advanced/pkg/log"
	"github.com/go-nunu/nunu-layout-advanced/pkg/sid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	logger *log.Logger
	j      *jwt.JWT
	sf     *sid.Sid
)

func TestMain(m *testing.M) {
	fmt.Println("begin")

	err := os.Setenv("APP_CONF", "../../../config/local.yml")
	if err != nil {
		panic(err)
	}

	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	logger = log.NewLog(conf)
	j = jwt.NewJwt(conf)
	sf = sid.NewSid()

	code := m.Run()
	fmt.Println("test end")

	os.Exit(code)
}

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)

	userService := service.NewUserService(srv, mockUserRepo)

	ctx := context.Background()
	req := &v1.RegisterRequest{
		Password: "password",
		Email:    "test@example.com",
	}

	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, nil)
	mockTm.EXPECT().Transaction(ctx, gomock.Any()).Return(nil)

	err := userService.Register(ctx, req)

	assert.NoError(t, err)
}

func TestUserService_Register_UsernameExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo)

	ctx := context.Background()
	req := &v1.RegisterRequest{
		Password: "password",
		Email:    "test@example.com",
	}

	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(&model.User{}, nil)

	err := userService.Register(ctx, req)

	assert.Error(t, err)
}

func TestUserService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo)

	ctx := context.Background()
	req := &v1.LoginRequest{
		Email:    "xxx@gmail.com",
		Password: "password",
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		t.Error("failed to hash password")
	}

	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(&model.User{
		Password: string(hashedPassword),
	}, nil)

	token, err := userService.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo)

	ctx := context.Background()
	req := &v1.LoginRequest{
		Email:    "xxx@gmail.com",
		Password: "password",
	}

	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, errors.New("user not found"))

	_, err := userService.Login(ctx, req)

	assert.Error(t, err)
}

func TestUserService_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo)

	ctx := context.Background()
	userId := "123"

	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(&model.User{
		UserId: userId,
		Email:  "test@example.com",
	}, nil)

	user, err := userService.GetProfile(ctx, userId)

	assert.NoError(t, err)
	assert.Equal(t, userId, user.UserId)
}

func TestUserService_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo)

	ctx := context.Background()
	userId := "123"
	req := &v1.UpdateProfileRequest{
		Nickname: "testuser",
		Email:    "test@example.com",
	}

	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(&model.User{
		UserId: userId,
		Email:  "old@example.com",
	}, nil)
	mockUserRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

	err := userService.UpdateProfile(ctx, userId, req)

	assert.NoError(t, err)
}

func TestUserService_UpdateProfile_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo)

	ctx := context.Background()
	userId := "123"
	req := &v1.UpdateProfileRequest{
		Nickname: "testuser",
		Email:    "test@example.com",
	}

	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(nil, errors.New("user not found"))

	err := userService.UpdateProfile(ctx, userId, req)

	assert.Error(t, err)
}
```

## File: web/index.html
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>

</body>
</html>
```

## File: .dockerignore
```
.git
.idea
```

## File: .gitattributes
```
* text=auto

# Force the following filetypes to have unix eols, so Windows does not break them
*.* text eol=lf

# Windows forced line-endings
/.idea/* text eol=crlf

#
## These files are binary and should be left untouched
#

# (binary is a macro for -text -diff)
*.png binary
```

## File: .gitignore
```
storage/logs
.idea
*.log
deploy/docker-compose/conf
deploy/docker-compose/data
```

## File: go.mod
```
module github.com/go-nunu/nunu-layout-advanced

go 1.23.10

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/duke-git/lancet/v2 v2.3.6
	github.com/gavv/httpexpect/v2 v2.17.0
	github.com/gin-gonic/gin v1.10.1
	github.com/glebarez/sqlite v1.11.0
	github.com/go-co-op/gocron v1.37.0
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/golang/mock v1.6.0
	github.com/google/wire v0.6.0
	github.com/redis/go-redis/v9 v9.10.0
	github.com/sony/sonyflake v1.2.1
	github.com/spf13/viper v1.20.1
	github.com/stretchr/testify v1.10.0
	github.com/swaggo/files v1.0.1
	github.com/swaggo/gin-swagger v1.6.0
	github.com/swaggo/swag v1.16.4
	go.mongodb.org/mongo-driver v1.17.4
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.39.0
	google.golang.org/grpc v1.73.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gorm.io/driver/mysql v1.6.0
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.30.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/TylerBrock/colorjson v0.0.0-20200706003622-8a50f05110d2 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/bytedance/sonic v1.11.6 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/glebarez/go-sqlite v1.21.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hpcloud/tail v1.0.0 // indirect
	github.com/imkira/go-interpol v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.6.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sanity-io/litter v1.5.5 // indirect
	github.com/sergi/go-diff v1.0.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.40.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	github.com/yudai/gojsondiff v1.0.0 // indirect
	github.com/yudai/golcs v0.0.0-20170316035057-ecda9a501e82 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/exp v0.0.0-20221208152030-732eee02a75a // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/libc v1.22.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.5.0 // indirect
	modernc.org/sqlite v1.23.1 // indirect
	moul.io/http2curl/v2 v2.3.0 // indirect
)
```

## File: LICENSE
```
MIT License

Copyright (c) 2023 Nunu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## File: Makefile
```
.PHONY: init
init:
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: bootstrap
bootstrap:
	cd ./deploy/docker-compose && docker compose up -d && cd ../../
	go run ./cmd/migration
	nunu run ./cmd/server

.PHONY: mock
mock:
	mockgen -source=internal/service/user.go -destination test/mocks/service/user.go
	mockgen -source=internal/repository/user.go -destination test/mocks/repository/user.go
	mockgen -source=internal/repository/repository.go -destination test/mocks/repository/repository.go

.PHONY: test
test:
	go test -coverpkg=./internal/handler,./internal/service,./internal/repository -coverprofile=./coverage.out ./test/server/...
	go tool cover -html=./coverage.out -o coverage.html

.PHONY: build
build:
	go build -ldflags="-s -w" -o ./bin/server ./cmd/server

.PHONY: docker
docker:
	docker build -f deploy/build/Dockerfile --build-arg APP_RELATIVE_PATH=./cmd/task -t 1.1.1.1:5000/demo-task:v1 .
	docker run --rm -i 1.1.1.1:5000/demo-task:v1

.PHONY: swag
swag:
	swag init  -g cmd/server/main.go -o ./docs
```

## File: README_zh.md
```markdown
# Nunu — A CLI tool for building go aplication.


Nunu是一个基于Golang的应用脚手架，它的名字来自于英雄联盟中的游戏角色，一个骑在雪怪肩膀上的小男孩。和努努一样，该项目也是站在巨人的肩膀上，它是由Golang生态中各种非常流行的库整合而成的，它们的组合可以帮助你快速构建一个高效、可靠的应用程序。

[英文介绍](https://github.com/go-nunu/nunu/blob/main/README.md)

![Nunu](https://github.com/go-nunu/nunu/blob/main/.github/assets/banner.png)

## 文档
* [使用指南](https://github.com/go-nunu/nunu/blob/main/docs/zh/guide.md)
* [分层架构](https://github.com/go-nunu/nunu/blob/main/docs/zh/architecture.md)
* [上手教程](https://github.com/go-nunu/nunu/blob/main/docs/zh/tutorial.md)
* [高效编写单元测试](https://github.com/go-nunu/nunu/blob/main/docs/zh/unit_testing.md)

## 许可证

Nunu是根据MIT许可证发布的。有关更多信息，请参见[LICENSE](LICENSE)文件。
```

## File: README.md
```markdown
# Nunu — A CLI tool for building Go applications.

Nunu is a scaffolding tool for building Go applications. Its name comes from a game character in League of Legends, a little boy riding on the shoulders of a Yeti. Just like Nunu, this project stands on the shoulders of giants, as it is built upon a combination of popular libraries from the Go ecosystem. This combination allows you to quickly build efficient and reliable applications.

[简体中文介绍](https://github.com/go-nunu/nunu/blob/main/README_zh.md)

![Nunu](https://github.com/go-nunu/nunu/blob/main/.github/assets/banner.png)

## Documentation
* [User Guide](https://github.com/go-nunu/nunu/blob/main/docs/en/guide.md)
* [Architecture](https://github.com/go-nunu/nunu/blob/main/docs/en/architecture.md)
* [Getting Started Tutorial](https://github.com/go-nunu/nunu/blob/main/docs/en/tutorial.md)
* [Unit Testing](https://github.com/go-nunu/nunu/blob/main/docs/en/unit_testing.md)


## License

Nunu is released under the MIT License. For more information, see the [LICENSE](LICENSE) file.
```
