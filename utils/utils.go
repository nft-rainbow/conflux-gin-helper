package utils

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"runtime/debug"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"

	"github.com/Conflux-Chain/go-conflux-sdk/types/cfxaddress"
	"github.com/Conflux-Chain/go-conflux-sdk/types/unit"
	"github.com/ethereum/go-ethereum/common"
)

var (
	U256Max = MustParseStrToBig("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
)

func Bytes2Hex(data []byte) string {
	return "0x" + common.Bytes2Hex(data)
}

func CheckCfxAddresses(chain string, addrs []string) ([]*cfxaddress.Address, []error) {
	var results []*cfxaddress.Address
	var errs []error

	for _, addr := range addrs {
		v, err := CheckCfxAddress(chain, addr)
		results = append(results, v)
		errs = append(errs, err)
	}
	return results, errs
}

func CheckCfxAddress(chain string, addr string) (*cfxaddress.Address, error) {
	chainType, chainId, err := ChainInfoByName(chain)
	if err != nil {
		return nil, err
	}
	if chainType != CHAIN_TYPE_CFX {
		return nil, errors.New("not cfx chain")
	}
	addrItem, err := cfxaddress.NewFromBase32(addr)
	if err != nil {
		return nil, err
	}
	if addrItem.GetNetworkID() != uint32(chainId) {
		return nil, fmt.Errorf("invalid conflux network address, want %v, got %v", uint32(chainId), addrItem.GetNetworkID())
	}
	return &addrItem, nil
}

func CurrentMonthStr() string {
	now := time.Now()
	return MonthOfTime(now)
}

func MonthOfTime(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}

func DateStrBeforeToday(day int) string {
	date := time.Now().Add(time.Hour * 24 * time.Duration(day) * -1)
	return fmt.Sprintf("%04d-%02d-%02d", date.Year(), date.Month(), date.Day())
}

func DateStrAfterToday(day int) string {
	date := time.Now().Add(time.Hour * 24 * time.Duration(day))
	return fmt.Sprintf("%04d-%02d-%02d", date.Year(), date.Month(), date.Day())
}

func DateStrOfTime(t time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

func TodayDateStr() string {
	now := time.Now()
	return fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())
}

func TodayBegin() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func TomorrowDateStr() string {
	t := time.Now().Add(time.Hour * 24)
	return fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

func TomorrowBegin() time.Time {
	t := time.Now().Add(time.Hour * 24)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// BeginningOfMonth return the begin of the month of t
func BeginningOfMonth(date time.Time) time.Time {
	year, month, _ := date.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, date.Location())
}

// BeginnigOfNextMonth return the end of the month of t
func BeginnigOfNextMonth(date time.Time) time.Time {
	return BeginningOfMonth(date).AddDate(0, 1, 0)
}

func EarlistDate() time.Time {
	return time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
}

func UintPtrToBig(val *uint) *big.Int {
	var result *big.Int
	if val != nil {
		result = big.NewInt(int64(*val))
	}
	return result
}

func Uint64PtrToBig(val *uint64) *big.Int {
	var result *big.Int
	if val != nil {
		result = new(big.Int).SetUint64(*val)
	}
	return result
}

func UintPtrToUint(val *uint) uint {
	result := uint(0)
	if val != nil {
		result = uint(*val)
	}
	return result
}

func Uint64Ptr(val uint64) *uint64 {
	return &val
}

func UintPtr(val uint) *uint {
	return &val
}

func MustParseStrToBig(s string) *big.Int {
	val, ok := new(big.Int).SetString(s, 0)
	if !ok {
		panic(fmt.Sprintf("failed to parse %s as big int", s))
	}
	return val
}

func ParseStrToBig(s string) (*big.Int, bool) {
	val, ok := new(big.Int).SetString(s, 0)
	return val, ok
}

func InUint256(val *big.Int) bool {
	return val.BitLen() <= 256
}

func BigIntToHexBig(val *big.Int) *hexutil.Big {
	return (*hexutil.Big)(val)
}

func MustNewBigIntByString(val string) *big.Int {
	b, _ := new(big.Int).SetString(val, 0)
	return b
}

func ParseCfxAddressesToCommon(addrs []string) ([]common.Address, error) {
	var results []common.Address
	for _, str := range addrs {
		v, err := cfxaddress.NewFromBase32(str)
		if err != nil {
			return nil, err
		}
		results = append(results, v.MustGetCommonAddress())
	}
	return results, nil
}

func MapHexAddressesToCfxAddressesStr(addrs []common.Address, networkId uint32) ([]string, error) {
	var results []string
	for _, str := range addrs {
		v, err := cfxaddress.NewFromCommon(str, networkId)
		if err != nil {
			return nil, err
		}
		results = append(results, v.String())
	}
	return results, nil
}

func CompactFormatTime(t time.Time) string {
	return t.Format("20060102150405")
}

func RandomNumber(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func Retry(count int, interval time.Duration, fn func() error) error {
	for i := 0; i < count; i++ {
		err := fn()
		if err != nil {
			logrus.WithError(err).WithField("stack", string(debug.Stack())).WithField("retry cnt", i).Info("retry function error")
			if i == count-1 {
				return err
			} else {
				time.Sleep(interval)
				continue
			}
		}
		return nil
	}
	return nil
}

func OneWayPassword(password string) string {
	return crypto.Keccak256Hash([]byte(password)).String()
}

func LeftPadZero(input string, l int) string {
	inputLen := len(input)
	if inputLen >= l {
		return input
	}

	for i := 0; i < l-inputLen; i++ {
		input = "0" + input
	}
	return input
}

func MustNewDripFromString(prettyValue string) *unit.Drip {
	v, err := unit.NewDripFromString(prettyValue)
	if err != nil {
		panic(err)
	}
	return v
}

func MapSlice[T any, R any](items []T, mapFunc func(item T) (R, error)) ([]R, error) {
	var result []R
	for _, item := range items {
		r, err := mapFunc(item)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

func MustMapSlice[T any, R any](items []T, mapFunc func(item T) R) []R {
	fn := func(item T) (R, error) {
		return mapFunc(item), nil
	}
	return Must(MapSlice(items, fn))
}

func GetMapKeys[KT comparable, VT any](value map[KT]VT) []KT {
	keys := make([]KT, 0, len(value))
	for k := range value {
		keys = append(keys, k)
	}
	return keys
}

func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}
