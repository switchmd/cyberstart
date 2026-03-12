package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Bat2ExeStep returns the Batch to Exe Converter installation step.
func Bat2ExeStep() Step {
	return Step{
		ID:   "bat2exe",
		Name: "Batch to Exe Converter",
		Run:  installBat2Exe,
	}
}

func installBat2Exe(logFn func(string)) error {
	logFn("Batch to Exe Converter를 다운로드합니다...")

	dest := filepath.Join(TempDir(), "Bat_To_Exe_Converter.zip")
	if err := DownloadFile(
		"https://ipfs.io/ipfs/QmPBp7wBSC9GukPUcp7LXFCGXBvc2e45PUfWUbCJzuLG65?filename=Bat_To_Exe_Converter.zip",
		dest, logFn,
	); err != nil {
		return fmt.Errorf("Bat2Exe 다운로드 실패: %w", err)
	}

	logFn("압축을 해제합니다...")
	extractPath := filepath.Join(os.Getenv("USERPROFILE"), "Desktop", "Bat_To_Exe_Converter")

	cmd := exec.Command("powershell", "-NoProfile", "-Command",
		fmt.Sprintf(`Expand-Archive -Path '%s' -DestinationPath '%s' -Force`, dest, extractPath))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("압축 해제 실패: %w", err)
	}

	logFn("바탕화면에 압축 해제 완료")
	return nil
}
