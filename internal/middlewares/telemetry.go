package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ronaldalds/base-go-api/internal/config/envs"
	"github.com/ronaldalds/base-go-api/internal/utils"

	"github.com/gofiber/fiber/v2"
)

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

func sendLogToLoki(logData LogTelemetry) {
	lokiURL := fmt.Sprintf("%v:%v/loki/api/v1/push", envs.Env.LogsUrl, envs.Env.LogsPort)
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())

	// Converte os dados do log para JSON
	jsonLog, err := json.Marshal(logData)
	if err != nil {
		log.Println("Erro ao converter log para JSON:", err)
		return
	}
	fmt.Println(string(jsonLog))

	// Monta a estrutura esperada pelo Loki
	params := utils.HttpRequestParams{
		Method: utils.POST,
		URL:    lokiURL,
		Headers: utils.Headers{
			ContentType: "application/json",
		},
		Body: map[string]any{
			"streams": []map[string]any{
				{
					"stream": map[string]any{
						"app": envs.Env.AppName,
					},
					"values": [][]string{
						{timestamp, string(jsonLog)},
					},
				},
			},
		},
	}
	res, err := utils.SendHttpRequest(params)
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
