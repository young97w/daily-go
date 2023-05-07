package v1

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDB_DoTx(t *testing.T) {
	db := memoryDB(t)
	fn := func(ctx context.Context, tx *Tx) error {
		return nil
	}
	err := db.DoTx(context.Background(), fn, &sql.TxOptions{})
	require.NoError(t, err)
}
