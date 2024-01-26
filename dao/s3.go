package dao

import (
	"PInfo-server/config"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"PInfo-server/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (d *Dao) MakeRawUrl(ctx context.Context, disableSSL bool, domain, bucket, key string) (string, error) {
	schema := "https"
	if disableSSL {
		schema = "http"
	}
	urlStr := fmt.Sprintf("%s://%s/%s/%s", schema, domain, bucket, key)
	return urlStr, nil
}

// GetPresignUrl .
func (d *Dao) GetPresignUrl(ctx context.Context, bucket, key string, expireHour time.Duration) (string, error) {
	log.DebugContextf(ctx, "Create Presign client, key:%s", key)

	if !config.AppConfig().CosInfo.SignFlag {
		return d.MakeRawUrl(ctx, config.AppConfig().CosInfo.DisableSSL,
			config.AppConfig().CosInfo.Domain, bucket, key)
	}
	svc := s3.New(d.sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(expireHour * time.Hour)

	if err != nil {
		log.ErrorContextf(ctx, "Create Presigned URL fail:%v", err)
	} else {
		log.InfoContextf(ctx, "Presigned URL For object: %s, expire in: %dh", urlStr, expireHour)
	}

	return urlStr, err
}

// UploadFile .
func (d *Dao) UploadFile(ctx context.Context, bucket, key string, reader io.Reader, ftype, md5Str string) error {

	uploader := s3manager.NewUploader(d.sess)
	ui := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   reader,
	}
	if len(ftype) > 0 {
		ui.ContentType = aws.String(ftype)
	}
	if len(md5Str) > 0 {
		ui.ContentMD5 = aws.String(md5Str)
	}
	_, err := uploader.Upload(ui)
	if err != nil {
		log.ErrorContextf(ctx, "Unable to upload %q to %q, %v", key, bucket, err)
		return err
	}
	log.InfoContextf(ctx, "Successfully uploaded %q to %q", key, bucket)
	return err
}

// Download 从cos下载文件，返回字节流
func (d *Dao) Download(ctx context.Context, bucket, key string, timeout int) ([]byte, error) {

	svc := s3.New(d.sess)
	objOutput := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond) // 单位毫秒
	defer cancel()

	result, err := svc.GetObjectWithContext(ctx, objOutput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.ErrorContextf(ctx, "GetObject err, bucket:%s, key:%s, err:%+v", bucket, key, err)
			}
		}
		return nil, err
	} else {
		log.InfoContextf(ctx, "GetObject ok, bucket:%s, key:%s", bucket, key)
		return io.ReadAll(result.Body)
	}
}

// DownloadFile 从cos下载文件，存储到本地
func (d *Dao) DownloadFile(ctx context.Context, bucket, key string) error {

	file, err := os.Create(key)
	if err != nil {
		log.ErrorContextf(ctx, "Unable to open file %q, %v\n", key, err)
		return err
	}
	defer file.Close()
	downloader := s3manager.NewDownloader(d.sess)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		log.ErrorContextf(ctx, "Unable to download key %q, %v", key, err)
		return err
	}
	log.InfoContextf(ctx, "Downloaded %s %d bytes", file.Name(), numBytes)
	return err
}

// ParseUrlKey 从url解析path路径剥离bucket
func (d *Dao) ParseUrlKey(ctx context.Context, urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		log.ErrorContextf(ctx, "parse url failed, url: %s", urlStr)
		return "", errors.New("parse url failed")
	}

	return u.Path, nil
}

// RawDownload 直接从url下载字节流
func (d *Dao) RawDownload(ctx context.Context, url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		defer func() {
			log.InfoContextf(ctx, "done download data")
			resp.Body.Close()
		}()
	} else {
		log.ErrorContextf(ctx, "resp err: %v", err)
		return nil, err
	}

	log.InfoContextf(ctx, "start download data from url:%s", url)
	return io.ReadAll(resp.Body)
}
