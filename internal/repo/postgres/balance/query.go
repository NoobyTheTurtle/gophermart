package balance

const (
	createBalanceQuery = `
	INSERT INTO user_balances (user_id, current, withdrawn, updated_at)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (user_id) DO NOTHING
`

	getBalanceByUserIDQuery = `
	SELECT user_id, current, withdrawn, updated_at
	FROM user_balances
	WHERE user_id = $1
`

	updateBalanceQuery = `
	UPDATE user_balances
	SET current = $1, withdrawn = $2, updated_at = $3
	WHERE user_id = $4
`

	addAccrualQuery = `
	UPDATE user_balances
	SET current = current + $1, updated_at = $2
	WHERE user_id = $3
`
)
