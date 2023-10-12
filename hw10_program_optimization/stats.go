package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	result := make(DomainStat)
	for scanner.Scan() {
		email := fastjson.GetString(scanner.Bytes(), "Email")
		if strings.HasSuffix(email, domain) {
			fields := strings.SplitN(email, "@", 2)
			if len(fields) != 2 {
				return result, fmt.Errorf("invalid email: %s", email)
			}
			result[strings.ToLower(fields[1])]++
		}
	}

	return result, nil
}
