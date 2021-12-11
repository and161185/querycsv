package clausechecker

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestReadConfigFile(t *testing.T) {

	c := &checker{}
	c.log = logrus.New()

	parameters := make(map[string]string)
	parameters["age"] = "30"
	parameters["salary"] = "10000"

	c.expression = c.normalizeExpression("Salary >= 9000	AND	Age < 40")
	result := c.Check(&parameters)

	assert.Equal(t, true, result)

	c.expression = c.normalizeExpression("Salary < 9000 OR Age > 40")
	result = c.Check(&parameters)

	assert.Equal(t, false, result)
}
