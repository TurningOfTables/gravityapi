package main

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GravityResponse struct {
	Data   interface{}       `json:"data"`
	Meta   map[string]string `json:"meta"`
	Errors []GravityError    `json:"errors"`
}

type GravityError struct {
	Id     string `json:"id"`     // A unique Id for the instance of the error. Is added automatically and does not need to be provided.
	Status string `json:"status"` // The HTTP status code applicable to the problem
	Code   string `json:"code"`   // An application specific error code
	Title  string `json:"title"`  // A short summary of the problem that is the same for each occurrence of the problem
	Detail string `json:"detail"` // A longer explanation specific to this occurrence of the problem
}

func (ge *GravityError) Error() string {
	return ge.Detail
}

func SendGravityResponse(c fiber.Ctx, gr *GravityResponse) error {
	httpStatus := fiber.StatusOK
	gr.Meta = make(map[string]string)
	gr.Meta["timestamp"] = time.Now().Format(time.RFC3339)

	if gr.Data == nil {
		gr.Data = []interface{}{}
	}

	if len(gr.Errors) == 0 {
		gr.Errors = []GravityError{}
	} else {
		statusInt, err := strconv.Atoi(gr.Errors[0].Status)
		if err == nil {
			httpStatus = statusInt
		}
		for k := range gr.Errors {
			gr.Errors[k].Id = uuid.New().String()
		}
	}

	return c.Status(httpStatus).JSON(gr)
}
