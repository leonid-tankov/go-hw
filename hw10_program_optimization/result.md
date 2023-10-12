## Последовательность оптимизации

Все результаты команды
```bash
go test -v -count=1 -timeout=30s -tags bench .
```
находятся в папке results.

1) Исходный вариант: файл result-1.txt
```text
time used: 653.505478ms / 300ms
memory used: 308Mb / 30Mb
```
2) Вместо чтения всего контента сразу используем bufio.NewScanner, меняем функцию getUsers на следующее:
<details><summary>getUsers</summary>

```go
func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)
	var i int
	for scanner.Scan() {
		var user User
		err := json.Unmarshal(scanner.Bytes(), &user)
		if err != nil {
			continue
		}
		result[i] = user
		i++
	}
	return result, nil
}
```
</details>

результат: файл result-2.txt
```text
time used: 670.008246ms / 300ms
memory used: 175Mb / 30Mb
```
3) Переходим на fastjson, используя функцию GetString, вследствие чего структура Users больше не нужна, начинаем хранить только слайс строк (emails), код приобретает следующий вид:
<details><summary>with fastjson</summary>

```go
package hw10programoptimization

import (
	"bufio"
	"fmt"
	"github.com/valyala/fastjson"
	"io"
	"regexp"
	"strings"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getEmails(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type emails [100_000]string

func getEmails(r io.Reader) (result emails, err error) {
	scanner := bufio.NewScanner(r)
	var i int
	for scanner.Scan() {
		email := fastjson.GetString(scanner.Bytes(), "Email")
		result[i] = email
		i++
	}
	return result, nil
}

func countDomains(e emails, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, email := range e {
		matched, err := regexp.Match("\\."+domain, []byte(email))
		if err != nil {
			return nil, err
		}

		if matched {
			num := result[strings.ToLower(strings.SplitN(email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(email, "@", 2)[1])] = num
		}
	}
	return result, nil
}


```
</details>

результат: файл result-3.txt
```text
time used: 398.256503ms / 300ms
memory used: 136Mb / 30Mb
```
4) Вместо regexp.Match используем strings.HasSuffix, так как нас интересует только конец почты, добавляем проверку с strings.SplitN, чтобы проверить корректность почты.

Также данные строки можно объединить
```go
num := result[strings.ToLower(strings.SplitN(email, "@", 2)[1])]
num++
result[strings.ToLower(strings.SplitN(email, "@", 2)[1])] = num
```
Функция countDomains приобретает вид:
<details><summary>countDomains</summary>

```go
func countDomains(e emails, domain string) (DomainStat, error) {
    result := make(DomainStat)
    for _, email := range e {
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
```
</details>

результат: файл result-4.txt
```text
time used: 251.631352ms / 300ms
memory used: 5Mb / 30Mb
```
5) Убираем функционал getEmails в countDomains, поэтому можно не использовать слайс emails

результат: файл result-5.txt
```text
time used: 244.284411ms / 300ms
memory used: 3Mb / 30Mb
```
