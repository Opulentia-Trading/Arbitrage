package validator

import "fmt"

type ValidateArb struct {
    validator Receiver
	assets    []string
	platforms []string
}

// Extends Command
func (v *ValidateArb) execute() {
	v.validator.execute_arb(assets, platforms)
}

func ValidateArb(r Receiver, assets string, []platforms []string) *ValidateArb {
	return &ValidateArb{
		receiver: r,
		assets: assets,
		platforms: platforms
	}
}