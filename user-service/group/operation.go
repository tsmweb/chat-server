package group

const (
	CREATE = iota
	UPDATE
)

// Operation type of operation such as CREATE and UPDATE.
type Operation int

func (o Operation) String() string {
	return [...]string{"CREATE", "UPDATE"}[o]
}
