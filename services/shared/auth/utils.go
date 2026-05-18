package auth

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgtype"
)

func StringToUUID(rawUUID string) (pgtype.UUID, error) {
	var uuid pgtype.UUID

	if err := uuid.Scan(rawUUID); err != nil {
		slog.Error("UUID parsing failed", "error", err)
		return pgtype.UUID{}, err
	}

	return uuid, nil
}
