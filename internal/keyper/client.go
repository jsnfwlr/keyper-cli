package keyper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jsnfwlr/keyper-cli/internal/app"
)

type Config struct {
	Address  string `env:"KEYPER_ADDRESS" validate:"required"`
	Username string `env:"KEYPER_USERNAME" validate:"required"`
	Password string `env:"KEYPER_PASSWORD" validate:"required"`
}

type Client struct {
	config Config
	token  string
	http   *http.Client
}

func NewClient() (client *Client, fault error) {
	http := &http.Client{}

	c := &Client{
		http: http,
	}

	err := c.Config()
	if err != nil {
		return nil, err
	}

	err = c.Auth()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) Do(method, route string, auth bool, payload any, response any) (fault error) {
	req, err := http.NewRequest(method, c.config.Address+route, nil)
	if err != nil {
		return err
	}

	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		br := bytes.NewReader(b)

		req, err = http.NewRequest(method, c.config.Address+route, br)
		if err != nil {
			return err
		}

		req.Header.Add("Content-Type", "application/json")
		// req.Header.Set("Accept", "application/json")
	}

	if auth {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read %[4]s response body for %[1]s (status code: %[2]d): %[3]w", c.config.Address+route, resp.StatusCode, err, method)
	}

	if resp.StatusCode > app.StatusError {
		return fmt.Errorf("problem with %[4]s request (%[5]v) for %[1]s (status code: %[2]d): %[3]s ", c.config.Address+route, resp.StatusCode, string(out), method, req.Header)
	}

	err = json.Unmarshal(out, response)
	if err != nil {
		return fmt.Errorf("could not unmarshal response: %w", err)
	}

	return nil
}

func (c *Client) Config() (fault error) {
	config := Config{}
	err := app.Config(&config)
	if err != nil {
		return err
	}

	c.config = config

	return nil
}
