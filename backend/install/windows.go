package install

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	ps "github.com/bhendo/go-powershell"
	"github.com/bhendo/go-powershell/backend"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	log "github.com/sirupsen/logrus"
)

func createWindowsShortcuts(src, dst string) error {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	defer ole.CoUninitialize()

	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer oleShellObject.Release()
	wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer wshell.Release()
	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", dst)
	if err != nil {
		return err
	}
	idispatch := cs.ToIDispatch()
	oleutil.PutProperty(idispatch, "TargetPath", src)
	oleutil.CallMethod(idispatch, "Save")

	return nil
}

func installJava() error {
	// local backend
	back := &backend.Local{}

	// start a local powershell process
	shell, err := ps.New(back)
	if err != nil {
		return fmt.Errorf("unable to initialize PowerShell: %v", err)
	}
	defer shell.Exit()

	// setting the execution policy to the current user
	stdout, stderr, err := shell.Execute("Set-ExecutionPolicy RemoteSigned -scope CurrentUser")
	if err != nil {
		return fmt.Errorf("unable to set ExecutionPolicy to CurrentUser. %v", err)
	}
	if stderr != "" {
		return fmt.Errorf("unable to set ExecutionPolicy to CurrentUser. stderr: %s", stderr)
	}
	log.Infof("setting ExecutionPolicy to CurrentUser. stdout: %s", stdout)

	// installing scoop package manager for windows: https://github.com/lukesampson/scoop
	stdout, stderr, err = shell.Execute("iwr -useb get.scoop.sh | iex")
	if err != nil {
		return fmt.Errorf("unable to install scoop. %v", err)
	}
	if stderr != "" {
		return fmt.Errorf("errors occured when installing scoop. stderr: %s", stderr)
	}
	log.Infof("installing scoop. stdout: %s", stdout)

	// installing git as a dependancy
	stdout, stderr, err = shell.Execute("scoop install git")
	if err != nil {
		return fmt.Errorf("unable to install git using scoop. %v", err)
	}
	if stderr != "" {
		return fmt.Errorf("errors occured when installing git through scoop. stderr: %s", stderr)
	}
	log.Infof("installing git using scoop. stdout: %s", stdout)

	// adding java bucket to scoop https://github.com/lukesampson/scoop/wiki/Java
	stdout, stderr, err = shell.Execute("scoop bucket add java")
	if err != nil {
		return fmt.Errorf("unable to add java bucket to scoop. %v", err)
	}
	if stderr != "" {
		return fmt.Errorf("errors occured while adding the java bucket to scoop. stderr: %s", stderr)
	}
	log.Infof("adding java bucket to scoop. stdout: %s", stdout)

	// installing adoptopenjdk-hotspot
	stdout, stderr, err = shell.Execute("scoop install adoptopenjdk-hotspot")
	if err != nil {
		return fmt.Errorf("unable to install java using scoop. %v", err)
	}
	if stderr != "" {
		return fmt.Errorf("errors occured when installing java through scoop. stderr: %s", stderr)
	}
	log.Infof("installing java using scoop. stdout: %s", stdout)

	// clean up excessive dependencies
	stdout, stderr, err = shell.Execute("scoop uninstall git")
	if err != nil {
		return fmt.Errorf("unable to uninstall git using scoop. %v", err)
	}
	if stderr != "" {
		return fmt.Errorf("errors occured when uninstalling git through scoop. stderr: %s", stderr)
	}

	log.Infof("uninstalling git using scoop. stdout: %s", stdout)

	return err

}

func javaInstalled() bool {
	javaPath, err := detectJavaPath()
	if javaPath == "" {
		return false
	}
	if err != nil {
		log.Errorln(err)
	}
	var javaInstalled bool
	if javaPath[len(javaPath)-9:] != "javaw.exe" {
		javaInstalled = false
	} else {
		javaInstalled = true
	}
	return javaInstalled
}

func detectJavaPath() (string, error) {

	var jwPath string

	cmd := exec.Command("cmd", "/c", "where", "java")
	log.Infoln("Running command: ", cmd)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out    // Captures STDOUT
	cmd.Stderr = &stderr // Captures STDERR

	err := cmd.Run()
	if err != nil {
		errFormatted := fmt.Sprint(err) + ": " + stderr.String()
		err = fmt.Errorf("unable to run command: %s", errFormatted)

		return "", err
	}
	jPath := out.String() // May contain multiple
	if jPath == "" {
		err = errors.New("unable to find Java Installation")
		return "", err
	}
	s := strings.Split(strings.Replace(jPath, "\r\n", "\n", -1), "\n")
	jwPath = string(s[0][:len(s[0])-4]) + "w.exe" // Shifting to javaw.exe
	if s[1] != "" {
		jwPath = string(s[1][:len(s[1])-4]) + "w.exe" // Shifting to javaw.exe
		log.Infoln("Detected a secondary java path. Using that over the first one.")
	}
	log.Infoln("Java path selected: " + jwPath)
	log.Debugln(cmd)
	return jwPath, nil

}
