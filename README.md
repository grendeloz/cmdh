# cmdh

Helper for cobra commander applications.


## RunParameters

This type helps capture information about the run environment.

```
package main

import (
    "log"
    "github.com/grendeloz/cmdh"
)

var MyRp cmdh.RunParameters

func init() {
    cmdh.SetTool(`myapp`)
    cmdh.SetVersion(`v0.1.0dev`)
    MyRp = cmdh.NewRunParameters()
}

func main() {
    log.Info("Started: ", MyRp.StartTime)
}
```
