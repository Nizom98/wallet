package models

type WalletRepository interface {
	Create(name string, balance float64, status bool) Walleter
	ByID(id string) (Walleter, error)
	All() []Walleter
	Transaction(fn func(repo WalletRepository) error) error
	UpdateByID(id string, name *string, balance *float64, status *bool) error
}
