package repository

type wallet struct {
	id      string
	name    string
	balance float64
	status  bool
}

func (wal *wallet) ID() string {
	return wal.id
}

func (wal *wallet) Name() string {
	return wal.name
}

func (wal *wallet) Balance() float64 {
	return wal.balance
}

func (wal *wallet) Status() bool {
	return wal.status
}
