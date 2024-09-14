package handlers

import "strconv"

func getNthDigit(fourDigitInteger int, nth int) int64 {
	strFourDigitInteger := strconv.Itoa(fourDigitInteger)
	if nth < 0 || nth >= len(strFourDigitInteger) {
		return -1
	}
	nthDigitChar := strFourDigitInteger[nth]
	nthDigitInt := int(nthDigitChar - '0')
	return int64(nthDigitInt)
}

func permissionIntToString(permissionInt int64) string {
	switch permissionInt {
	case 0:
		return "No Access"
	case 1:
		return "View Only"
	case 2:
		return "Edit Only"
	case 3:
		return "Edit and Approve"
	default:
		return "Invalid Permission"
	}
}
