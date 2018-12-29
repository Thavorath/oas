//@openapi: "a sample api"
//sample documentation
//on multiple lines!
//@version:0.2.0
//@apiKey: header api_key
//@security: api_key
//@server: "https://api.nirandas.com/"
package main


//User object
//@dto User
type User struct {
	//unique id
	ID int64 `json:"id"`
	//email
	EmailAddress string `json:"email"`
	//name
	Name string
}

//LoginResponseDTO response dto
//@dto
type LoginResponseDTO struct {
	//User auth token
	Token string `json:"auth_token"`
	//User information
	User *User `json:"user"`
	Foo  func(string) int
}

//LoginRequestDTO request
//@dto
type LoginRequestDTO struct {
	//login email
	Email string
	//login password
	Password string
	//list of scope
	Scope []Scope `json:"scope_values"`
	priv  string
}

//Scope test struct
//@dto
//This is about testing, all about tsting!
//then, what do you do???
type Scope struct {
	ID    int `json:"ID"`
	Value float32
}

//LoginAPI handler
//@route: GET "/api/account/login"
//Accepts the login credentials and returns the user object and token
//if valid login.
//@input: LoginRequestDTO
//Provide the login email and password to authenticate
//@response: 200 LoginResponseDTO
//Responds with the user auth token and user info.
//The user auth token is generated during login hence all previously issued auth tokens for this user gets invalidated
//@response: content="text/json" 200 Foo
//@response: 422 ValidationError
//@parameters: limit skip
//@security: -
func LoginAPI() {}

//@parameter: limit query
//Specifies the number of results to return.
type _ int64

//@parameter: skip query
//Specifies the number of results to skip.
type _ int64

//Contact test struct
//@dto: Contact
type Contact struct {
	Name  string
	Email string
}

//Person test struct
//@dto: Person
type Person struct {
	Contact
	Age int
}
