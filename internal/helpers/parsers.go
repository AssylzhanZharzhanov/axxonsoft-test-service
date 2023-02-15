package helpers

import (
	"fmt"
	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/internal/domain"
	"strconv"
)

// ExtractInt64Route - returns the route int64 variable for the current request.
func ExtractInt64Route(vars map[string]string, key string) (int64, error) {
	str, ok := vars[key]
	if !ok {
		return 0, fmt.Errorf("bad routing")
	}
	result, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return int64(result), nil
}

// ParsePageOrGetDefault - parses the page number, if present, otherwise it will return a default value.
func ParsePageOrGetDefault(pageStr string) uint32 {
	if len(pageStr) == 0 {
		return domain.DefaultPageNumber
	}
	page, _ := strconv.Atoi(pageStr)
	return uint32(page)
}

// ParseSizeOrGetDefault - parses the page size, if present, otherwise it will return a default value.
func ParseSizeOrGetDefault(pageStr string) uint32 {
	if len(pageStr) == 0 {
		return domain.DefaultPageNumber
	}
	page, _ := strconv.Atoi(pageStr)
	return uint32(page)
}
