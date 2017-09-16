package utils

import "testing"

func TestRandString(t *testing.T) {
	str := RandString(32)
	if len(str) != 32 {
		t.Fatalf("长度不匹配,目标长度32, 实际%d-%s", len(str), str)
	}
}

func TestMd5(t *testing.T) {
	str := Md5("123456")
	if len(str) != 32 {
		t.Fatalf("长度不匹配,目标长度32, 实际%d-%s", len(str), str)
	}
}

func TestRandNumber(t *testing.T) {
	num := RandNumber(10000)
	if num <= 0 && num >= 10000 {
		t.Fatalf("随机数不在有效范围内-%d", num)
	}
}
