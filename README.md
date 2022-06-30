# cmdh

Helper for cobra commander applications.

This file centralises initialisation and logging for a style of
cobra commander applications as used by grendeloz. It may suit or
interest nobody else.

## Usage

To use cmdh, you should call cmdh.Initialise from an init()  - probably
in cmd/root.go. This will add global (Persistent) flags to the cobra root
command ("rootCmd" in our example) and set strings for the name and
version of the application. 

```
func init() {
    cmdh.Initialise(rootCmd, "myapp", "v0.1.0-dev")
}
```

Add calls to start and stop logging around any cobra "CmdRun" functions
where you want logging. Note that this can not happen until after
cobra flags have been parsed, assuming that we are allowing users to
set values for logfile and loglevel.  For example:

```
var myCmd = &cobra.Command{
    Use:   "do-stuff",
    Short: "Do some stuff",
    Long:  `Do some clever stuff.`,
    Run: func(cmd *cobra.Command, args []string) {
        cmdh.StartLogging()
        myCmdRun(cmd, args)
        cmdh.FinishLogging()
    },
}
```

## RunParameters

This type helps capture information about the run environment. It is
part of the internal funvtioning of the cmdh package but can also be
used externally.

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
