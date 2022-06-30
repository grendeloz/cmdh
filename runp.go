// This type provides helper code for capturing information about a
// a program's run environment.

package cmdh

import (
	"os"
	"os/user"
	"strconv"
	"time"
)

// Run environment
type RunParameters struct {
	Tool      string
	Version   string
	StartTime time.Time
	Args      []string
	UserId    int
	UserName  string
	GroupId   int
	GroupName string
	HostName  string
}

// Return a record with execution parameters
func NewRunParameters() RunParameters {
	userId := os.Getuid()
	groupId := os.Getgid()

	// Systems that use LDAP for user management (e.g. Avalon) bork when
	// trying to get the names to match the UID/GID numbers so we are
	// going to silently ignore errors on those functions.
	userName := ""
	tmpUserName, err := user.LookupId(strconv.Itoa(userId))
	if err == nil {
		userName = tmpUserName.Name
	}
	groupName := ""
	tmpGroupName, err := user.LookupGroupId(strconv.Itoa(groupId))
	if err == nil {
		groupName = tmpGroupName.Name
	}

	hostName, err := os.Hostname()
	if err != nil {
		// If we can get a hostname, we will set to the error message
		hostName = err.Error()
	}

	// Setup and return RunParameters
	var run RunParameters
	run.StartTime = time.Now()
	run.Args = os.Args
	run.UserId = userId
	run.UserName = userName
	run.GroupId = groupId
	run.GroupName = groupName
	run.HostName = hostName
	return run
}
