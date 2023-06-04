package parseutils

import (
	"errors"
)

func ParseUintFromPrefix(str string) (int, error, int) {
	ptr := 0
	for ptr < len(str) && str[ptr] == ' ' {
		ptr++
	} 

	if ptr == len(str) || !('0' <= str[ptr] && str[ptr] <= '9') {
		return 0, errors.New("str prefix is not valid number"), -1
	}

	res := 0
	for ptr < len(str) && ('0' <= str[ptr] && str[ptr] <= '9') {
		res = res * 10 + int(str[ptr] - '0')
		ptr++
	}

	return res, nil, ptr
}
