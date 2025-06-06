package migrations

import (
	"deepjudge/utils"
)

func AddUserSignature() error {
	return utils.DB.Exec("ALTER TABLE users ADD COLUMN signature TEXT DEFAULT ''").Error
}
