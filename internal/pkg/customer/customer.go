package customer

type Customer string

const (
	CustomerInternal Customer = "internal"
)

func (c Customer) String() string {
	return string(c)
}
