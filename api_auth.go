package qvapay

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p *LoginRequest) ToReader() *strings.Reader {
	byte, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}
	return strings.NewReader(string(byte))
}

func (c *apiClient) Login(ctx context.Context, payload LoginRequest) (APIResult, error) {

	url := fmt.Sprintf("%s/%s", c.server, login)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, payload.ToReader())
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	apiCallDebugger(req, c.debug)

	req.WithContext(ctx)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %v", err)
	}
	defer DrainBody(res.Body)
	if c.debug != nil {
		c.dumpResponse(res)
	}

	var result APIResult
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to create HTTP response: %v", err)
	}
	return result, nil
}
