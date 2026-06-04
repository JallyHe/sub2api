package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/Wei-Shaw/sub2api/migrations"
)

func main() {
	entries, _ := migrations.FS.ReadDir(".")
	vals := []string{}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".sql") {
			raw, _ := migrations.FS.ReadFile(e.Name())
			// Mirror migrations_runner.go: TrimSpace before SHA256
			content := strings.TrimSpace(string(raw))
			h := sha256.Sum256([]byte(content))
			checksum := hex.EncodeToString(h[:])
			fname := strings.ReplaceAll(e.Name(), "'", "''")
			vals = append(vals, fmt.Sprintf("('%s','%s')", fname, checksum))
		}
	}
	sql := "INSERT INTO schema_migrations(filename,checksum) VALUES " +
		strings.Join(vals, ",") +
		" ON CONFLICT(filename) DO UPDATE SET checksum=EXCLUDED.checksum;"
	fmt.Print(sql)
}
