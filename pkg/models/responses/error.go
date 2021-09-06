package responses

//ErrProblem resembles the minimal structure
//used by the specification https://datatracker.ietf.org/doc/html/rfc7807
//as a way to carry machine-readable details of errors in an HTTP response to avoid the need to
//define new error response formats for HTTP APIs.
type ErrProblem struct {
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}