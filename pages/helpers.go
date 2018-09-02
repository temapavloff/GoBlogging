package pages

import "bytes"

var l2l = map[string]string{
	"А": "A", "а": "a", "В": "V", "в": "v", "Г": "G", "г": "g", "Д": "D", "д": "d",
	"Е": "E", "е": "e", "Ё": "Yo", "ё": "yo", "Ж": "Zh", "ж": "zh", "З": "Z", "з": "z",
	"И": "I", "и": "i", "Й": "J", "й": "j", "К": "K", "к": "k", "Л": "L", "л": "l",
	"М": "M", "м": "m", "Н": "N", "н": "n", "О": "O", "о": "o", "П": "P", "п": "p",
	"Р": "R", "р": "r", "С": "S", "с": "s", "Т": "T", "т": "t", "У": "U", "у": "u",
	"Ф": "F", "ф": "f", "Х": "H", "х": "h", "Ц": "Ts", "ц": "ts", "Ч": "Ch", "ч": "ch",
	"Ш": "Sh", "ш": "sh", "Щ": "Sch", "щ": "sch", "Ъ": "", "ъ": "", "Ы": "Y", "ы": "y",
	"Ь": "", "ь": "", "Э": "E", "э": "e", "Ю": "Yu", "ю": "yu", "Я": "Ya", "я": "ya",
}

func translite(text string) string {
	input := bytes.NewBufferString(text)
	output := bytes.NewBuffer(nil)

	for {
		r, _, err := input.ReadRune()
		if err != nil {
			break
		}

		if engLet, has := l2l[string(r)]; has {
			output.WriteString(engLet)
		} else {
			output.WriteRune(r)
		}
	}

	return output.String()
}

func slug(text string) string {
	input := bytes.NewBufferString(translite(text))
	output := bytes.NewBuffer(nil)

	for {
		r, _, err := input.ReadRune()
		if err != nil {
			break
		}

		if r == ' ' {
			output.WriteRune('-')
		} else {
			output.WriteRune(r)
		}
	}

	return output.String()
}
