package uuid

import uid "github.com/google/uuid"

func New() string {
	return uid.NewString()
}
