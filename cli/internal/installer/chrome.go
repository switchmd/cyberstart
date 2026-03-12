package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// ChromeStep returns the Google Chrome installation step.
func ChromeStep() Step {
	return Step{
		ID:   "chrome",
		Name: "Google Chrome",
		Run:  installChrome,
	}
}

func installChrome(logFn func(string)) error {
	logFn("Chrome 설치 상태를 확인합니다...")

	chromePaths := []string{
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		filepath.Join(os.Getenv("LOCALAPPDATA"), `Google\Chrome\Application\chrome.exe`),
	}

	var installedPath string
	for _, p := range chromePaths {
		if _, err := os.Stat(p); err == nil {
			installedPath = p
			break
		}
	}

	needsInstall := true

	if installedPath != "" {
		logFn("Chrome이 설치되어 있습니다. 버전을 확인합니다...")
		version, err := getChromeVersion(installedPath)
		if err != nil {
			logFn("버전 확인 실패. 재설치합니다...")
		} else {
			major := getMajorVersion(version)
			if major >= 130 {
				logFn(fmt.Sprintf("최신 버전 (v%s) 설치됨. 건너뜁니다.", version))
				needsInstall = false
			} else {
				logFn(fmt.Sprintf("레거시 버전 (v%s) 감지. 업데이트합니다.", version))
				uninstallChrome(logFn)
			}
		}
	} else {
		logFn("Chrome이 설치되지 않았습니다.")
	}

	if !needsInstall {
		return nil
	}

	// Download
	dest := filepath.Join(TempDir(), "ChromeInstaller.exe")
	if err := DownloadFile(
		"https://dl.google.com/chrome/install/375.126/chrome_installer.exe",
		dest, logFn,
	); err != nil {
		return fmt.Errorf("Chrome 다운로드 실패: %w", err)
	}

	// Install
	logFn("Chrome을 설치합니다...")
	if err := RunElevated(dest, "/install", "--do-not-launch-chrome"); err != nil {
		return fmt.Errorf("Chrome 설치 실패: %w", err)
	}

	logFn("Chrome 설치 완료")
	return nil
}

func getChromeVersion(chromePath string) (string, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-Command",
		fmt.Sprintf(`(Get-Item '%s').VersionInfo.ProductVersion`, chromePath))
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func getMajorVersion(version string) int {
	parts := strings.SplitN(version, ".", 2)
	if len(parts) == 0 {
		return 0
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}
	return major
}

func uninstallChrome(logFn func(string)) {
	logFn("레거시 Chrome을 제거합니다...")

	// Kill Chrome processes
	exec.Command("taskkill", "/f", "/im", "chrome.exe").Run()
	exec.Command("taskkill", "/f", "/im", "GoogleUpdate.exe").Run()

	// WMIC uninstall
	logFn("Windows 설치 프로그램을 통한 제거 시도 중...")
	exec.Command("wmic", "product", "where", "name like '%Google Chrome%'",
		"call", "uninstall", "/nointeractive").Run()

	// Direct folder cleanup
	logFn("잔여 파일 정리 중...")
	for _, folder := range []string{
		`C:\Program Files\Google`,
		`C:\Program Files (x86)\Google`,
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Google"),
	} {
		os.RemoveAll(folder)
	}

	logFn("레거시 Chrome 제거 완료")
}
