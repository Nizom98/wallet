package models

type Walleter interface {
	ID() string
	Name() string
	Balance() float64
	Status() bool
}

type WalletManager interface {
	Create(name string) (Walleter, error)
	ByID(id string) (Walleter, error)
	List() []Walleter
	IncreaseBalanceBy(id string, amount float64) error
	DecreaseBalanceBy(id string, amount float64) error
	TransferBalance(fromID, toID string, amount float64) error
	DeactivateByID(id string) error
	UpdateName(id, name string) error
}


