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

// PageOffset - converts a page number to the page offset.
func PageOffset(page uint32, size uint32) int32 {
	if page == 0 || page == 1 {
		return 0
	}
	offset := (page-1)*size + 1
	return int32(offset)
}
