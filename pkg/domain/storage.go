package domain

import "io"

type StorageObject struct {
	Bucket    string
	Path      string
	Blob      io.Reader
	SizeBytes uint64
}
