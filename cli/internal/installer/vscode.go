package installer

import (
	"fmt"
	"os"
	"path/filepath"
)

// VSCodeStep returns the Visual Studio Code installation step.
func VSCodeStep() Step {
	return Step{
		ID:   "vscode",
		Name: "Visual Studio Code",
		Run:  installVSCode,
	}
}

func installVSCode(logFn func(string)) error {
	logFn("VSCode 설치 상태를 확인합니다...")

	vscodePaths := []string{
		filepath.Join(os.Getenv("LOCALAPPDATA"), `Programs\Microsoft VS Code\Code.exe`),
		`C:\Program Files\Microsoft VS Code\Code.exe`,
	}

	for _, p := range vscodePaths {
		if _, err := os.Stat(p); err == nil {
			logFn("Visual Studio Code가 이미 설치되어 있습니다. 건너뜁니다.")
			return nil
		}
	}

	dest := filepath.Join(TempDir(), "VSCodeSetup.exe")
	if err := DownloadFile(
		"https://code.visualstudio.com/sha/download?build=stable&os=win32-x64-user",
		dest, logFn,
	); err != nil {
		return fmt.Errorf("VSCode 다운로드 실패: %w", err)
	}

	logFn("VSCode를 설치합니다...")
	if err := RunElevated(dest, "/VERYSILENT", "/MERGETASKS=!runcode"); err != nil {
		return fmt.Errorf("VSCode 설치 실패: %w", err)
	}

	logFn("VSCode 설치 완료")
	return nil
}
