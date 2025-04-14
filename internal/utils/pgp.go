package utils

func EncryptPGP(data string) (string, error) {
	// В реальном приложении — шифрование через OpenPGP
	return "[PGP_ENCRYPTED]:" + data, nil
}
