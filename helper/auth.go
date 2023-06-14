package helper

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func GenOtp(secret string) (string, error) {
	// Chuyển đổi khóa bí mật từ base32 sang byte slice
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	// Tính toán số giây hiện tại
	now := time.Now().Unix()

	// Chia số giây cho 30 để lấy thời gian tính bằng 30 giây
	// Sau đó, chuyển đổi số giây thành byte slice theo định dạng big-endian
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(now/30))

	// Tạo mã HMAC-SHA1 từ khóa bí mật và byte slice thời gian
	hash := hmac.New(sha1.New, key)
	hash.Write(buf)
	sum := hash.Sum(nil)

	// Lấy 4 bit cuối của mã HMAC-SHA1 để tính toán chỉ số offset
	offset := sum[len(sum)-1] & 0x0F

	// Chuyển đổi 4 byte tiếp theo của mã HMAC-SHA1 thành một số nguyên 32-bit theo định dạng big-endian
	code := binary.BigEndian.Uint32(sum[offset:offset+4]) & 0x7FFFFFFF

	// Định dạng mã OTP thành chuỗi 6 chữ số
	otp := fmt.Sprintf("%06d", code%1000000)

	return otp, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
