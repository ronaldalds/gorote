package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/ronaldalds/gorote-core/core"
)

type Middleware struct {
	App *fiber.App
	// RedisStore *databases.RedisStore
}

type LogTelemetry struct {
	Timestamp    string              `json:"timestamp"`
	Method       string              `json:"method"`
	Path         string              `json:"path"`
	Headers      map[string][]string `json:"headers"`
	IP           string              `json:"ip"`
	RequestBody  map[string]any      `json:"request_body"`
	Status       int                 `json:"status"`
	Latency      int64               `json:"latency"`
	ResponseBody string              `json:"response_body"`
}

func NewMiddleware(app *fiber.App) *Middleware {
	return &Middleware{
		App: app,
		// RedisStore: databases.DB.RedisStore,
	}
}

func (m *Middleware) CorsMiddleware() {
	m.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false,
		MaxAge:           300,
	}))
}

func (m *Middleware) SecurityMiddleware() {
	m.App.Use(helmet.New(helmet.Config{
		XSSProtection: "1; mode=block",
	}))
}

func sendLogToLoki(logData LogTelemetry) {
	lokiURL := fmt.Sprintf("%v:%v/loki/api/v1/push", Env.LogsUrl, Env.LogsPort)
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())

	// Converte os dados do log para JSON
	jsonLog, err := json.Marshal(logData)
	if err != nil {
		log.Println("Erro ao converter log para JSON:", err)
		return
	}
	fmt.Println(string(jsonLog))

	// Monta a estrutura esperada pelo Loki
	params := core.HttpRequestParams{
		Method: core.POST,
		URL:    lokiURL,
		Headers: core.Headers{
			ContentType: "application/json",
		},
		Body: map[string]any{
			"streams": []map[string]any{
				{
					"stream": map[string]any{
						"app": Env.AppName,
					},
					"values": [][]string{
						{timestamp, string(jsonLog)},
					},
				},
			},
		},
	}
	res, err := core.SendHttpRequest(params)
	if err != nil {
		log.Println("Erro ao enviar log para o Loki:", err)
		return
	}
	defer res.Body.Close()

	log.Println("Log enviado para o Loki com sucesso!")
}

func (m *Middleware) Telemetry(ConfidentialPath ...string) {
	m.App.Use(func(ctx *fiber.Ctx) error {
		start := time.Now()
		// request
		var body map[string]any
		var logData LogTelemetry
		ctx.BodyParser(&body)

		logData.Timestamp = time.Now().Format(time.RFC3339)
		logData.Method = ctx.Method()
		logData.Path = ctx.Path()
		logData.Headers = ctx.GetReqHeaders()
		logData.IP = ctx.IP()
		logData.RequestBody = body

		for _, path := range ConfidentialPath {
			if strings.Contains(ctx.Path(), path) {
				logData.RequestBody = map[string]any{"confidential": true}
				break
			}
		}
		// end request
		if err := ctx.Next(); err != nil {
			var e *fiber.Error
			if errors.As(err, &e) {
				// response error
				logData.Status = e.Code
				logData.Latency = time.Since(start).Milliseconds()
				logData.ResponseBody = e.Message

				// Imprime o log formatado para ser capturado pelo Loki
				sendLogToLoki(logData)
				// end response error
			}
			return err
		}
		// response
		logData.Status = ctx.Response().StatusCode()
		logData.Latency = time.Since(start).Milliseconds()
		logData.ResponseBody = string(ctx.Response().Body())

		// Imprime o log formatado para ser capturado pelo Loki
		sendLogToLoki(logData)
		//end response
		return nil
	})
}
