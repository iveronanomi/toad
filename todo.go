package todo

import (
	"errors"
	"log"
	"regexp"
	"strings"
)

// prepared regexp for a collecting entry fields
var todoRegexp *regexp.Regexp

// build in placeholders
const (
	username    = "{{username}}"
	project     = "{{project_name}}"
	title       = "{{title}}"
	description = "{{description}}"
)

// regexp for a placeholders
var placeholdersRE = map[string]string{
	username:    "(?P<uName>.+?)",
	project:     "(?P<pName>.+?)",
	title:       "(?P<iTitle>.*)",
	description: "(?P<iDescription>.*)",
}

var order []string

// PreparePattern regular expression for extraction to-do fields
func PreparePattern(pattern string) *regexp.Regexp {
	order = make([]string, 0, len(placeholdersRE))
	defer func() {
		if r := recover(); r != nil {
			log.Printf("combine todo regexp error, %v", r)
		}
	}()

	// Prepare pattern before reg exp execution
	pattern = strings.NewReplacer(".", `\.`, " ", `\s`, "(", `\(`, ")", `\)`, "|", `\|`).Replace(pattern)

	p := regexp.MustCompile(`{{2}(?:[\w]+?)}{2}`).
		ReplaceAllFunc([]byte(pattern), func(src []byte) []byte {
			if p, ok := placeholdersRE[string(src)]; ok {
				order = append(order, string(src))
				return []byte(p)
			}
			log.Printf("unknown placeholder `%s`", src)
			return src
		})
	log.Printf("combined reg exp string %s", string(p))
	p = []byte(strings.NewReplacer(`{`, `\{`, `}`, `\}`, `[`, `\[`, `]`, `\]`).Replace(string(p)))

	log.Printf("combined reg exp string %s", p)
	//log.Printf("%s", string(p)+`(\n|\r)?`)
	// `\bTODO\((.*?){2,}\)\.\{(.*?){2,}\}\s(.*?)\n|\r`

	todoRegexp = regexp.MustCompile(string(p))
	return todoRegexp
}

// Extract entry fields from given string
func Extract(from string) (*todo, bool, error) {
	if todoRegexp == nil {
		return nil, false, errors.New("combine regexp mut be invoked before extraction")
	}
	matches := todoRegexp.FindAllStringSubmatch(from, -1)

	log.Printf("line: %#v", from)
	log.Printf("fields order: %#v", order)
	log.Printf("matches found: %#v", matches)

	if len(matches) < 1 {
		return nil, false, nil
	}
	if len(matches[0]) < (len(order) + 1) { // `found count of fields` + `matched string` (1?)
		return nil, false, nil
	}

	result := &todo{}
	for i := range order {
		switch order[i] {
		case username:
			result.Username = matches[0][i+1]
		case description:
			result.IssueDescription = matches[0][i+1]
		case title:
			result.IssueTitle = matches[0][i+1]
		case project:
			result.ProjectName = matches[0][i+1]
		}
	}

	return result, true, nil
}

type todo struct {
	Username         string
	ProjectName      string
	IssueTitle       string
	IssueDescription string

	File string
	Line int
}
