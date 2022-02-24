package main

import "fmt"

type Validator struct {}

func NewValidator() *Maintainer {
	return &Validator{}
}

func (v *Validator) execute_arb(pair string) {
    // Grab from DB here
    fmt.Println("Calculating :", pair)
    
    // Execute arb here
    fmt.Println("Executing :", pair)
}