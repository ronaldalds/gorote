package utils

import (
	"fmt"
)

func Pagination[T any](page, limit uint, data *[]T) error {
	count := uint(len(*data))
	start := (page - 1) * limit
	end := page * limit

	if start >= count {
		return fmt.Errorf("page not exist")
	}

	if end > count {
		end = count
	}
	*data = (*data)[start:end]
	return nil
}
