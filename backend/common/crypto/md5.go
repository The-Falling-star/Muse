package crypto

import "golang.org/x/crypto/bcrypt"

// Md5HashStr 对字符串进行哈希
func Md5HashStr(str string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyMd5Hash 验证md5哈希字符串
func VerifyMd5Hash(hashStr, sourceStr string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashStr), []byte(sourceStr))
}
