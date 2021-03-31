package query

import (
	"context"
	"database/sql"
	"fmt"

	"git.cana.pw/avalonbits/fball"
)

type Handler struct {
	db     *sql.DB
	cfball *fball.Client
}

func New(db *sql.DB, cfball *fball.Client) *Handler {
	return &Handler{
		db:     db,
		cfball: cfball,
	}
}

var noParamQuery = `
SELECT Response, Timestamp from RequestCache
	WHERE
		Endpoint = ?
	ORDER BY
		Timestamp DESC;
`

func (h *Handler) Timezone(ctx context.Context) (fball.TimezoneResponse, error) {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return fball.TimezoneResponse{}, fmt.Errorf("error creating transaction: %w", err)
	}
	defer tx.Rollback()

	return fball.TimezoneResponse{}, nil
}
