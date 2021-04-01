package db

import (
	"database/sql"

	"git.cana.pw/avalonbits/fball"
)

type Inserter struct {
	DB *sql.DB
}

var inserNoParam = `
INSERT INTO RequestCache(Endpoint, Params, Timestamp, Response)
				  VALUES(?, '', ?, ?);`

func (i *Inserter) Timezone(tr fball.TimezoneResponse) error {
	return nil
}
