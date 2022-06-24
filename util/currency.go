package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupportedCurrenct(currenct string) bool {
	switch currenct {
	case USD, EUR, CAD:
		return true
	}
	return false
}
