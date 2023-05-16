package pkg

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrNotExist           = &MyError{Code: 1, Message: "not exist"}
	ErrDbWrite            = &MyError{Code: 2, Message: "error db write"}
	ErrDbRead             = &MyError{Code: 3, Message: "error db read"}
	ErrEmptyFieldDetected = &MyError{Code: 4, Message: "empty field detected"}
	ErrCard               = &MyError{Code: 101, Message: "card substitution"}
	ErrNotIpDevice        = &MyError{Code: 102, Message: "not IP device in request"}
	ErrNotWiegandCode     = &MyError{Code: 103, Message: "not Wiegand code in request"}

	ErrBadRequest       = &MyError{Code: 400, Message: "bad request"}
	ErrNoRebootMethod   = &MyError{Code: 9999, Message: "no Reboot method"}
	ErrNotAuthenticated = &MyError{Code: 401, Message: "not Authenticated"}

	ErrDublicateExternalId = &MyError{Code: 1010, Message: "dublicate ExternalID"}
	ErrCreateExternalId    = &MyError{Code: 1011, Message: "error create ExternalID"}
	ErrListEmpty           = &MyError{Code: 1012, Message: "list empty"}
	ErrNotUserPhoto        = &MyError{Code: 1013, Message: "not photo for User"}
	ErrFaceMatchFoundUser  = &MyError{Code: 1014, Message: "face match found User"}

	ErrZoneNotExist = &MyError{Code: 1015, Message: "not exists zone"}
	ErrPermission   = &MyError{Code: 1100, Message: "insufficient access rights"}
)

type MyError struct {
	Code    int
	Message string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("Error: %d - %s", e.Code, e.Message)
}

func HandleError(c *gin.Context, err error) {
	var Err string
	if myErr, ok := err.(*MyError); ok {
		Err = myErr.Error()
	} else {
		Err = err.Error()
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": Err})
}
