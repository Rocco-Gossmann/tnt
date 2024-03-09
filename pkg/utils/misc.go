package utils

const SQL_DATETIME_FORMAT = "2006-01-02 15:04:05"

func Suffix(cnt int, singular string, plural string) string {
	ret := plural
	if cnt == 1 || cnt == -1 {
		ret = singular
	}

	return ret
}
