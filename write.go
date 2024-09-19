package gsheets

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/sheets/v4"
)

func (c *Client[T]) AppendToSheet(ctx context.Context, sheetName string, entries []T) error {
	nextRow, err := c.FindNextRow(ctx, sheetName)
	if err != nil {
		return err
	}
	writeRange := fmt.Sprintf("A%d", nextRow)
	var vr sheets.ValueRange
	for _, m := range entries {
		vr.Values = append(vr.Values, c.formatRowFn(ctx, m))
	}
	_, err = c.Service.Spreadsheets.Values.Update(c.spreadsheetId, writeRange, &vr).ValueInputOption(c.valueInputOption).Do()
	if err != nil {
		log.Error().Err(err).Msg("unable to write data to sheet")
		return fmt.Errorf("unable to write data to sheet: %v", err)
	}
	return nil
}
