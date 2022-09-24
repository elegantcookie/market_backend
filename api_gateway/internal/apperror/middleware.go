package apperror

import (
	"api_gateway/internal/gateway/routes/auth_service"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type appHandler func(http.ResponseWriter, *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
	log.Println("got into middleware")
	return func(w http.ResponseWriter, r *http.Request) {
		var appErr *AppError
		err := h(w, r)
		if err != nil {
			if errors.As(err, &appErr) {
				if errors.Is(err, ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					w.Write(ErrNotFound.Marshal())
					return
				}
				err := err.(*AppError)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(err.Marshal())
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(systemError(err.Error()).Marshal())
		}
	}
}

func AuthMiddleware(h appHandler) http.HandlerFunc {
	//log.Println("got into auth middleware")
	return func(w http.ResponseWriter, r *http.Request) {
		var appErr *AppError
		headerVal := r.Header.Get("Authorization")
		if headerVal == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(ErrWrongToken.Marshal())
			return
		}

		authHeaderArr := strings.Split(headerVal, " ")
		if len(authHeaderArr) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(ErrWrongToken.Marshal())
			return
		}
		tokenString := authHeaderArr[1]
		tmp := make(map[string]string)
		tmp["token"] = tokenString
		bytes, err := json.Marshal(tmp)
		if err != nil {
			return
		}
		body := io.NopCloser(strings.NewReader(string(bytes)))
		url := fmt.Sprintf("http://%s:%s%s", auth_service.DockerName, auth_service.Port, auth_service.CheckTokenURL)
		log.Printf("url: %s", url)
		request, err := http.NewRequestWithContext(r.Context(), http.MethodPost, url, body)
		if err != nil {
			return
		}
		var client http.Client
		response, err := client.Do(request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("%v", err)
			//w.Write(err.Marshal())
			return
		}
		if response == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(ErrResponseIsNil.Marshal())
			return
		}
		if response.StatusCode != http.StatusOK {
			bytes, err = io.ReadAll(response.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				err := err.(*AppError)
				w.Write(err.Marshal())
				return
			}
			w.WriteHeader(response.StatusCode)
			w.Write(bytes)
			return
		}

		//if response.StatusCode != 200 {
		//	w.WriteHeader(http.StatusBadRequest)
		//	instance := ErrWrongStatusCode
		//	instance.DeveloperMessage = fmt.Sprintf("status code: %d", response.StatusCode)
		//	w.Write(instance.Marshal())
		//	return
		//}

		err = h(w, r)
		if err != nil {
			if errors.As(err, &appErr) {
				if errors.Is(err, ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					w.Write(ErrNotFound.Marshal())
					return
				}
				err := err.(*AppError)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(err.Marshal())
				return
			}
			w.WriteHeader(http.StatusTeapot)
			w.Write(systemError(err.Error()).Marshal())
			return
		}
	}
}
