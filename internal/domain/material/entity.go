package material

import "strconv"

const (
	NameMaxLen = 128

	DescriptionMaxLen = 4096
)

func ProcessName(name string) (string, error) {
	if name == "" {
		return "", ErrNameEmpty
	}

	if len(name) >= NameMaxLen {
		return "", ErrNameEmpty
	}

	return name, nil
}

func ProcessECampus(ecampus int) (int, error) {
	return ecampus, nil
}

func ProcessCATMAT(catmat int) (int, error) {
	return catmat, nil
}

func ProcessSIADS(siads int) (int, error) {
	return siads, nil
}

func ProcessDescription(description string) (string, error) {
	if len(description) >= DescriptionMaxLen {
		return "", ErrDescriptionTooLong
	}

	return description, nil
}

func ProcessUnit(unit string) (string, error) {
	return unit, nil
}

func StatusAvailable(amount, min float64) Stock {
	switch {
	case min <= 0:
		return StockFine
	case amount < min:
		return StockWarning
	}

	return StockFine
}

func ProcessMin(quantity float64) (float64, error) {
	if quantity < 0 {
		return 0, ErrMinNegative
	}

	return quantity, nil
}

type Stock int

const (
	StockEmpty Stock = iota
	StockWarning
	StockFine
)

var stocks = [...]string{
	StockFine:    "fine",
	StockWarning: "warning",
	StockEmpty:   "empty",
}

func (s Stock) String() string {
	if 0 <= int(s) && int(s) < len(stocks) {
		str := stocks[s]
		if str != "" {
			return str
		}
	}

	return "stock(" + strconv.Itoa(int(s)) + ")"
}

func (s Stock) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}
