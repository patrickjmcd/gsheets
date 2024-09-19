package gsheets

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var valueInputOptions = map[string]bool{
	"RAW":          true,
	"USER_ENTERED": true,
}

type AuthConfig struct {
	credentialsFilePath *string
	b64ServiceAccount   *string
}

type Client[T any] struct {
	Service          *sheets.Service
	authConfig       *AuthConfig
	spreadsheetId    string
	parseRowFn       func(context.Context, []interface{}) (T, error)
	formatRowFn      func(context.Context, T) []interface{}
	valueInputOption string
}

type ClientOption[T any] func(*Client[T])

func WithParseRowFn[T any](fn func(context.Context, []interface{}) (T, error)) ClientOption[T] {
	return func(c *Client[T]) {
		c.parseRowFn = fn
	}
}

func WithFormatRowFn[T any](fn func(context.Context, T) []interface{}) ClientOption[T] {
	return func(c *Client[T]) {
		c.formatRowFn = fn
	}
}

func WithCredentialsFilePath[T any](credentialsFilePath string) ClientOption[T] {
	return func(c *Client[T]) {
		if c.authConfig == nil {
			c.authConfig = &AuthConfig{}
		}
		c.authConfig.credentialsFilePath = &credentialsFilePath
	}
}

func WithB64ServiceAccount[T any](b64ServiceAccount string) ClientOption[T] {
	return func(c *Client[T]) {
		if c.authConfig == nil {
			c.authConfig = &AuthConfig{}
		}
		c.authConfig.b64ServiceAccount = &b64ServiceAccount
	}
}

func WithValueInputOption[T any](valueInputOption string) ClientOption[T] {
	return func(c *Client[T]) {
		c.valueInputOption = valueInputOption
	}
}

func New[T any](ctx context.Context, spreadsheetId string, opts ...ClientOption[T]) (*Client[T], error) {
	if spreadsheetId == "" {
		log.Error().Msg("spreadsheetId must be provided")
		return nil, fmt.Errorf("spreadsheetId must be provided")
	}

	client := &Client[T]{
		spreadsheetId:    spreadsheetId,
		valueInputOption: "USER_ENTERED",
	}
	for _, opt := range opts {
		opt(client)
	}

	if client.authConfig.credentialsFilePath != nil && client.authConfig.b64ServiceAccount != nil {
		log.Error().Msg("both credentialsFilePath and b64ServiceAccount cannot be provided")
		return nil, fmt.Errorf("both credentialsFilePath and b64ServiceAccount cannot be provided")
	}

	if client.authConfig.credentialsFilePath == nil && client.authConfig.b64ServiceAccount == nil {
		log.Error().Msg("either credentialsFilePath or b64ServiceAccount must be provided")
		return nil, fmt.Errorf("either credentialsFilePath or b64ServiceAccount must be provided")
	}

	if _, ok := valueInputOptions[client.valueInputOption]; !ok {
		log.Error().Str("valueInputOption", client.valueInputOption).Msg("invalid valueInputOption")
		return nil, fmt.Errorf("invalid valueInputOption")
	}

	if client.authConfig.credentialsFilePath != nil {
		s, err := sheets.NewService(ctx, option.WithCredentialsFile(*client.authConfig.credentialsFilePath))
		if err != nil {
			log.Error().Err(err).Msg("unable to retrieve Sheets client")
			return nil, fmt.Errorf("unable to retrieve Sheets client: %w", err)
		}
		client.Service = s
	}

	if client.authConfig.b64ServiceAccount != nil {
		sa, err := tokenFromB64(*client.authConfig.b64ServiceAccount)
		if err != nil {
			log.Error().Err(err).Msg("unable to retrieve Service Account")
			return nil, fmt.Errorf("unable to retrieve Service Account: %w", err)
		}
		s, err := sheets.NewService(ctx, option.WithCredentialsJSON(sa))
		if err != nil {
			log.Error().Err(err).Msg("unable to retrieve Sheets client")
			return nil, fmt.Errorf("unable to retrieve Sheets client: %w", err)
		}
		client.Service = s
	}

	return client, nil
}
