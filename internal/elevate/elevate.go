package elevate

import (
	"errors"
	"os"
)

//func ElevatedOrRerun() {
//	if !isAdmin() {
//		runMeElevated()
//		common.Exit(0)
//	}
//}

// https://stackoverflow.com/questions/31558066/how-to-ask-for-administer-privileges-on-windows-with-go
//func runMeElevated() {
//	verb := "runas"
//	exe, _ := os.Executable()
//	args := strings.Join(os.Args[1:], " ")
//	cwd, _ := os.Getwd()
//
//	verbPtr, _ := syscall.UTF16PtrFromString(verb)
//	filePtr, _ := syscall.UTF16PtrFromString(exe)
//	argsPtr, _ := syscall.UTF16PtrFromString(args)
//	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
//	showCmd := int32(windows.SW_NORMAL)
//
//	err := windows.ShellExecute(0, verbPtr, filePtr, argsPtr, cwdPtr, showCmd)
//	if err != nil {
//		fmt.Println("Cannot elevate:", err)
//	}
//}

func ElevatedOrError() error {
	if !isAdmin() {
		return errors.New("please run this command as an administrator")
	}
	return nil
}

func isAdmin() bool {
	file, err := os.Open(`\\.\PHYSICALDRIVE0`)
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}
