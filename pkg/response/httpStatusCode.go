package response

const (
	Success             = 200
	BadRequest          = 400
	Unauthorized        = 401
	InvalidRequest      = 402
	NotFound            = 404
	InternalServerError = 500
)

// message
var msg = map[int]string{
	Success:             "Success",
	BadRequest:          "Bad Request",
	Unauthorized:        "Unauthorized",
	NotFound:            "Not Found",
	InternalServerError: "Internal Server Error",
}
