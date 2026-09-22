package uuid

import (
	std_uuid "uuid"
)

func New() string {
	return std_uuid.New().String()
}
