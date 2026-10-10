package member

const postgresCreateTableSQL = `CREATE TABLE IF NOT EXISTS admin_member (
  id bigserial PRIMARY KEY,
  username varchar(128) NOT NULL,
  password_hash varchar(128) NOT NULL,
  name varchar(32) NOT NULL DEFAULT '',
  last_login_ip varchar(64) NOT NULL DEFAULT '',
  last_login_at timestamptz NULL,
  token_version bigint NOT NULL DEFAULT 0,
  status smallint NOT NULL DEFAULT 1,
  deleted_at timestamptz NULL,
  created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT uk_admin_member_username UNIQUE (username)
)`

const mysqlCreateTableSQL = "CREATE TABLE IF NOT EXISTS admin_member (" +
	"id bigint NOT NULL AUTO_INCREMENT," +
	"username varchar(128) NOT NULL," +
	"password_hash varchar(128) NOT NULL," +
	"name varchar(32) NOT NULL DEFAULT ''," +
	"last_login_ip varchar(64) NOT NULL DEFAULT ''," +
	"last_login_at datetime NULL," +
	"token_version bigint NOT NULL DEFAULT 0," +
	"status smallint NOT NULL DEFAULT 1," +
	"deleted_at datetime NULL," +
	"created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP," +
	"updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
	"PRIMARY KEY (id)," +
	"UNIQUE KEY uk_admin_member_username (username)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci"

func createTableSQL(driver string) (string, error) {
	switch driver {
	case "postgres":
		return postgresCreateTableSQL, nil
	case "mysql":
		return mysqlCreateTableSQL, nil
	default:
		return "", ErrInvalidConfig
	}
}
