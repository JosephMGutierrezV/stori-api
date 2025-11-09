package out

import "context"

type ObjectStorage interface {
	PutObject(
		ctx context.Context,
		bucket string,
		key string,
		contentType string,
		data []byte,
	) error
}
