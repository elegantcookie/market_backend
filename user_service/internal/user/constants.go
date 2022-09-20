package user

import (
	"time"
)

const (
	deleteRegisterNumberTime = 5 * time.Minute
	twilioStatusApproved     = "approved"
	verificationCodeLength   = 6
	checkTokenURL            = "https://api.vk.com/method/secure.checkToken?token=%s&access_token=%s&v=%s"
	checkTokenVersion        = "5.81"
)

var verificationMessage = "Ваш код для входа в <название приложения> %s"
