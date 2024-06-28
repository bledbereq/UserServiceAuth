package suite

import (
	ssov1 "UserServiceAuth/gen/go"
	"UserServiceAuth/internal/config"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.GetPublicKeyClient
}
