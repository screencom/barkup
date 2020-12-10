package barkup

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var (
	// PGDumpCmd is the path to the `pg_dump` executable
	PGDumpCmd = "pg_dump"
)

// Postgres is an `Exporter` interface that backs up a Postgres database via the `pg_dump` command
type Postgres struct {
	// DB Host (e.g. 127.0.0.1)
	Host string
	// DB Port (e.g. 5432)
	Port string
	// DB Name
	DB string
	// Connection Username
	Username string
	// Extra pg_dump options
	// e.g []string{"--inserts"}
	Options []string
}

// Export produces a `pg_dump` of the specified database, and creates a gzip compressed tarball archive.
func (x Postgres) Export() *ExportResult {
	optionString := strings.Join(x.dumpOptions(), "|")
	var result *ExportResult
	if strings.Contains(optionString, "-Fp") {
		result = &ExportResult{MIME: "text/plain"}
	} else {
		result = &ExportResult{MIME: "application/x-tar"}
	}

	// Export without extension. The manual renaming can add an extension later on
	result.Path = fmt.Sprintf(`bu_%v_%v`, x.DB, time.Now().Unix())
	options := append(x.dumpOptions(), fmt.Sprintf(`-f%v`, result.Path))
	out, err := exec.Command(PGDumpCmd, options...).Output()
	if err != nil {
		result.Error = makeErr(err, string(out))
	}
	result.Stdout = string(out)
	return result
}

func (x Postgres) dumpOptions() []string {
	options := x.Options

	if x.DB != "" {
		options = append(options, fmt.Sprintf(`-d%v`, x.DB))
	}

	if x.Host != "" {
		options = append(options, fmt.Sprintf(`-h%v`, x.Host))
	}

	if x.Port != "" {
		options = append(options, fmt.Sprintf(`-p%v`, x.Port))
	}

	if x.Username != "" {
		options = append(options, fmt.Sprintf(`-U%v`, x.Username))
	}

	return options
}
