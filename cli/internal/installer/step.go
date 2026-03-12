package installer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

// Step represents a single installation step.
type Step struct {
	ID   string
	Name string
	Run  func(logFn func(string)) error
}

// TempDir returns (and creates) the temporary directory for downloads.
func TempDir() string {
	dir := filepath.Join(os.TempDir(), "CyberstartInstallers")
	os.MkdirAll(dir, 0o755)
	return dir
}

// DownloadFile downloads a file from url to dest, reporting progress via logFn.
func DownloadFile(url, dest string, logFn func(string)) error {
	logFn(fmt.Sprintf("다운로드 중: %s", filepath.Base(dest)))

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("다운로드 요청 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("다운로드 실패 (HTTP %d)", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("파일 생성 실패: %w", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("파일 저장 실패: %w", err)
	}

	logFn("다운로드 완료")
	return nil
}

// RunElevated runs an executable with administrator privileges via ShellExecute "runas".
// It blocks until the spawned process exits.
func RunElevated(exe string, args ...string) error {
	verb, _ := syscall.UTF16PtrFromString("runas")
	exePath, _ := syscall.UTF16PtrFromString(exe)
	params, _ := syscall.UTF16PtrFromString(strings.Join(args, " "))
	cwd, _ := syscall.UTF16PtrFromString(".")

	shell32 := syscall.NewLazyDLL("shell32.dll")
	shellExecuteEx := shell32.NewProc("ShellExecuteExW")

	type SHELLEXECUTEINFO struct {
		CbSize       uint32
		FMask        uint32
		Hwnd         uintptr
		LpVerb       *uint16
		LpFile       *uint16
		LpParameters *uint16
		LpDirectory  *uint16
		NShow        int32
		HInstApp     uintptr
		LpIDList     uintptr
		LpClass      uintptr
		HkeyClass    uintptr
		DwHotKey     uint32
		HIcon        uintptr
		HProcess     syscall.Handle
	}

	const (
		SEE_MASK_NOCLOSEPROCESS = 0x00000040
		SW_SHOWNORMAL           = 1
	)

	info := SHELLEXECUTEINFO{
		FMask:        SEE_MASK_NOCLOSEPROCESS,
		LpVerb:       verb,
		LpFile:       exePath,
		LpParameters: params,
		LpDirectory:  cwd,
		NShow:        SW_SHOWNORMAL,
	}
	info.CbSize = uint32(unsafe.Sizeof(info))

	ret, _, _ := shellExecuteEx.Call(uintptr(unsafe.Pointer(&info)))
	if ret == 0 {
		return fmt.Errorf("ShellExecuteEx 실패")
	}

	if info.HProcess != 0 {
		syscall.WaitForSingleObject(info.HProcess, syscall.INFINITE)
		syscall.CloseHandle(info.HProcess)
	}

	return nil
}

// CleanupTempDir removes the temporary download directory.
func CleanupTempDir() {
	os.RemoveAll(TempDir())
}

// AllSteps returns every available installation step.
func AllSteps() []Step {
	return []Step{
		ChromeStep(),
		KollusStep(),
		PlaynPlayStep(),
		VSCodeStep(),
		Bat2ExeStep(),
	}
}
