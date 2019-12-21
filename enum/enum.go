package enum

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type Enum string

// Scan set value of enum
func (enum *Enum) Scan(value interface{}) error {
	*enum = Enum(value.([]byte))
	return nil
}

// Value return enum value
func (enum Enum) Value() (driver.Value, error) {
	return string(enum), nil
}

// GeneratePGScript return postgres script for generate enum
func GeneratePGScript(name string, params ...string) string {
	return fmt.Sprintf(`
		do $$
		begin
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = '%v') THEN
			CREATE TYPE %v AS ENUM ('%v');
			END IF;
		end
		$$`, name, name, strings.Join(params, "', '"))
}
