package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestTommowDateStr(t *testing.T) {
	v := TomorrowDateStr()
	fmt.Println(v)

	tomorrow, _ := time.Parse("2006-01-02", TomorrowDateStr())
	fmt.Println(tomorrow)
}

func TestCompactFormatTime(t *testing.T) {
	v := CompactFormatTime(time.Now())
	fmt.Println(v)
}

func TestRandomNumber(t *testing.T) {
	fmt.Println(RandomNumber(100, 200))
}

func TestLeftPadZero(t *testing.T) {
	fmt.Println(LeftPadZero("123", 10))
}
