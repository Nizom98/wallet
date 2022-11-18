package repository

import (
	"example.com/wallet/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	repo := NewRepo()
	expect := &wallet{
		name:    "test_name",
		balance: 9999,
		status:  true,
	}
	got := repo.Create(expect.name, expect.balance, expect.status)

	assert.NotNil(t, got)
	assert.True(t, expect.Name() == got.Name())
	assert.True(t, expect.Balance() == got.Balance())
	assert.True(t, expect.Status() == got.Status())
}

func TestUpdateByID(t *testing.T) {
	repo := NewRepo()
	name, balance, status := "test_name", float64(9999), true
	oldWal := repo.Create(name, balance, status)

	expectName := name+"postfix"

	err := repo.UpdateByID(oldWal.ID(), utils.PtrString(expectName), nil, nil)
	assert.Nil(t, err)
	if err != nil {
		return
	}

	updWal, err := repo.ByID(oldWal.ID())
	assert.Nil(t, err)
	assert.True(t, updWal.Name() == expectName)
	assert.True(t, balance == updWal.Balance())
	assert.True(t, status == updWal.Status())
}

func TestByID_found(t *testing.T) {
	repo := NewRepo()

	repo.Create("test_name", 9999, true)
	expect := repo.Create("test_name_2", 8888, true)

	got, err := repo.ByID(expect.ID())
	assert.Nil(t, err)
	if err != nil {
		return
	}

	assert.True(t, expect.ID() == got.ID())
	assert.True(t, expect.Name() == got.Name())
	assert.True(t, expect.Balance() == got.Balance())
	assert.True(t, expect.Status() == got.Status())
}

func TestByID_notFound(t *testing.T) {
	repo := NewRepo()
	expect := repo.Create("test_name_2", 8888, true)

	got, err := repo.ByID(expect.ID())
	assert.Nil(t, err)
	if err != nil {
		return
	}

	assert.True(t, expect.ID() == got.ID())
	assert.True(t, expect.Name() == got.Name())
	assert.True(t, expect.Balance() == got.Balance())
	assert.True(t, expect.Status() == got.Status())
}
