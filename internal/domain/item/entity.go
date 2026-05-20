package item

import (
	"strconv"
	"time"

	"github.com/alan-b-lima/almodon/pkg/money"
)

const (
	ExpiresWarnThreshold = 30 * 24 * time.Hour
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

func StatusExpires(expires time.Time) Expiration {
	if expires.IsZero() {
		return ExpirationNone
	}

	diff := time.Until(expires)
	if diff <= 0 {
		return ExpirationExpired
	}
	if diff < ExpiresWarnThreshold {
		return ExpirationWarning
	}
	return ExpirationFine
}

type Expiration int

const (
	ExpirationNone Expiration = iota
	ExpirationFine
	ExpirationWarning
	ExpirationExpired
)

var expirations = [...]string{
	ExpirationNone:    "none",
	ExpirationFine:    "fine",
	ExpirationWarning: "warning",
	ExpirationExpired: "expired",
}

func (s Expiration) String() string {
	if 0 <= int(s) && int(s) < len(expirations) {
		str := expirations[s]
		if str != "" {
			return str
		}
	}

	return "expiration(" + strconv.Itoa(int(s)) + ")"
}

func (s Expiration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}
