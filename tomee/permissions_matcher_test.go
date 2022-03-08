package tomee_test

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	"os"
)

type permissionsDontMatchError struct {
	fileInfo os.FileInfo
}

func (t permissionsDontMatchError) Error() string {
	return fmt.Sprintf("file mode is: %s", t.fileInfo.Mode().String())
}

//HavePermissions succeeds if a file exists and has the permissions on the file match the expected permissions.
//Actual must be a string representing the abs path to the file being checked.
func HavePermissions(mode os.FileMode) types.GomegaMatcher {
	return &HavePermissionsMatcher{mode: mode}
}

type HavePermissionsMatcher struct {
	expected interface{}
	err      error
	mode     os.FileMode
}

func (matcher *HavePermissionsMatcher) Match(actual interface{}) (success bool, err error) {
	actualFilename, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("HavePermissionsMatcher matcher expects a file path")
	}

	fileInfo, err := os.Stat(actualFilename)
	if err != nil {
		matcher.err = err
		return false, nil
	}

	if fileInfo.Mode() != matcher.mode {
		matcher.err = permissionsDontMatchError{fileInfo}
		return false, nil
	}
	return true, nil
}

func (matcher *HavePermissionsMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("to have permissions %s: %s", matcher.mode, matcher.err))
}

func (matcher *HavePermissionsMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("not to have permissions"))
}
