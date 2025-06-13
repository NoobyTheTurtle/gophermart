package user

const createUserQuery = `
	INSERT INTO users (login, password_hash, created_at)
	VALUES ($1, $2, $3)
	RETURNING id
`

const getUserByLoginQuery = `
	SELECT id, login, password_hash, created_at
	FROM users
	WHERE login = $1
`

const getUserByIDQuery = `
	SELECT id, login, password_hash, created_at
	FROM users
	WHERE id = $1
`
