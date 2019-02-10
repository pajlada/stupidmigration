package stupidmigration

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type migration struct {
	path    string
	version uint64
}

func (m *migration) readData() (data string, err error) {
	file, err := os.Open(m.path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		data += line + "\n"
	}

	return
}

func (m *migration) Migrate(db *sql.DB) error {
	fmt.Println("Applying migration", m.version)

	data, err := m.readData()
	if err != nil {
		return err
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	if _, err = tx.Exec(data); err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`UPDATE migrations SET version=$1`, strconv.FormatUint(m.version, 10))
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
