package transformer

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDString(id pgtype.UUID) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", id.Bytes[0:4], id.Bytes[4:6], id.Bytes[6:8], id.Bytes[8:10], id.Bytes[10:16])
}
