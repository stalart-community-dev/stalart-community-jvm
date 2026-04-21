// Package javamatch decides whether an invoked java.exe / javaw.exe is the
// STALART bundled runtime by comparing SHA-256 of the on-disk image with the
// reference copy under %AppData%\STALART\updates\...
package javamatch

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// RuntimeDirName is the folder under %AppData%\STALART\updates\ that
// contains bin\javaw.exe. Update when the launcher ships another JDK build.
const RuntimeDirName = "java25-windows-x86-64"

func referenceBin(name string) (string, error) {
	roaming, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("user config dir: %w", err)
	}
	return filepath.Join(roaming, "STALART", "updates", RuntimeDirName, "bin", name), nil
}

// ReferenceJavaw returns the absolute path to the known-good javaw.exe
// under %AppData%\STALART\updates\<RuntimeDirName>\bin\ .
func ReferenceJavaw() (string, error) {
	return referenceBin("javaw.exe")
}

// ReferenceJava returns the absolute path to java.exe in the same STALART
// runtime bin directory (Forge/Gradle launchers often use java.exe).
func ReferenceJava() (string, error) {
	return referenceBin("java.exe")
}

func sha256File(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// Match reports whether invokedExe is the STALART-bundled java.exe or
// javaw.exe (SHA-256 equals the reference under AppData\...\updates\...\bin).
// If the matching reference file is missing, it returns (false, err) so
// callers can log once. Any other image name returns (false, nil).
func Match(invokedExe string) (ok bool, err error) {
	base := filepath.Base(invokedExe)
	var refPath string
	var refErr error
	switch {
	case strings.EqualFold(base, "javaw.exe"):
		refPath, refErr = ReferenceJavaw()
	case strings.EqualFold(base, "java.exe"):
		refPath, refErr = ReferenceJava()
	default:
		return false, nil
	}
	if refErr != nil {
		return false, refErr
	}
	if _, err := os.Stat(refPath); err != nil {
		if os.IsNotExist(err) {
			return false, fmt.Errorf("reference JVM image missing: %s", refPath)
		}
		return false, err
	}
	want, err := sha256File(refPath)
	if err != nil {
		return false, fmt.Errorf("hash reference image: %w", err)
	}
	got, err := sha256File(invokedExe)
	if err != nil {
		return false, fmt.Errorf("hash invoked image: %w", err)
	}
	if subtle.ConstantTimeCompare(want, got) == 1 {
		return true, nil
	}
	return false, nil
}
