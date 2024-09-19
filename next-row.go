package gsheets

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
)

func (c *Client[T]) FindNextRow(ctx context.Context, sheetName string) (int, error) {
	startRow := 1
	readRange := fmt.Sprintf("%s!A%d:A", sheetName, startRow)
	resp, err := c.Service.Spreadsheets.Values.Get(c.spreadsheetId, readRange).Do()
	if err != nil {
		log.Error().Err(err).Msg("unable to retrieve data from sheet")
		return 0, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}
	if len(resp.Values) == 0 {
		return startRow, nil
	}
	return startRow + len(resp.Values), nil
}
