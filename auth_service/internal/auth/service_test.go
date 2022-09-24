package auth

import (
	"context"
	"fmt"
	"testing"
)

func TestService_CheckToken(t *testing.T) {
	service := service{
		logger:  nil,
		storage: nil,
	}
	ctx := context.Background()
	tokens, err := service.CreateToken(ctx, "new_id_in_hex123")
	if err != nil {
		t.Fail()
	}

	t.Run("normal token test", func(t *testing.T) {
		err = service.CheckToken(ctx, tokens.AccessToken)
		if err != nil {
			t.Fail()
		}
	})

	t.Run("blank access token test", func(t *testing.T) {
		tokens.AccessToken = ""

		err = service.CheckToken(ctx, tokens.AccessToken)
		if err == nil {
			t.Fail()
		}
	})

}

func TestService_CreateToken(t *testing.T) {
	service := service{
		logger:  nil,
		storage: nil,
	}
	ctx := context.Background()

	t.Run("correct token test", func(t *testing.T) {
		token, err := service.CreateToken(ctx, "userid123inhex")
		if err != nil {
			t.Fail()
		}
		fmt.Printf("token: %s\n", token)
	})

	t.Run("empty userID test", func(t *testing.T) {
		_, err := service.CreateToken(ctx, "")
		if err == nil {
			t.Fail()
		}
	})
}
