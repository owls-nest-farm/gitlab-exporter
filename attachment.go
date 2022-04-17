package main

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/xanzy/go-gitlab"
)

type AttachmentService struct {
	*BaseService
}

type Attachment struct {
	Type             string     `json:"type"`
	URL              string     `json:"url"`
	DownloadURL      string     `json:"download_url"`
	User             string     `json:"user"`
	AssetName        string     `json:"asset_name"`
	AssetContentType string     `json:"asset_content_type"`
	AssetURL         string     `json:"asset_url"`
	CreatedAt        *time.Time `json:"created_at"`
}

func NewAttachmentService(e *Exporter) *AttachmentService {
	return &AttachmentService{
		BaseService: &BaseService{
			exporter: e,
			filename: "attachments.json",
		},
	}
}

func (a *AttachmentService) Get(matches []string, service interface{}) *Attachment {
	project := a.exporter.CurrentProject
	filename := matches[0]
	path := matches[1]

	switch srv := service.(type) {
	case *gitlab.Issue:
		return &Attachment{
			Type:             "attachment",
			URL:              fmt.Sprintf("%s%s", project.WebURL, path),
			DownloadURL:      fmt.Sprintf("%s%s", project.WebURL, path),
			User:             srv.Author.WebURL,
			AssetName:        filename,
			AssetContentType: "",
			AssetURL:         "",
			CreatedAt:        srv.CreatedAt,
		}
	case *gitlab.MergeRequest:
		return &Attachment{
			Type:             "attachment",
			URL:              fmt.Sprintf("%s%s", project.WebURL, path),
			DownloadURL:      fmt.Sprintf("%s%s", project.WebURL, path),
			User:             srv.Author.WebURL,
			AssetName:        filename,
			AssetContentType: "",
			AssetURL:         "",
			CreatedAt:        srv.CreatedAt,
		}
	}

	return nil
}

func (a *AttachmentService) Export(field string, service interface{}) {
	r := reflect.ValueOf(service)
	f := reflect.Indirect(r).FieldByName(field)

	if matches := a.exporter.AttachmentRegex.FindStringSubmatch(f.String()); matches != nil {
		path := matches[2]
		err := a.exporter.DownloadFile(path)
		// https://gitlab.com/gl-group1/foo/uploads/1e690825814e23a29bd8810f567829e7/invite.go
		// https://gitlab.com/gl-group1/foo/uploads/448016f60dba8d6580eeeea894dc1d09/user.go`
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v\n", *a.Get(matches[1:], service))
//		a.exporter.State.Attachments = append(a.exporter.State.Attachments, *a.Get(matches[1:], service))
	}
}
