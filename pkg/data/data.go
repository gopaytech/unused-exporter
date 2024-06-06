package data

import "github.com/gopaytech/unused-exporter/pkg/model"

type Data interface {
	GetUnusedIP() ([]model.IPAddress, error)
}
