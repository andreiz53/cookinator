package util

import "github.com/google/uuid"

func NullUUID(arg uuid.UUID) *uuid.UUID {
	var id *uuid.UUID

	if arg != uuid.Nil {
		id = &arg
	}

	return id
}
