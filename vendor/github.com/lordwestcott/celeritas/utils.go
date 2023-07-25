package celeritas

import (
	"fmt"
	"regexp"
	"runtime"
	"time"
)

// Cool Function to get the name of the function that called this function and the time it took to run.
func (c *Celeritas) LoadTime(start time.Time) {
	elapsed := time.Since(start)
	pc, _, _, _ := runtime.Caller(1) //Gets the caller of this function
	funcObj := runtime.FuncForPC(pc) //Gets the function object of the caller
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1") //Gets the name of the caller

	c.InfoLog.Println(fmt.Sprintf("Load Time: %s took %s", name, elapsed))
}
