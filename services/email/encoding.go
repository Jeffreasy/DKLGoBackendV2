package email

import (
	"bytes"
	"regexp"
	"strings"
)

var (
	// Voorgecompileerde regex patterns voor betere performance
	scriptStyleRegex = regexp.MustCompile(`(?i)<(script|style)[^>]*>[\s\S]*?</($1)>`)
	commentRegex     = regexp.MustCompile(`<!--[\s\S]*?-->`)
	tagRegex         = regexp.MustCompile(`<[^>]*>`)
	entityRegex      = regexp.MustCompile(`&[#a-zA-Z0-9]+;`)
	whitespaceRegex  = regexp.MustCompile(`[ \t]+`)

	// Common HTML entities map
	htmlEntities = map[string]string{
		"&nbsp;":  " ",
		"&amp;":   "&",
		"&lt;":    "<",
		"&gt;":    ">",
		"&quot;":  "\"",
		"&apos;":  "'",
		"&cent;":  "¢",
		"&pound;": "£",
		"&euro;":  "€",
		"&copy;":  "©",
		"&reg;":   "®",
		"&trade;": "™",
		"&#8216;": "'",
		"&#8217;": "'",
		"&#8220;": `"`,
		"&#8221;": `"`,
		"&#8230;": "...",
		"&bull;":  "•",
		"&ndash;": "–",
		"&mdash;": "—",
		"&lsquo;": "'",
		"&rsquo;": "'",
		"&ldquo;": `"`,
		"&rdquo;": `"`,
	}

	// Semantic replacements voor HTML elementen
	blockElements = map[string]string{
		"</p>":       "\n\n",
		"</div>":     "\n",
		"</tr>":      "\n",
		"</table>":   "\n\n",
		"</h1>":      "\n\n",
		"</h2>":      "\n\n",
		"</h3>":      "\n\n",
		"</h4>":      "\n\n",
		"</h5>":      "\n\n",
		"</h6>":      "\n\n",
		"</pre>":     "\n\n",
		"</form>":    "\n\n",
		"</ul>":      "\n\n",
		"</ol>":      "\n\n",
		"</li>":      "\n",
		"<br>":       "\n",
		"<br/>":      "\n",
		"<br />":     "\n",
		"</article>": "\n\n",
		"</section>": "\n\n",
	}

	// Inline element replacements voor opmaak behoud
	inlineElements = map[string]string{
		"</em>":     "_", // Markering voor emphasis
		"</i>":      "_",
		"</b>":      "*", // Markering voor bold
		"</strong>": "*",
	}
)

// ProcessHTML verwerkt HTML content naar leesbare platte tekst met behoud van basis opmaak
func (s *EmailService) ProcessHTML(html string) string {
	if strings.TrimSpace(html) == "" {
		return ""
	}

	// Maximum grootte check om oneindige loops te voorkomen
	if len(html) > 10*1024*1024 { // 10MB limit
		html = html[:10*1024*1024]
	}

	var buf bytes.Buffer
	buf.Grow(len(html)) // Pre-allocate buffer

	// 1. Verwijder scripts, styles en comments
	html = scriptStyleRegex.ReplaceAllString(html, "")
	html = commentRegex.ReplaceAllString(html, "")

	// 2. Vervang block elements met newlines
	for tag, replacement := range blockElements {
		html = strings.ReplaceAll(html, strings.ToLower(tag), replacement)
		html = strings.ReplaceAll(html, strings.ToUpper(tag), replacement)
	}

	// 3. Vervang inline elements met hun markeringen
	for tag, replacement := range inlineElements {
		html = strings.ReplaceAll(html, strings.ToLower(tag), replacement)
		html = strings.ReplaceAll(html, strings.ToUpper(tag), replacement)
	}

	// 4. Verwijder alle overige HTML tags, maar behoud hun inhoud
	html = tagRegex.ReplaceAllString(html, "")

	// 5. Decode HTML entities
	html = entityRegex.ReplaceAllStringFunc(html, func(entity string) string {
		if replacement, ok := htmlEntities[entity]; ok {
			return replacement
		}
		return entity
	})

	// 6. Clean up whitespace met behoud van opmaak
	lines := strings.Split(html, "\n")
	var lastLineWasEmpty bool

	for i, line := range lines {
		line = whitespaceRegex.ReplaceAllString(strings.TrimSpace(line), " ")

		if line == "" {
			if !lastLineWasEmpty {
				if i > 0 {
					buf.WriteString("\n")
				}
				lastLineWasEmpty = true
			}
			continue
		}

		if i > 0 && !lastLineWasEmpty {
			buf.WriteString(" ")
		}
		buf.WriteString(line)
		lastLineWasEmpty = false
	}

	return buf.String()
}
