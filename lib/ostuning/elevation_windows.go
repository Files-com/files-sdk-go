//go:build windows

package ostuning

import "golang.org/x/sys/windows"

func currentProcessElevated() bool {
	adminSID, err := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
	if err != nil {
		return false
	}
	member, err := windows.Token(0).IsMember(adminSID)
	return err == nil && member
}
