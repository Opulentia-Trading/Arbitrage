package main

// The Invoker class is responsible for initiating requests. 
// This class stores a reference to a queue of command objects. 
// Invoker triggers that command instead of sending the request directly to the receiver. 
// Invoker doesn't create the command object, it gets a pre-created command from the client via the constructor.
type Invoker struct {
    commands []Command
}

func NewInvoker() *Invoker {
	return &Invoker{}
}

// Adds a new command to current queue of commands
func (i *Invoker) Add(command Command) {
    i.commands = append(i.commands, command)
}

// The ExecuteCommands method executes all the commands one by one
func (i *Invoker) ExecuteCommands() {
	for _, command := range i.commands {
		command.execute()
	}
    i.commands = []Command{}
}