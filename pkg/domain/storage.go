package domain

import "io"

type StorageObject struct {
	Bucket    string
	Path      string
	Blob      io.ReadSeeker
	SizeBytes uint64
}
