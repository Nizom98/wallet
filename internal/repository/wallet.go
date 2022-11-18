package repository

type wallet struct {
	id string `json:"id"`
	name string `json:"name"`
	balance float64 `json:"balance"`
	status bool `json:"status"`
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

