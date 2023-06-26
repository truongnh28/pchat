package helper

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"github.com/whatvn/denny"
	"github.com/whatvn/denny/log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const ActorCtxKey = "actor"

func getActorFromContext(context context.Context) (string, bool) {
	if ctx, ok := context.(*denny.Context); ok {
		iActor, ok := ctx.Get(ActorCtxKey)
		if !ok {
			return "", false
		}
		return iActor.(string), true
	}

	iActor := context.Value(ActorCtxKey)
	if iActor == nil {
		return "", false
	}
	return iActor.(string), true
}

func GetUserAndLogger(ctx context.Context) (string, *log.Log) {
	actor, _ := getActorFromContext(ctx)
	logger := denny.GetLogger(ctx).WithField("actor", actor)
	return actor, logger
}

func GenOtp(secret string) (string, error) {
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}
	now := time.Now().Unix()
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(now/30))
	hash := hmac.New(sha1.New, key)
	hash.Write(buf)
	sum := hash.Sum(nil)
	offset := sum[len(sum)-1] & 0x0F
	code := binary.BigEndian.Uint32(sum[offset:offset+4]) & 0x7FFFFFFF
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
