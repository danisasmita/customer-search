package message

const (
	Success             = "success"
	InternalServerError = "internal server error"
	BadRequest          = "bad request"
	Unauthorized        = "unauthorized"

	UserRegistered     = "user registered successfully"
	UserNotFound       = "user not found"
	InvalidCredentials = "invalid credentials"

	CustomerNotFound = "customer not found"

	NameRequired     = "name is required"
	EmailRequired    = "email is required"
	UsernameRequired = "username is required"
	PasswordRequired = "password is required"

	SearchCustomer = "Please provide at least name, email, or account_number for the search"
)
