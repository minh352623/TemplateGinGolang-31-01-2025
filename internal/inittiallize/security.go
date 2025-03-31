package inittiallize

import (
	"ecom/global"
	"ecom/pkg/security"
)

func initSecurity() {
	global.SecurityService = &security.SecurityService{
		PrivKey: global.Config.Security.CryptoKeys.Asymmetric.PrivKey,
		PubKey:  global.Config.Security.CryptoKeys.Asymmetric.PubKey,
		AESKey:  global.Config.Security.CryptoKeys.Symmetric.AESKey,
	}
}
