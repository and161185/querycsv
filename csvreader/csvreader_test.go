package csvreader

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestReadConfigFile(t *testing.T) {

	r := NewReader("..\\testing\\table.csv", logrus.New(), "Salary >= 9000	AND	Age < 40")
	expectedColumns := map[string]int{"id": 0, "name": 1, "age": 2, "salary": 3}
	expectedClauseColumns := map[string]int{"age": 2, "salary": 3}
	assert.Equal(t, expectedColumns, r.columns)
	assert.Equal(t, expectedClauseColumns, r.clauseColumns)
}

func TestFindRows(t *testing.T) {

	r := NewReader("..\\testing\\table.csv", logrus.New(), "Salary >= 9000	AND	Age < 40")
	r.log.SetLevel(logrus.DebugLevel)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(5)) //nolint
	r.FindRows(ctx)
	expectedRows := []map[string]string{
		{"id": "4", "name": "Ivan", "age": "28", "salary": "10000"},
		{"id": "5", "name": "Ann", "age": "29", "salary": "10000"},
	}

	result := []map[string]string{} //nolint
	if r.result[0]["id"] == "5" {
		result = []map[string]string{
			r.result[1], r.result[0],
		}
	} else {
		result = r.result
	}

	assert.Equal(t, expectedRows, result)
}
