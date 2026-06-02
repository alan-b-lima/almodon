package lotitem

import (
	"time"

	"github.com/alan-b-lima/almodon/pkg/money"
)

func ProcessAmount(amount int64) (int64, error) {
	if amount < 0 {
		return 0, ErrAmountNegative
	}
	return amount, nil
}

func ProcessUnitCost(cost money.Money) (money.Money, error) {
	if cost < 0 {
		return 0, ErrUnitCostNegative
	}
	return cost, nil
}

func ProcessExpires(expires time.Time) (time.Time, error) {
	return expires, nil
}
