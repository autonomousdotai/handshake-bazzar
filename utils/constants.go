package utils

import (
	"errors"
	"strings"
)

var DEFAULT_PAGE_SIZE = "20"
var DEFAULT_PAGE = "1"

var ORDER_STATUS_NEW = 0
var ORDER_STATUS_SHAKED = 1
var ORDER_STATUS_SHAKED_PROCESS = -11
var ORDER_STATUS_SHAKED_FAILED = -21
var ORDER_STATUS_DELIVERED = 3
var ORDER_STATUS_DELIVERED_PROCESS = 13
var ORDER_STATUS_DELIVERED_FAILED = 23
var ORDER_STATUS_CANCELED = 4
var ORDER_STATUS_CANCELED_PROCESS = 14
var ORDER_STATUS_CANCELED_FAILED = 24
var ORDER_STATUS_ACCEPTED = 5
var ORDER_STATUS_ACCEPTED_PROCESS = 15
var ORDER_STATUS_ACCEPTED_FAILED = 25
var ORDER_STATUS_REJECTED = 6
var ORDER_STATUS_REJECTED_PROCESS = 16
var ORDER_STATUS_REJECTED_FAILED = 26
var ORDER_STATUS_WITHDRAWED = 7
var ORDER_STATUS_WITHDRAWED_PROCESS = 17
var ORDER_STATUS_WITHDRAWED_FAILED = 27

var OFFCHAIN_BAZZAR = "bz"
var OFFCHAIN_BAZZAR_SHAKE = "bzs"
var OFFCHAIN_USER = "usr"

func ParseOffchain(offchain string) (string, string, error) {
	offchains := strings.Split(offchain, "_")
	if len(offchains) >= 2 {
		return strings.Trim(offchains[0], " "), strings.Trim(offchains[1], " "), nil
	}
	return "", "", errors.New("offchain is invalid")
}
