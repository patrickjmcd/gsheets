package gsheets

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Client struct {
	Service *sheets.Service
}

func New(ctx context.Context, credentialsFilePath, b64ServiceAccount *string) (*Client, error) {

	var srv *sheets.Service

	if credentialsFilePath != nil && b64ServiceAccount != nil {
		log.Error().Msg("both credentialsFilePath and b64ServiceAccount cannot be provided")
		return nil, fmt.Errorf("both credentialsFilePath and b64ServiceAccount cannot be provided")
	}

	if credentialsFilePath == nil && b64ServiceAccount == nil {
		log.Error().Msg("either credentialsFilePath or b64ServiceAccount must be provided")
		return nil, fmt.Errorf("either credentialsFilePath or b64ServiceAccount must be provided")
	}

	if credentialsFilePath != nil {
		s, err := sheets.NewService(ctx, option.WithCredentialsFile(*credentialsFilePath))
		if err != nil {
			log.Error().Err(err).Msg("unable to retrieve Sheets client")
			return nil, fmt.Errorf("unable to retrieve Sheets client: %w", err)
		}
		srv = s
	}

	if b64ServiceAccount != nil {
		sa, err := tokenFromB64(*b64ServiceAccount)
		if err != nil {
			log.Error().Err(err).Msg("unable to retrieve Service Account")
			return nil, fmt.Errorf("unable to retrieve Service Account: %w", err)
		}
		s, err := sheets.NewService(ctx, option.WithCredentialsJSON(sa))
		if err != nil {
			log.Error().Err(err).Msg("unable to retrieve Sheets client")
			return nil, fmt.Errorf("unable to retrieve Sheets client: %w", err)
		}
		srv = s
	}
	client := &Client{
		Service: srv,
	}
	return client, nil
}
