package helpers

import (
	"fmt"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"

	"github.com/caarlos0/env"
)

// LoadConfig - reads configuration from the environment variables.
func LoadConfig() (config domain.AppConfig, _ error) {
	if err := env.Parse(&config); err != nil {
		return config, fmt.Errorf("failed to read configuration: %v", err)
	}
	return config, nil
}
