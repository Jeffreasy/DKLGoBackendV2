package email

import (
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
)

func (s *EmailService) extractCharset(contentType string) string {
	if strings.Contains(strings.ToLower(contentType), "charset") {
		parts := strings.Split(contentType, ";")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(strings.ToLower(part), "charset=") {
				charset := strings.TrimPrefix(part, "charset=")
				charset = strings.Trim(charset, `"'`)
				return charset
			}
		}
	}
	return ""
}

func (s *EmailService) decodeCharset(content, charset string) (string, error) {
	charset = strings.ToLower(charset)
	var decoder *encoding.Decoder

	switch charset {
	case "utf-8", "us-ascii":
		return content, nil
	case "iso-8859-1":
		decoder = charmap.ISO8859_1.NewDecoder()
	case "iso-8859-2":
		decoder = charmap.ISO8859_2.NewDecoder()
	case "iso-8859-3":
		decoder = charmap.ISO8859_3.NewDecoder()
	case "iso-8859-4":
		decoder = charmap.ISO8859_4.NewDecoder()
	case "iso-8859-5":
		decoder = charmap.ISO8859_5.NewDecoder()
	case "iso-8859-6":
		decoder = charmap.ISO8859_6.NewDecoder()
	case "iso-8859-7":
		decoder = charmap.ISO8859_7.NewDecoder()
	case "iso-8859-8":
		decoder = charmap.ISO8859_8.NewDecoder()
	case "iso-8859-9":
		decoder = charmap.ISO8859_9.NewDecoder()
	case "iso-8859-10":
		decoder = charmap.ISO8859_10.NewDecoder()
	case "iso-8859-13":
		decoder = charmap.ISO8859_13.NewDecoder()
	case "iso-8859-14":
		decoder = charmap.ISO8859_14.NewDecoder()
	case "iso-8859-15":
		decoder = charmap.ISO8859_15.NewDecoder()
	case "iso-8859-16":
		decoder = charmap.ISO8859_16.NewDecoder()
	case "windows-1250":
		decoder = charmap.Windows1250.NewDecoder()
	case "windows-1251":
		decoder = charmap.Windows1251.NewDecoder()
	case "windows-1252":
		decoder = charmap.Windows1252.NewDecoder()
	case "windows-1253":
		decoder = charmap.Windows1253.NewDecoder()
	case "windows-1254":
		decoder = charmap.Windows1254.NewDecoder()
	case "windows-1255":
		decoder = charmap.Windows1255.NewDecoder()
	case "windows-1256":
		decoder = charmap.Windows1256.NewDecoder()
	case "windows-1257":
		decoder = charmap.Windows1257.NewDecoder()
	case "windows-1258":
		decoder = charmap.Windows1258.NewDecoder()
	case "koi8r":
		decoder = charmap.KOI8R.NewDecoder()
	case "koi8u":
		decoder = charmap.KOI8U.NewDecoder()
	case "gbk", "gb2312":
		decoder = simplifiedchinese.GBK.NewDecoder()
	case "gb18030":
		decoder = simplifiedchinese.GB18030.NewDecoder()
	case "big5":
		decoder = traditionalchinese.Big5.NewDecoder()
	case "euc-jp":
		decoder = japanese.EUCJP.NewDecoder()
	case "iso-2022-jp":
		decoder = japanese.ISO2022JP.NewDecoder()
	case "shift-jis":
		decoder = japanese.ShiftJIS.NewDecoder()
	case "euc-kr":
		decoder = korean.EUCKR.NewDecoder()
	case "utf-16be":
		decoder = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder()
	case "utf-16le":
		decoder = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	default:
		return content, nil
	}

	result, err := decoder.String(content)
	if err != nil {
		return content, err
	}
	return result, nil
}
