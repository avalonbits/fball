package db

import (
	"context"
	"database/sql"

	"git.cana.pw/avalonbits/fball"
)

type Inserter struct {
	DB *sql.DB
}

var insertSQL = `
INSERT INTO RequestCache(Endpoint, Params, Timestamp, Response)
				  VALUES(?, ?, ?, ?);`

func (i *Inserter) Timezone(ctx context.Context, tr fball.TimezoneResponse) error {
	return transact(ctx, i.DB, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, insertSQL)
		if err != nil {
			return err
		}
		defer stmt.Close()

		res, err := stmt.ExecContext(ctx, fball.EP_Timezone, "", tr.Timestamp, &tr)
		if err != nil {
			return err
		}
		if _, err := res.RowsAffected(); err != nil {
			return err
		}
		return nil
	})
}
