package flow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventNameFactory(t *testing.T) {
	assert.Equal(t, "A.7e60df042a9c0868.FlowToken.AccountCreated", NewEvent().
		WithEventName("AccountCreated").
		WithAddressString("7e60df042a9c0868").
		WithContractName("FlowToken").
		Build())
}
