package autop2p

type Runner interface {
	ListProducts() []Product
	InvestProduct(product *Product, amount int) *InvestError
}

type InvestError struct {
	Code string
}

const (
	Duplicated           = "Duplicated"
	InsufficientCapacity = "InsufficientCapacity"
	InsufficientBalance  = "InsufficientBalance"
)

func (err *InvestError) Error() string {
	switch err.Code {
	case Duplicated:
		return "duplicated investment"
	case InsufficientCapacity:
		return "insufficient residual capacity"
	case InsufficientBalance:
		return "Insufficient balance"
	default:
		return "unsupported InvestError Code"
	}
}
