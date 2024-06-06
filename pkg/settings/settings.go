package settings

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Settings represents the configuration settings for the application.
type Settings struct {
	EnableGCP bool `envconfig:"ENABLE_GCP" default:"false"`
	EnableAWS bool `envconfig:"ENABLE_AWS" default:"false"`
	// strict
	GCPServiceAccount     string   `envconfig:"GCP_SERVICE_ACCOUNT"`
	GCPProjects           []string `envconfig:"GCP_PROJECTS"`
	AWSRegion             string   `envconfig:"AWS_REGION"`
	AWSIAMSecret          string   `envconfig:"AWS_IAM_SECRET"`
	AWSIAMKey             string   `envconfig:"AWS_IAM_KEY"`
	AWSAssumeRolesRegions []string `envconfig:"AWS_ASSUME_ROLES_REGIONS"`
	AWSAssumeRoleDuration int      `envconfig:"AWS_ASSUME_ROLE_DURATION" default:"120"`
	Port                  int      `envconfig:"PORT" default:"8080"`
}

func NewSettings() Settings {
	var settings Settings

	err := envconfig.Process("", &settings)
	if err != nil {
		log.Fatalln(err)
	}

	return settings
}
