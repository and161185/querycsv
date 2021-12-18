package csvreader

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"os"
	"regexp"
	"sync"

	"querycsv/csvreader/clausechecker"
	"strings"

	"github.com/sirupsen/logrus"
)

type reader struct {
	log           *logrus.Logger
	filename      string
	query         string
	clause        string
	columns       map[string]int
	clauseColumns map[string]int
	res           chan *[]string
	fileIsRead    chan bool
	wg            sync.WaitGroup
	result        []map[string]string
	checker       checkerI
}

type checkerI interface {
	Check(params *map[string]string) bool
}

func NewReader(filename string, log *logrus.Logger, query string) *reader {
	r := reader{
		filename:      filename,
		log:           log,
		query:         query,
		clause:        query,
		columns:       make(map[string]int),
		clauseColumns: make(map[string]int),
		fileIsRead:    make(chan bool),
		res:           make(chan *[]string),
		checker:       clausechecker.NewChecker(log, query),
	}

	r.getFileColumns()
	r.getClauseColumns()

	return &r
}

func (r *reader) getFileColumns() {
	f, err := os.Open(r.filename)
	if err != nil {
		log.Fatalf("Couldn't read file : %v", err)
	}
	defer f.Close()

	csvr := csv.NewReader(f)

	record, err := csvr.Read()
	if err != nil {
		log.Fatalf("Couldn't get columns : %v", err)
	}

	if len(record) == 0 {
		log.Fatalf("File must have at least 1 column!")
	}

	columnsMap := map[string]int{}
	for i := 0; i < len(record); i++ {
		record[i] = strings.Trim(record[i], " ")
		colname := strings.ToLower(record[i])
		if colname == "" {
			log.Fatalf("Wrong column's %d name!", i+1)
		}

		if _, ok := columnsMap[colname]; ok {
			log.Fatalf("Duplicate columns %s !", colname)
		}

		columnsMap[colname] = i
	}

	r.columns = columnsMap
}

func (r *reader) getClauseColumns() {
	rgx := regexp.MustCompile(`[^'"a-zA-Z0-9]+([^'"a-zA-Z]+)?`)

	elems := rgx.Split(strings.ToLower(r.clause), -1)

	clauseColumns := map[string]int{}
	columnsMap := r.columns
	for i := 0; i < len(elems); i++ {
		colname := elems[i]
		if colname == "" || wordIsReserved(colname) {
			continue
		}
		if string([]rune(colname)[0]) == "'" || string([]rune(colname)[0]) == "\"" {
			continue
		}

		if val, ok := columnsMap[colname]; !ok {
			log.Fatalf("There is no column %s in file!", colname)
		} else {
			clauseColumns[colname] = val
		}
	}

	r.clauseColumns = clauseColumns
}

func wordIsReserved(word string) bool {
	test := strings.ToLower(word)
	return test == "select" || test == "from" || test == "where" || test == "and" || test == "or" || test == "top"
}

func (r *reader) FindRows(ctx context.Context) []map[string]string {

	f, err := os.Open(r.filename)
	if err != nil {
		log.Fatalf("Couldn't read file : %v", err)
	}
	defer f.Close()

	csvr := csv.NewReader(f)

	go r.collect(ctx)

	row := 0
loop:
	for {
		row++

		record, err := csvr.Read()
		if row == 1 {
			continue
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Couldn't read row %d : %v", row, err)
		}

		select {
		case <-ctx.Done(): //Если контекст завершен - прекращаем выполнение
			r.log.Debug("context is done")
			break loop
		default:
			r.wg.Add(1)
			go r.checkRecord(record)
		}

	}

	r.wg.Wait()

	r.wg.Add(1)

	r.fileIsRead <- true
	r.wg.Wait()

	return r.result
}

func (r *reader) checkRecord(record []string) {
	defer func() {
		r.wg.Done()
		r.log.Debug("checkRecord defer")
	}()

	r.log.Debug("checkRecord start ", record[0])

	params := map[string]string{}
	for k, v := range r.clauseColumns {
		params[k] = strings.ToLower(record[v])
	}

	if r.checker.Check(&params) {
		r.log.Debug("checkRecord start ", record[0], " pass")
		r.res <- &record
	}

}

func (r *reader) collect(ctx context.Context) {
	defer r.wg.Done()
	fileIsRead := false
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case refRecord := <-r.res:
			record := (*refRecord)
			r.log.Debug("record ", record[0], " collected")
			resRow := map[string]string{}
			for k, v := range r.columns {

				resRow[k] = record[v]
			}
			r.result = append(r.result, resRow)
		case fileIsRead = <-r.fileIsRead:
			r.log.Debug("got  fileIsRead")
		default:
			if fileIsRead {
				r.log.Debug("end collecting")
				break loop
			}
		}
	}
}
