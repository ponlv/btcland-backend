package utils

import (
	cryptoRand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"runtime/debug"
	"strconv"
	"time"
)

func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func EndOfDay(t time.Time) time.Time {
	startOfDay := StartOfDay(t)
	return startOfDay.Add(24 * 3600 * time.Second)
}

func FloatToBigInt(val float64, decimal int32) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)

	coin := big.NewFloat(math.Pow10(int(decimal)))
	bigval.Mul(bigval, coin)

	result := new(big.Int)
	result.SetString(fmt.Sprintf("%f", bigval), 10)
	return result
}

func I2JsonString(data interface{}) string {
	body, _ := json.Marshal(data)
	return string(body)
}

func Recover() {
	if err := recover(); err != nil {
		fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
	}
}

func GenerateCode(length int) string {
	code := ""
	for i := 0; i < length; i++ {
		code = code + strconv.Itoa(rand.Intn(9-0)+0)
	}
	return code
}

func ItoBool(value interface{}) (bool, error) {
	if value == nil {
		return false, errors.New("empty value")
	}
	v, err := value.(bool)
	if !err {
		return false, errors.New("error convert Interface to bool")
	}
	return v, nil
}

func ItoString(value interface{}) string {
	if value == nil {
		return ""
	}
	str := fmt.Sprintf("%v", value)
	return str
}

func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		for {
			num, err := cryptoRand.Int(cryptoRand.Reader, big.NewInt(int64(len(charset))))
			if err == nil {
				result[i] = charset[num.Int64()]
				break
			}
		}
	}
	return string(result)
}
