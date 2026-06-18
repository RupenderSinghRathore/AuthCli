package main

import "github.com/pquerna/otp/totp"

func generateTOTP(username string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      AppName,
		AccountName: username,
	})
	if err != nil {
		return "", err
	}

	return key.Secret(), nil
}

func verifyTOTP(secret, code string) bool {
	return totp.Validate(code, secret)
}
