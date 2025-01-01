package types

// Even though `http` package already has constants with these
// values, I thought it was better to create a separated type to distinguish
// any strings of "GET", "POST" and so on
type HTTPMethod string
