package user

import (
	"time"
)

const (
	deleteRegisterNumberTime = 5 * time.Minute
	twilioStatusApproved     = "approved"
	verificationCodeLength   = 6
)

var verificationMessage = "Ваш код для входа в <название приложения> %s"
