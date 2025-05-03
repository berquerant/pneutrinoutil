package infra

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
	"github.com/berquerant/pneutrinoutil/pkg/alog"
)

type StorageParam struct {
	UseS3   bool
	RootDir string
	Debug   bool // enable debug logging
}

func NewStorage(ctx context.Context, param *StorageParam) (Object, error) {
	if param.UseS3 {
		c, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return nil, err
		}

		return NewS3(s3.NewFromConfig(c, func(opt *s3.Options) {
			if param.Debug {
				// https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/configure-logging.html
				opt.ClientLogMode = aws.LogRequest | aws.LogResponse
				opt.Logger = logging.LoggerFunc(func(classification logging.Classification, format string, v ...any) {
					msg := fmt.Sprintf(format, v...)
					switch classification {
					case logging.Warn:
						alog.L().Warn(msg)
					case logging.Debug:
						alog.L().Debug(msg)
					default:
						alog.L().Info(msg)
					}
				})
			}

			if x := os.Getenv("AWS_ENDPOINT_URL"); x != "" {
				opt.BaseEndpoint = aws.String(x)
				opt.EndpointOptions.DisableHTTPS = os.Getenv("AWS_S3_DISABLE_HTTPS") == "true"
			}
		})), nil
	}

	return NewFileSystem(param.RootDir), nil
}
