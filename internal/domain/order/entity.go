package order

import "github.com/alan-b-lima/almodon/pkg/uuid"

var Canon = uuid.Max

func NextVersion(version int) int {
	return max(version+1, 1)
}
