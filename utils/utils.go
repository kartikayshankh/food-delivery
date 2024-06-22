package utils

type GenericResponse struct {
	Status           string                 `json:"status,omitempty"`
	Message          string                 `json:"message"`
	DeveloperMessage string                 `json:"developerMessage,omitempty"`
	Data             []string               `json:"data,omitempty"`
	UpdatedAt        string                 `json:"updatedAt,omitempty"`
	Error            string                 `json:"error,omitempty"`
	Code             int                    `json:"code,omitempty"`
	Request          map[string]interface{} `json:"request,omitempty"`
	Results          interface{}            `json:"results,omitempty"`
}

type ErrorHandler struct {
	DevMessage string                 `json:"developerMessage,omitempty"`
	Request    map[string]interface{} `json:"request,omitempty"`
	Response   map[string]interface{} `json:"response,omitempty"`
	UserId     string                 `json:"userId,omitempty"`
	Message    string                 `json:"message,omitempty"`
	Method     string                 `json:"method,omitempty"`
	Code       int                    `json:"code,omitempty"`
	Status     string                 `json:"status,omitempty"`
}

const (

	//Server Response
	INTERNAL_SERVER_ERROR = "Internal Server Error"
	SOMETHING_WENT_WRONG  = "Unable to process your request, please try again later"
	DATA_NOT_FOUND        = "Data Not Found"
)

type Role string

const (
	User  Role = "user"
	Rider Role = "rider"
)
