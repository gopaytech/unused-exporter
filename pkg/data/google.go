package data

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/gopaytech/unused-exporter/pkg/model"
	"github.com/gopaytech/unused-exporter/pkg/settings"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

var (
	ErrMissingToken = errors.New("google cloud dns token is missing")
)

type GoogleData struct {
	computeService *compute.Service
	projects       []string
}

func (g *GoogleData) GetUnusedIP() ([]model.IPAddress, error) {
	var IPs []model.IPAddress

	for _, project := range g.projects {
		aggregatedList, err := g.computeService.Addresses.AggregatedList(project).Do()
		if err != nil {
			return nil, err
		}

		for region, addressesScopedList := range aggregatedList.Items {
			for _, address := range addressesScopedList.Addresses {
				if address.Status == "RESERVED" {
					IPs = append(IPs, model.IPAddress{
						Cloud:    "GCP",
						Region:   region,
						Identity: project + "/" + address.Name,
						Value:    address.Address,
						Type:     address.AddressType,
						Used:     false,
					})
				}
			}
		}
	}

	return IPs, nil
}

func (g *GoogleData) GetUsedIP() ([]model.IPAddress, error) {
	var IPs []model.IPAddress

	for _, project := range g.projects {
		aggregatedList, err := g.computeService.Addresses.AggregatedList(project).Do()
		if err != nil {
			return nil, err
		}

		for region, addressesScopedList := range aggregatedList.Items {
			for _, address := range addressesScopedList.Addresses {
				if address.Status != "RESERVED" {
					IPs = append(IPs, model.IPAddress{
						Cloud:    "GCP",
						Region:   region,
						Identity: project + "/" + address.Name,
						Value:    address.Address,
						Type:     address.AddressType,
						Used:     true,
					})
				}
			}
		}
	}

	return IPs, nil
}

func NewGoogleData(settings settings.Settings) (*GoogleData, error) {
	g := &GoogleData{
		projects: settings.GCPProjects,
	}

	creds, err := newWithAPITokenSource(settings.GCPServiceAccount, []string{compute.ComputeScope})
	if err != nil {
		return nil, err
	}

	computeService, err := compute.NewService(context.Background(), option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}

	g.computeService = computeService

	return g, nil
}

func newWithAPITokenSource(token string, scopes []string) (*google.Credentials, error) {
	ctx := context.Background()

	if token == "" {
		return nil, ErrMissingToken
	}

	data, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	creds, err := google.CredentialsFromJSON(ctx, data, scopes...)
	if err != nil {
		return nil, err
	}

	return creds, nil
}
