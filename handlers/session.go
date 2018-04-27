// Package handlers предназначен для перехвата http запросов и их обслуживания.
package handlers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/qa-dev/jsonwire-grid-wda-agent/command"
	"github.com/qa-dev/jsonwire-grid-wda-agent/config"
	"github.com/qa-dev/jsonwire-grid-wda-agent/wda"
	"log"
	"net/http"
	"unicode/utf8"
)

var (
	acl         string = s3.BucketCannedACLPublicRead
	contentType string = "video/quicktime"
	forcePath   bool   = true
)

// SessionHandler содержит в себе конфиг приложения.
type SessionHandler struct {
	cfg   *config.Config
	proxy *wda.Proxy
}

// NewSessionHandler описывает конструктор со структурой SessionHandler.
func NewSessionHandler(cfg *config.Config, pr *wda.Proxy) *SessionHandler {
	return &SessionHandler{
		cfg:   cfg,
		proxy: pr,
	}
}

// ServeHTTP перехватывает http запрос /session/*
func (h *SessionHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete && h.cfg.Video.Enable {
		go func() {
			sessionID := r.URL.Path[utf8.RuneCount([]byte("/session/")):]

			videoFile, err := command.FinishVideo()
			if err != nil {
				log.Printf("Error while finishing video: %v\n", err)
				return
			}

			if videoFile != nil {
				creds := credentials.NewStaticCredentials(h.cfg.Video.S3.AccessKey, h.cfg.Video.S3.SecretKey, "")
				sess := session.Must(session.NewSession(&aws.Config{
					Credentials:      creds,
					Endpoint:         h.cfg.Video.S3.Endpoint,
					Region:           &h.cfg.Video.S3.Region,
					S3ForcePathStyle: &forcePath,
				}))
				svc := s3.New(sess)

				key := sessionID + ".mov"

				_, err = svc.PutObject(&s3.PutObjectInput{
					Bucket:      &h.cfg.Video.S3.Bucket,
					Body:        videoFile,
					Key:         &key,
					ContentType: &contentType,
					ACL:         &acl,
				})
				if err != nil {
					log.Printf("Error while sending video: %v\n", err)
				}
			}
		}()

	}

	h.proxy.ServeHTTP(rw, r)
}
