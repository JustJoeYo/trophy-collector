package steamid

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const steam64Base uint64 = 76561197960265728

var (
	reSteam2  = regexp.MustCompile(`^STEAM_[0-5]:([01]):(\d+)$`)
	reSteam3  = regexp.MustCompile(`^\[U:1:(\d+)\]$`)
	reSteam64 = regexp.MustCompile(`^7656119\d{10}$`)
	reSteam32 = regexp.MustCompile(`^\d{1,10}$`)
)

func Resolve(input string) (uint32, error) {
	input = strings.TrimSpace(input)

	if m := reSteam2.FindStringSubmatch(input); m != nil {
		y, _ := strconv.ParseUint(m[1], 10, 64)
		z, _ := strconv.ParseUint(m[2], 10, 64)
		return uint32(z*2 + y), nil
	}

	if m := reSteam3.FindStringSubmatch(input); m != nil {
		v, _ := strconv.ParseUint(m[1], 10, 64)
		return uint32(v), nil
	}

	if reSteam64.MatchString(input) {
		v, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid steam64 id: %w", err)
		}
		if v < steam64Base {
			return 0, fmt.Errorf("steam64 id too small")
		}
		return uint32(v - steam64Base), nil
	}

	if reSteam32.MatchString(input) {
		v, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid steam32 id: %w", err)
		}
		if v > 0xFFFFFFFF {
			return 0, fmt.Errorf("steam32 id out of range")
		}
		return uint32(v), nil
	}

	return 0, fmt.Errorf("unrecognised steam id format: %q", input)
}

func ToSteam64(accountID uint32) uint64 {
	return steam64Base + uint64(accountID)
}
