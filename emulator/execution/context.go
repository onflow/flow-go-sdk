package execution

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/dapperlabs/cadence"
	"github.com/dapperlabs/cadence/runtime"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/emulator/types"
	"github.com/dapperlabs/flow-go-sdk/keys"
)

type CheckerFunc func([]byte, runtime.Location) error

// RuntimeContext implements host functionality required by the Cadence runtime.
//
// A context is short-lived and is intended to be used when executing a single transaction.
//
// The logic in this runtime context is specific to the emulator and is designed to be
// used with a Blockchain instance.
type RuntimeContext struct {
	ledger          *types.LedgerView
	signingAccounts []runtime.Address
	checker         CheckerFunc
	logs            []string
	events          []runtime.Event
}

// NewRuntimeContext returns a new RuntimeContext instance.
func NewRuntimeContext(ledger *types.LedgerView) *RuntimeContext {
	return &RuntimeContext{
		ledger:  ledger,
		checker: func([]byte, runtime.Location) error { return nil },
		events:  make([]runtime.Event, 0),
	}
}

// SetSigningAccounts sets the signing accounts for this context.
//
// Signing accounts are the accounts that signed the transaction executing
// inside this context.
func (r *RuntimeContext) SetSigningAccounts(accounts []flow.Address) {
	signingAccounts := make([]runtime.Address, len(accounts))

	for i, account := range accounts {
		signingAccounts[i] = runtime.Address(cadence.NewAddress(account))
	}

	r.signingAccounts = signingAccounts
}

// GetSigningAccounts gets the signing accounts for this context.
//
// Signing accounts are the accounts that signed the transaction executing
// inside this context.
func (r *RuntimeContext) GetSigningAccounts() []runtime.Address {
	return r.signingAccounts
}

// SetChecker sets the semantic checker function for this context.
func (r *RuntimeContext) SetChecker(checker CheckerFunc) {
	r.checker = checker
}

// Events returns all events emitted by the runtime to this context.
func (r *RuntimeContext) Events() []runtime.Event {
	return r.events
}

// Logs returns all logs emitted by the runtime to this context.
func (r *RuntimeContext) Logs() []string {
	return r.logs
}

// GetValue gets a register value from the world state.
func (r *RuntimeContext) GetValue(owner, controller, key []byte) ([]byte, error) {
	v, _ := r.ledger.Get(fullKey(string(owner), string(controller), string(key)))
	return v, nil
}

// SetValue sets a register value in the world state.
func (r *RuntimeContext) SetValue(owner, controller, key, value []byte) error {
	r.ledger.Set(fullKey(string(owner), string(controller), string(key)), value)
	return nil
}

// CreateAccount creates a new account and inserts it into the world state.
//
// This function returns an error if the input is invalid.
//
// After creating the account, this function calls the onAccountCreated callback registered
// with this context.
func (r *RuntimeContext) CreateAccount(publicKeys [][]byte) (runtime.Address, error) {
	latestAccountID, _ := r.ledger.Get(keyLatestAccount)

	accountIDInt := big.NewInt(0).SetBytes(latestAccountID)
	accountIDBytes := accountIDInt.Add(accountIDInt, big.NewInt(1)).Bytes()

	accountAddress := flow.BytesToAddress(accountIDBytes)

	accountID := accountAddress[:]

	r.ledger.Set(fullKey(string(accountID), "", keyBalance), big.NewInt(0).Bytes())

	err := r.setAccountPublicKeys(accountID, publicKeys)
	if err != nil {
		return runtime.Address{}, err
	}

	r.ledger.Set(keyLatestAccount, accountID)

	r.Log("Creating new account\n")
	r.Log(fmt.Sprintf("Address: %x", accountAddress))

	return runtime.Address(cadence.NewAddress(accountAddress)), nil
}

// AddAccountKey adds a public key to an existing account.
//
// This function returns an error if the specified account does not exist or
// if the key insertion fails.
func (r *RuntimeContext) AddAccountKey(address runtime.Address, publicKey []byte) error {
	accountID := address[:]

	bal, err := r.ledger.Get(fullKey(string(accountID), "", keyBalance))
	if err != nil {
		return err
	}
	if bal == nil {
		return fmt.Errorf("account with ID %s does not exist", accountID)
	}

	publicKeys, err := r.getAccountPublicKeys(accountID)
	if err != nil {
		return err
	}

	publicKeys = append(publicKeys, publicKey)

	return r.setAccountPublicKeys(accountID, publicKeys)
}

// RemoveAccountKey removes a public key by index from an existing account.
//
// This function returns an error if the specified account does not exist, the
// provided key is invalid, or if key deletion fails.
func (r *RuntimeContext) RemoveAccountKey(address runtime.Address, index int) (publicKey []byte, err error) {
	accountID := address[:]

	bal, err := r.ledger.Get(fullKey(string(accountID), "", keyBalance))
	if err != nil {
		return nil, err
	}
	if bal == nil {
		return nil, fmt.Errorf("account with ID %s does not exist", accountID)
	}

	publicKeys, err := r.getAccountPublicKeys(accountID)
	if err != nil {
		return publicKey, err
	}

	if index < 0 || index > len(publicKeys)-1 {
		return publicKey, fmt.Errorf("invalid key index %d, account has %d keys", index, len(publicKeys))
	}

	removedKey := publicKeys[index]

	publicKeys = append(publicKeys[:index], publicKeys[index+1:]...)

	err = r.setAccountPublicKeys(accountID, publicKeys)
	if err != nil {
		return publicKey, err
	}

	return removedKey, nil
}

