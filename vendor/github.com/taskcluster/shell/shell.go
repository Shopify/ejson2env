// Package shell provides safe escaping of command line tokens
package shell

import "strings"

const safeChars = `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_/:=-'`

func Escape(tokens ...string) (escaped string) {
	var escapedToken string
	escapedTokens := []string{}
	for _, j := range tokens {
		for _, k := range j {
			if !strings.ContainsRune(safeChars, k) {
				goto escape
			}
		}
		escapedToken = strings.Replace(j, `'`, `\'`, -1)
		escapedTokens = append(escapedTokens, escapedToken)
		continue
	escape:
		escapedToken = `'` + strings.Replace(j, `'`, `'\''`, -1) + `'`
		if strings.HasPrefix(escapedToken, `''`) {
			escapedToken = escapedToken[2:]
		}
		escapedToken = strings.Replace(escapedToken, `\'''`, `\'`, -1)
		escapedTokens = append(escapedTokens, escapedToken)
	}
	return strings.Join(escapedTokens, " ")
}
