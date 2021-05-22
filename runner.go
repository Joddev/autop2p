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
)

func (err *InvestError) Error() string {
	switch err.Code {
	case Duplicated:
		return "duplicated investment"
	case InsufficientCapacity:
		return "insufficient residual capacity"
	default:
		return "unsupported InvestError Code"
	}
}
