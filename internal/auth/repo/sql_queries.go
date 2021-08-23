package repo

const (
	createUserQuery = `INSERT INTO users(user_id, username, password, created_at) VALUES($1, $2, $3, $4) RETURNING *`

	findUserByNameQuery = `SELECT * FROM users WHERE username = $1`
)
