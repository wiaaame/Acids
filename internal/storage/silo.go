package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Лимиты файлов услуг: изображение до 5 МБ, короткое видео до 100 МБ.
const (
	MaxImageSize = 5 << 20
	MaxVideoSize = 100 << 20
)

var (
	imageExts = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	videoExts = map[string]bool{".mp4": true, ".webm": true, ".mov": true}
)

// Silo — обёртка над объектным хранилищем Silo (MinIO), совместимым с S3 API.
type Silo struct {
	client  *minio.Client
	bucket  string
	baseURL string
}

// NewSilo создаёт клиент, корзину (bucket) и политику публичного чтения.
// parameter baseURL используется для построения прямых ссылок на файлы.
func NewSilo(endpoint, accessKey, secretKey, bucket, baseURL string, useSSL bool) (*Silo, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio client: %w", err)
	}

	ctx := context.Background()
	ok, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("bucket exists: %w", err)
	}
	if !ok {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("make bucket: %w", err)
		}
	}

	// Публичный доступ на чтение, чтобы прямые ссылки /acids/images/... работали без токенов.
	if err := client.SetBucketPolicy(ctx, bucket, policyReadOnly(bucket)); err != nil {
		return nil, fmt.Errorf("set bucket policy: %w", err)
	}

	return &Silo{
		client:  client,
		bucket:  bucket,
		baseURL: strings.TrimRight(baseURL, "/") + "/",
	}, nil
}

// policyReadOnly — JSON-политика MinIO: чтение объектов всем (public read).
func policyReadOnly(bucket string) string {
	return fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, bucket)
}

// UploadImage загружает изображение услуги, возвращает ключ (латиница).
func (s *Silo) UploadImage(fh *multipart.FileHeader) (string, error) {
	return s.upload(fh, "images", imageExts, MaxImageSize)
}

// UploadVideo загружает короткое видео услуги, возвращает ключ (латиница).
func (s *Silo) UploadVideo(fh *multipart.FileHeader) (string, error) {
	return s.upload(fh, "videos", videoExts, MaxVideoSize)
}

// URL возвращает прямую публичную ссылку на объект по его ключу.
func (s *Silo) URL(key string) string {
	if key == "" {
		return ""
	}
	return s.baseURL + key
}

// upload — общая логика загрузки: проверка размера/расширения,
// генерация латинского имени и запись объекта в корзину.
func (s *Silo) upload(fh *multipart.FileHeader, folder string, allowed map[string]bool, maxSize int64) (string, error) {
	if fh == nil {
		return "", errors.New("файл не передан")
	}
	if fh.Size > maxSize {
		return "", fmt.Errorf("файл слишком большой: %d байт (лимит %d)", fh.Size, maxSize)
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !allowed[ext] {
		return "", fmt.Errorf("недопустимый тип файла %q", ext)
	}

	src, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Определяем Content-Type по содержимому (первые 512 байт), с фолбэком по расширению.
	sniff := make([]byte, 512)
	n, _ := src.Read(sniff)
	contentType := http.DetectContentType(sniff[:n])
	if contentType == "application/octet-stream" {
		if ct := mime.TypeByExtension(ext); ct != "" {
			contentType = ct
		}
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	// Имя генерируется на латинице: <unix>_<hex>.jpg (папка images/ или videos/ — отдельно).
	key := fmt.Sprintf("%d_%s%s", time.Now().Unix(), randomHex(6), ext)
	objectKey := folder + "/" + key

	_, err = s.client.PutObject(context.Background(), s.bucket, objectKey, src, fh.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("put object: %w", err)
	}
	return objectKey, nil
}

// randomHex возвращает n случайных байт в hex-виде (для уникальных латинских имён).
func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "000000"
	}
	return hex.EncodeToString(b)
}