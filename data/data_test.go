package data

// test data
var users = []User{
	{
		Name:     "Peter Jones",
		Email:    "peter@gmail.com",
		Password: "peter_pass",
		Role:     "user",
	},
	{
		Name:     "John Smith",
		Email:    "john@gmail.com",
		Password: "john_pass",
		Role:     "admin",
	},
}

func setup() {
	SessionDeleteAll()
	UserDeleteAll()
}
