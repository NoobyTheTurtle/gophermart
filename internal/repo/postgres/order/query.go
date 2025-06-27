package order

const (
	createOrderQuery = `
	INSERT INTO orders (user_id, number, status, uploaded_at)
	VALUES ($1, $2, $3, $4)
	RETURNING id
`

	getOrderByNumberQuery = `
	SELECT id, user_id, number, status, accrual, uploaded_at
	FROM orders
	WHERE number = $1
`

	getOrdersByUserIDQuery = `
	SELECT id, user_id, number, status, accrual, uploaded_at
	FROM orders
	WHERE user_id = $1
	ORDER BY uploaded_at DESC
`

	updateOrderStatusQuery = `
	UPDATE orders
	SET status = $1, accrual = $2
	WHERE number = $3
`

	getOrdersForProcessingQuery = `
	SELECT id, user_id, number, status, accrual, uploaded_at
	FROM orders
	WHERE status IN ('NEW', 'PROCESSING')
	ORDER BY uploaded_at ASC
`
)
