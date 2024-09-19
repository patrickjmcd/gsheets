package gsheets

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
)

func (c *Client[T]) ReadFromSheet(ctx context.Context, sheetName string, readRange string) ([]T, error) {
	readSheetRange := fmt.Sprintf("%s!%s", sheetName, readRange) // e.g. "Sheet1!A1:B"
	resp, err := c.Service.Spreadsheets.Values.Get(c.spreadsheetId, readSheetRange).Do()
	if err != nil {
		log.Error().Err(err).Str("sheetName", sheetName).Str("readRange", readRange).Str("readSheetRange", readSheetRange).Msg("unable to retrieve data from sheet")
		return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}
	var parsedRows []T

	for _, row := range resp.Values {
		log.Trace().Interface("row", row).Msg("row")
		if row == nil || len(row) == 0 {
			log.Warn().Msgf("empty row in sheet %s\n", sheetName)
			continue
		}

		parsed, err := c.parseRowFn(ctx, row)
		if err != nil {
			log.Error().Err(err).Msg("unable to parse row")
			continue
		}
		parsedRows = append(parsedRows, parsed)
	}
	return parsedRows, nil
}
