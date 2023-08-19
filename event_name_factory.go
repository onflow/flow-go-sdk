package flow

import (
	"fmt"
)

type eventNameFactory struct {
	address      string
	contractName string
	eventName    string
}

func (f eventNameFactory) WithAddressString(address string) eventNameFactory {
	f.address = address
	return f
}

func (f eventNameFactory) WithAddress(address Address) eventNameFactory {
	f.address = address.Hex()
	return f
}

func (f eventNameFactory) WithContractName(contract string) eventNameFactory {
	f.contractName = contract
	return f
}

func (f eventNameFactory) WithEventName(event string) eventNameFactory {
	f.eventName = event
	return f
}

func (f eventNameFactory) Build() string {
	return fmt.Sprintf("A.%s.%s.%s", f.address, f.contractName, f.eventName)
}

// NewEvent helper function for constructing event names
func NewEvent() eventNameFactory {
	return eventNameFactory{}
}
