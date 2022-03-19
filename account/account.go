package account

import "errors"

// Account struct
type Account struct {
	owner   string
	balance int
}

var errNotEnoughBalance = errors.New("not enough balance")

// NewAccount create Account
func NewAccount(owner string) *Account {
	account := Account{owner: owner, balance: 0}
	return &account
}

// Deposit x amount on the account
func (account *Account) Deposit(amount int) {
	account.balance += amount
}

// Withdraw x amount on the account
func (account *Account) Withdraw(amount int) error {
	if account.balance < amount {
		return errNotEnoughBalance
	}
	account.balance -= amount

	return nil
}

// ChangeOwner of the account
func (account *Account) ChangeOwner(newOwner string) {
	account.owner = newOwner
}

// Owner of the account
func (account Account) Owner() string {
	return account.owner
}

// Balance of the account
func (account Account) Balance() int {
	return account.balance
}
