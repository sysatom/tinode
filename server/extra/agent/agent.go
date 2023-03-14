package agent

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
	"os"
)

type Data struct {
	Id      string      `json:"id"`
	Content interface{} `json:"content"`
}

func StartInfo() {
	log.Println("Agent started")
}

func URI() string {
	if len(os.Args) < 2 {
		panic("args error")
	}
	return os.Args[1]
}

func PostData(agentURI string, data Data) error {
	c := resty.New()
	resp, err := c.R().
		SetContext(context.Background()).
		SetBody(data).
		Post(agentURI)
	if err != nil {
		return err
	}
	if resp.StatusCode() == http.StatusOK {
		return nil
	}
	return fmt.Errorf("%d", resp.StatusCode())
}
