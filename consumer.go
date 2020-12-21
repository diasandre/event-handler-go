package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var client = resty.New()

func defaultConsumer(event Event, c *fiber.Ctx) error {
	resp, err := client.R().SetBody(event.Payload).Post("https://example:8080/api")
	return responseHandler(event, c, resp, err)
}

func responseHandler(event Event, c *fiber.Ctx, resp *resty.Response, err error) error {
	if err != nil {
		zapLogger.Error(err.Error())
		if resp != nil && resp.StatusCode() != 0 {
			c.Status(resp.StatusCode())
		} else {
			c.Status(fiber.StatusBadRequest)
		}
		return c.Send([]byte(err.Error()))
	}

	logSuccessResponse(event, resp)
	return c.Send(resp.Body())
}

func logSuccessResponse(event Event, resp *resty.Response) {
	zapLogger.Info(
		successResponseLog,
		zap.String("eventType", string(event.Type)),
		zap.String("payload", event.Payload),
		zap.Int("statusCode", resp.StatusCode()),
		zap.ByteString("body", resp.Body()))
}