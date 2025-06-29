package store

import "github.com/jackc/pgx/v5/pgtype"

func toPgInt8(v *int64) pgtype.Int8 {
	if v == nil {
		return pgtype.Int8{Valid: false}
	}
	return pgtype.Int8{Int64: *v, Valid: true}
}

func toPgBool(v *bool) pgtype.Bool {
	if v == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *v, Valid: true}
}

func toPgText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func fromPgText(p pgtype.Text) *string {
	if !p.Valid {
		return nil
	}
	return &p.String
}

func fromPgBool(p pgtype.Bool) bool {
	if !p.Valid {
		return false
	}

	return p.Bool
}
