package formatter

import (
	"time"
)

// func StringToMonth(input string) time.Month {
// 	s := strings.ToLower(input)
// 	switch (s) {
// 	case "januari":
// 		return time.January
// 	case "februari":
// 		return time.February
// 	case "maret":
// 		return time.March
// 	case 4:
// 		return time.April
// 	case 5:
// 		return time.May
// 	case 6:
// 		return time.June
// 	case 7:
// 		return time.July
// 	case 8:
// 		return time.August
// 	case 9:
// 		return time.September
// 	case 10:
// 		return time.October
// 	case 11:
// 		return time.November
// 	case 12:
// 		return time.December
// 	default:
// 		return time.January
// 	}
// }

func IntToMonth(input int64) time.Month {

	switch input {
	case 1:
		return time.January
	case 2:
		return time.February
	case 3:
		return time.March
	case 4:
		return time.April
	case 5:
		return time.May
	case 6:
		return time.June
	case 7:
		return time.July
	case 8:
		return time.August
	case 9:
		return time.September
	case 10:
		return time.October
	case 11:
		return time.November
	case 12:
		return time.December
	default:
		return time.January
	}
}
