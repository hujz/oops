package oops

import (
	"io"
	"runtime"
)

// StackBufferSize max buffer size
const StackBufferSize = 4096

// PrintStack PrintStack
func PrintStack(writer io.Writer) {
	var buf [StackBufferSize]byte
	n := runtime.Stack(buf[:], false)
	io.WriteString(writer, string(buf[:n]))
}
