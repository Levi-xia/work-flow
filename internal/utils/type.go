package utils

import "strconv"

// uint转string
func UintToString(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}

// uint转int64
func UintToInt64(u uint) int64 {
	return int64(u)
}

// int转string
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// string转int64
func StringToInt64(s string) (int64, error){
	i, err :=  strconv.ParseInt(s, 10, 64)
	return i, err
}

// int64转string
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// float64转string
func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// string转float64
func StringToFloat64(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	return f, err
}

// string转int
func StringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	return i, err
}