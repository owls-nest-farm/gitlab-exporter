package main

import (
	"time"

	"github.com/xanzy/go-gitlab"
)

type LabelService struct {
	*BaseService
}

type Label struct {
	Type      string     `json:"type"`
	URL       string     `json:"url"`
	Name      string     `json:"body"`
	Color     string     `json:"assignee"`
	CreatedAt *time.Time `json:"created_at"`
}

func NewLabelService(e *Exporter) *LabelService {
	return &LabelService{
		BaseService: &BaseService{
			exporter: e,
			filename: "labels.json",
		},
	}
}

func (l *LabelService) GetAll() ([]Label, error) {
	labels, _, err := l.exporter.Client.Labels.ListLabels(l.exporter.CurrentProject.ID, &gitlab.ListLabelsOptions{})
	if err != nil {
		panic(err)
	}
	lbls := make([]Label, len(labels))
	for i, label := range labels {
		now := time.Now()
		lbls[i] = Label{
			Type:      "label",
			URL:       l.exporter.CurrentProject.WebURL,
			Name:      label.Name,
			Color:     label.Color,
			CreatedAt: &now,
		}
	}

	return lbls, nil
}

func (l *LabelService) Export() {
	labels, err := l.GetAll()
	if err != nil {
		panic(err)
	}
	if len(labels) > 0 {
		l.exporter.State.Labels = append(l.exporter.State.Labels, labels...)
	}
}
