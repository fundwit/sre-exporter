package oss

import (
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

//go:generate mockgen -source oss/oss_client.go -destination oss/oss_client_mock.go -package oss
type ObjectLister interface {
	ListObjects(bucket string, options ...oss.Option) (*oss.ListObjectsResult, error)
}

type OssClient struct {
	Endpoint  string
	AccessKey string
	Secret    string
}

func NewOssClient() *OssClient {
	return &OssClient{
		Endpoint:  os.ExpandEnv(os.Getenv("OSS_ENDPOINT")),
		AccessKey: os.ExpandEnv(os.Getenv("OSS_ACCESSKEY")),
		Secret:    os.ExpandEnv(os.Getenv("OSS_SECRET")),
	}
}

func (c *OssClient) ListObjects(bucketName string, options ...oss.Option) (*oss.ListObjectsResult, error) {
	// CRC check is enabled by default
	cli, err := oss.New(c.Endpoint, c.AccessKey, c.Secret)
	if err != nil {
		return nil, err
	}

	bucket, err := cli.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	r, err := bucket.ListObjects(options...)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
