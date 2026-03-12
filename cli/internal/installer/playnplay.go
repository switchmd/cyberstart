package installer

import (
	"fmt"
	"os"
	"path/filepath"
)

// PlaynPlayStep returns the PlaynPlay Player installation step.
func PlaynPlayStep() Step {
	return Step{
		ID:   "playnplay",
		Name: "PlaynPlay Player",
		Run:  installPlaynPlay,
	}
}

func installPlaynPlay(logFn func(string)) error {
	logFn("PlaynPlay 설치 상태를 확인합니다...")

	playnplayPaths := []string{
		filepath.Join(os.Getenv("LOCALAPPDATA"), `PlaynPlay\PlaynPlay.exe`),
		`C:\Program Files\PlaynPlay\PlaynPlay.exe`,
		`C:\Program Files (x86)\PlaynPlay\PlaynPlay.exe`,
	}

	for _, p := range playnplayPaths {
		if _, err := os.Stat(p); err == nil {
			logFn("PlaynPlay가 이미 설치되어 있습니다. 건너뜁니다.")
			return nil
		}
	}

	dest := filepath.Join(TempDir(), "PlaynPlaySetup.exe")
	if err := DownloadFile(
		"https://pnp-appdn.cdnetworks.com/releases/download/1.0.33/PlaynPlay_1.0.33_x64-setup.exe",
		dest, logFn,
	); err != nil {
		return fmt.Errorf("PlaynPlay 다운로드 실패: %w", err)
	}

	logFn("PlaynPlay를 설치합니다...")
	if err := RunElevated(dest); err != nil {
		return fmt.Errorf("PlaynPlay 설치 실패: %w", err)
	}

	logFn("PlaynPlay 설치 완료")
	return nil
}
