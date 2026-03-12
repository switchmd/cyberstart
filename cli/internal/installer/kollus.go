package installer

import (
	"fmt"
	"os"
	"path/filepath"
)

// KollusStep returns the Kollus PC Player Agent installation step.
func KollusStep() Step {
	return Step{
		ID:   "kollus",
		Name: "Kollus PC Player Agent",
		Run:  installKollus,
	}
}

func installKollus(logFn func(string)) error {
	logFn("Kollus Agent 설치 상태를 확인합니다...")

	kollusPaths := []string{
		filepath.Join(os.Getenv("LOCALAPPDATA"), `Kollus\KollusAgent.exe`),
		`C:\Program Files\Kollus\KollusAgent.exe`,
		`C:\Program Files (x86)\Kollus\KollusAgent.exe`,
	}

	for _, p := range kollusPaths {
		if _, err := os.Stat(p); err == nil {
			logFn("Kollus Agent가 이미 설치되어 있습니다. 건너뜁니다.")
			return nil
		}
	}

	dest := filepath.Join(TempDir(), "KollusAgentSetup.exe")
	if err := DownloadFile(
		"https://v.kr.kollus.com/pc_player_install/agent",
		dest, logFn,
	); err != nil {
		return fmt.Errorf("Kollus Agent 다운로드 실패: %w", err)
	}

	logFn("Kollus Agent를 설치합니다...")
	if err := RunElevated(dest); err != nil {
		return fmt.Errorf("Kollus Agent 설치 실패: %w", err)
	}

	logFn("Kollus Agent 설치 완료")
	return nil
}
