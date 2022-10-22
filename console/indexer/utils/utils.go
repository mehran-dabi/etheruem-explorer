package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

func ParseScanQuery(s string, latestBlock int64) (start int64, end int64, err error) {
	parts := strings.Split(s, ":")

	switch {
	case len(s) == 0:
		start = 0
		end = latestBlock
		return
	case len(parts) == 1:
		start, err = strconv.ParseInt(parts[0], 10, 0)
		if err != nil {
			return
		}

		end = latestBlock
		return

	case len(parts) == 2:
		start, err = strconv.ParseInt(parts[0], 10, 0)
		if err != nil {
			return
		}

		end, err = strconv.ParseInt(parts[1], 10, 0)
		if err != nil {
			return
		}

		return
	default:
		err = fmt.Errorf("range must be in format start:end")
		return
	}
}

func AddrToHex(addr *common.Address) string {
	const ZeroAddress = `0x0000000000000000000000000000000000000000`
	if addr == nil {
		return ZeroAddress
	}
	return addr.Hex()
}
