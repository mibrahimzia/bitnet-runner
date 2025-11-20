package windows

import (
	"embed"
)

//go:embed bitnet.exe *.dll
var BitNetFS embed.FS

// GetFS returns the embedded filesystem containing the Windows binaries
func GetFS() embed.FS {
	return BitNetFS
}