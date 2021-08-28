package repo

const (
	createUserQuery = `INSERT INTO auth.users(user_id, username, password, created_at) VALUES($1, $2, $3, $4) RETURNING *`

	findUserByNameQuery = `SELECT * FROM auth.users WHERE username = $1`
)
