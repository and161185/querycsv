package clausechecker

import (
	"regexp"
	"strings"

	"github.com/PaesslerAG/gval"
	"github.com/sirupsen/logrus"
)

type checker struct {
	log        *logrus.Logger
	expression string
}

func NewChecker(log *logrus.Logger, expression string) *checker {
	c := &checker{
		log: log,
	}

	c.expression = c.normalizeExpression(expression)
	return c
}

func (c *checker) Check(params *map[string]string) bool {

	value, err := gval.Evaluate(c.expression, (*params))
	if err != nil {
		c.log.Fatalf("Can't evaluate %s , got %v !", c.expression, err)
	}

	return value == true

}

func (c *checker) normalizeExpression(expression string) string {

	rgx := regexp.MustCompile(`\s`)

	elems := rgx.Split(expression, -1)
	string := strings.Join(elems, " ")
	string = strings.ToLower(string)
	string = strings.ReplaceAll(string, " and ", " && ")
	string = strings.ReplaceAll(string, " or ", " || ")

	return string
}