func (r *RuntimeContext) getAccountPublicKeys(accountID []byte) (publicKeys [][]byte, err error) {
	countBytes, err := r.ledger.Get(
		fullKey(string(accountID), string(accountID), keyPublicKeyCount),
	)
	if err != nil {
		return nil, err
	}

	if countBytes == nil {
		return nil, fmt.Errorf("key count not set")
	}

	count := int(big.NewInt(0).SetBytes(countBytes).Int64())

	publicKeys = make([][]byte, count)

	for i := 0; i < count; i++ {
		publicKey, err := r.ledger.Get(
			fullKey(string(accountID), string(accountID), keyPublicKey(i)),
		)
		if err != nil {
			return nil, err
		}

		if publicKey == nil {
			return nil, fmt.Errorf("failed to retrieve key from account %s", accountID)
		}

		publicKeys[i] = cadence.NewBytes(publicKey)
	}

	return publicKeys, nil
}

func (r *RuntimeContext) setAccountPublicKeys(accountID []byte, publicKeys [][]byte) error {
	var existingCount int

	countBytes, err := r.ledger.Get(
		fullKey(string(accountID), string(accountID), keyPublicKeyCount),
	)
	if err != nil {
		return err
	}

	if countBytes != nil {
		existingCount = int(big.NewInt(0).SetBytes(countBytes).Int64())
	} else {
		existingCount = 0
	}

	newCount := len(publicKeys)

	r.ledger.Set(
		fullKey(string(accountID), string(accountID), keyPublicKeyCount),
		big.NewInt(int64(newCount)).Bytes(),
	)

	for i, publicKey := range publicKeys {
		if err := keys.ValidateEncodedPublicKey(publicKey); err != nil {
			return err
		}

		r.ledger.Set(
			fullKey(string(accountID), string(accountID), keyPublicKey(i)),
			publicKey,
		)
	}

	// delete leftover keys
	for i := newCount; i < existingCount; i++ {
		r.ledger.Delete(fullKey(string(accountID), string(accountID), keyPublicKey(i)))
	}

	return nil
}

// CheckCode checks the code for its validity.
func (r *RuntimeContext) CheckCode(address runtime.Address, code []byte) (err error) {
	return r.checkProgram(code, address)
}

// UpdateAccountCode updates the deployed code on an existing account.
//
// This function returns an error if the specified account does not exist or is
// not a valid signing account.
func (r *RuntimeContext) UpdateAccountCode(address runtime.Address, code []byte, checkPermission bool) (err error) {
	accountID := address[:]

	if checkPermission && !r.isValidSigningAccount(address) {
		return fmt.Errorf("not permitted to update account with ID %x", accountID)
	}

	bal, err := r.ledger.Get(fullKey(string(accountID), "", keyBalance))
	if err != nil {
		return err
	}
	if bal == nil {
		return fmt.Errorf("account with ID %s does not exist", accountID)
	}

	r.ledger.Set(fullKey(string(accountID), string(accountID), keyCode), code)

	return nil
}

// GetAccount gets an account by address.
//
// The function returns nil if the specified account does not exist.
func (r *RuntimeContext) GetAccount(address flow.Address) *flow.Account {
	accountID := address.Bytes()

	balanceBytes, _ := r.ledger.Get(fullKey(string(accountID), "", keyBalance))
	if balanceBytes == nil {
		return nil
	}

	balanceInt := big.NewInt(0).SetBytes(balanceBytes)

	code, _ := r.ledger.Get(fullKey(string(accountID), string(accountID), keyCode))

	publicKeys, err := r.getAccountPublicKeys(accountID)
	if err != nil {
		panic(err)
	}

	accountPublicKeys := make([]flow.AccountPublicKey, len(publicKeys))
	for i, publicKey := range publicKeys {
		accountPublicKey, err := keys.DecodePublicKey(publicKey)
		if err != nil {
			panic(err)
		}

		accountPublicKeys[i] = accountPublicKey
	}

	return &flow.Account{
		Address: address,
		Balance: balanceInt.Uint64(),
		Code:    code,
		Keys:    accountPublicKeys,
	}
}

// ResolveImport imports code for the provided import location.
//
// This function returns an error if the import location is not an account address,
// or if there is no code deployed at the specified address.
func (r *RuntimeContext) ResolveImport(location runtime.Location) ([]byte, error) {
	addressLocation, ok := location.(runtime.AddressLocation)
	if !ok {
		return nil, errors.New("import location must be an account address")
	}

	address := flow.BytesToAddress(addressLocation)

	accountID := address.Bytes()

	code, err := r.ledger.Get(fullKey(string(accountID), string(accountID), keyCode))
	if err != nil {
		return nil, err
	}

	if code == nil {
		return nil, fmt.Errorf("no code deployed at address %x", accountID)
	}

	return code, nil
}

// Log captures a log message from the runtime.
func (r *RuntimeContext) Log(message string) {
	r.logs = append(r.logs, message)
}

// EmitEvent is called when an event is emitted by the runtime.
func (r *RuntimeContext) EmitEvent(event runtime.Event) {
	r.events = append(r.events, event)
}

func (r *RuntimeContext) isValidSigningAccount(address runtime.Address) bool {
	for _, accountAddress := range r.GetSigningAccounts() {
		if accountAddress == address {
			return true
		}
	}

	return false
}

// checkProgram checks the given code for syntactic and semantic correctness.
func (r *RuntimeContext) checkProgram(code []byte, address runtime.Address) error {
	if code == nil {
		return nil
	}

	location := runtime.AddressLocation(address[:])

	return r.checker(code, location)
}

const (
	keyLatestAccount  = "latest_account"
	keyBalance        = "balance"
	keyCode           = "code"
	keyPublicKeyCount = "public_key_count"
)

func fullKey(owner, controller, key string) string {
	return strings.Join([]string{owner, controller, key}, "__")
}

func keyPublicKey(index int) string {
	return fmt.Sprintf("public_key_%d", index)
}
