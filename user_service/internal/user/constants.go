package user

import (
	"time"
)

const (
	deleteRegisterNumberTime = 5 * time.Minute
	twilioStatusApproved     = "approved"
)

var verificationMessage = "Ваш код для входа в <название приложения> %s"
