package backup

import (
	"fmt"
	"path/filepath"
	"time"
)

type Monitor struct {
	// Paths with there associated hashes
	Paths       map[string]string
	Archiver    Archiver
	Destination string
}

// Now iterates over paths and act Archiver if hash changed
func (m *Monitor) Now() (int, error) {
	var counter int
	for path, lastHash := range m.Paths {
		newHash, err := DirHash(path)
		if err != nil {
			return counter, err
		}
		if newHash != lastHash {
			err := m.act(path)
			if err != nil {
				return counter, err
			}
			m.Paths[path] = newHash
			counter++
		}
	}
	return counter, nil
}

func (m *Monitor) act(path string) error {
	dirname := filepath.Base(path)
	filename := fmt.Sprintf("%s%s", time.Now().Format("2006-01-02 15:04:05"), m.Archiver.DestExt())
	return m.Archiver.Archiver(path, filepath.Join(m.Destination, dirname, filename))
}
