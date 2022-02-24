package chip8

const StackSize = 16

type Stack struct {
	Values []uint16
}

func NewStack() *Stack {
	stack := new(Stack)
	stack.Clear()
	return stack
}

func (stack *Stack) Clear() {
	stack.Values = make([]uint16, 0, StackSize)
}

func (stack *Stack) Push(value uint16) {
	if len(stack.Values) == cap(stack.Values) {
		panic("chip8: stack overflow")
	}
	stack.Values = append(stack.Values, value)
}

func (stack *Stack) Pop() uint16 {
	if len(stack.Values) == 0 {
		panic("chip8: nothing to pop from stack")
	}
	n := len(stack.Values) - 1
	defer func() {
		stack.Values[n] = 0x00
		stack.Values = stack.Values[:n]
	}()
	return stack.Values[n]
}
