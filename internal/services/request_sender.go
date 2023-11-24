package services

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type RequestSender struct {
}

func NewRequestSender() *RequestSender {
	return &RequestSender{}
}

func (r *RequestSender) Send(url string, reqNum int) (int, error) {

	client := resty.New()

	jsonBody := fmt.Sprintf("{ 'iteration': %d }", reqNum)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonBody).
		Post(url)

	if err != nil {
		return 0, err
	}

	return resp.StatusCode(), nil
}
