package iputil

import (
	"errors"
	"fmt"
	"math/big"
	"net"
	"strconv"
)

var ErrRatioInvalid = errors.New("ratio is invalid")

func GetIPRangeFromNetEnd(
	ipnet net.IPNet,
	ratio float64,
	offsetFromEnd int64,
) (net.IP, net.IP, error) {
	if ratio <= 0.0 || ratio >= 1.0 {
		return nil, nil, fmt.Errorf(
			"ratio %q: %w",
			strconv.FormatFloat(ratio, 'g', 2, 64),
			ErrRatioInvalid,
		)
	}

	ones, bits := ipnet.Mask.Size()
	hostBits := bits - ones
	totalIPs := new(big.Int).Lsh(big.NewInt(1), uint(hostBits))

	rangeTotalF := new(
		big.Float,
	).Mul(new(big.Float).SetInt(totalIPs), new(big.Float).SetFloat64(ratio))
	rangeTotal, _ := rangeTotalF.Int(nil)
	if rangeTotal.Cmp(big.NewInt(0)) <= 0 {
		return nil, nil, fmt.Errorf(
			"ratio %q too small: %w",
			strconv.FormatFloat(ratio, 'g', 2, 64),
			ErrRatioInvalid,
		)
	}

	IPNetworkStart := new(big.Int).SetBytes(ipnet.IP.Mask(ipnet.Mask))
	IPNetworkEnd := new(
		big.Int,
	).Add(IPNetworkStart, new(big.Int).Sub(totalIPs, big.NewInt(offsetFromEnd)))
	IPRangeFirst := new(big.Int).Sub(IPNetworkEnd, rangeTotal)
	IPRangeLast := new(big.Int).Sub(IPNetworkEnd, big.NewInt(1))

	return net.IP(IPRangeFirst.Bytes()), net.IP(IPRangeLast.Bytes()), nil
}
