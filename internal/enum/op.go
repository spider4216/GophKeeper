package enum

type OpType string

const (
	OpCreate OpType = "CREATE"
	OpUpdate OpType = "UPDATE"
	OpDelete OpType = "DELETE"
)

func (t OpType) String() string {
	return string(t)
}
