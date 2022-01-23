package lib

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func printStackTrace(err interface{}) string {
	var pcAddrs [32]uintptr
	pos := runtime.Callers(3, pcAddrs[:]) // skip 3 callers
	/*
		<Panic>						 <---+	<== panic position
			|Recovery()					 |	caller 3
				|return func(...)		 |	caller 2
					|defer func() {...}	 |	caller 1
						|printStackTrace()	<== current position
	*/
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("\nPanic: %v\n\n", err))
	for _, pc := range pcAddrs[:pos] {
		fn := runtime.FuncForPC(pc)   // retrieve function from its PC address
		file, line := fn.FileLine(pc) // acquire calling position
		builder.WriteString(fmt.Sprintf("\tat %s(%s:%d)\n", fn.Name(), file, line))
	}
	return builder.String()
}

func Recovery() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				// successfully recovered from a panic
				log.Printf("[Panic] %s\n", printStackTrace(err))
				ctx.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		ctx.Next()
	}
}
