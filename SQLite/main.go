package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

type metaCommandType int32

const (
	metaCommandSuccess metaCommandType = iota
	metaCommandUnRecognizedCommand
)

type StatementType int32

const (
	STATEMENT_SELECT StatementType = iota
	STATEMENT_INSERT
)

type Statement struct {
	statementType StatementType
	Row           *Row
}

type PrepareType int32

const (
	PREPARE_SUCCESS PrepareType = iota
	PREPARE_UNRECOGNIZED_STATEMENT
	PREPARE_NEGATIVE_ID
)

type ExecuteResultType int32

const (
	EXECUTE_SUCCESS ExecuteResultType = iota
	EXECUTE_TABLE_FULL
)

const (
	ID_SIZE         = 4
	USERNAME_SIZE   = 32
	EMAIL_SIZE      = 255
	ID_OFFSET       = 0
	USERNAME_OFFSET = ID_SIZE + ID_OFFSET
	EMAIL_OFFSET    = USERNAME_OFFSET + USERNAME_SIZE
	ROW_SIZE        = ID_SIZE + USERNAME_SIZE + EMAIL_SIZE
)

type Row struct {
	ID       int32
	UserName string
	Email    string
}

const (
	PAGE_SIZE       = 4094
	TABLE_MAX_PAGES = 100
	ROW_PER_PAGE    = PAGE_SIZE / ROW_SIZE
	TABLE_MAX_ROWS  = ROW_PER_PAGE * TABLE_MAX_PAGES
)

type Pager struct {
	fs         *os.File
	fileLength int64
	pages      []unsafe.Pointer
}

type Table struct {
	numRows int
	Pager   *Pager
}

type Cursor struct {
	table    *Table
	numRows  int
	tableEnd bool
}

func tableStart(table *Table) *Cursor {
	cursor := &Cursor{}
	cursor.table = table
	cursor.numRows = 0
	cursor.tableEnd = table.numRows == 0
	return cursor
}

func tableEnd(table *Table) *Cursor {
	cursor := &Cursor{}
	cursor.table = table
	cursor.numRows = table.numRows
	cursor.tableEnd = true
	return cursor
}

func cursorAdvance(cursor *Cursor) {
	cursor.numRows += 1
	if cursor.numRows >= cursor.table.numRows {
		cursor.tableEnd = true
	}
}
func pageOpen(filename string) (*Pager, error) {
	fs, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	fileLength, _ := fs.Seek(0, io.SeekEnd)
	pager := &Pager{
		fs:         fs,
		fileLength: fileLength,
		pages:      make([]unsafe.Pointer, TABLE_MAX_PAGES),
	}
	return pager, nil
}

func dbOpen(filename string) *Table {
	table := &Table{}
	table.Pager, _ = pageOpen(filename)
	table.numRows = int(table.Pager.fileLength / ROW_SIZE)
	return table
}

func dbClose(table *Table) {
	fullPages := table.numRows / ROW_PER_PAGE
	for i := 0; i < fullPages; i++ {
		if table.Pager.pages[i] == nil {
			continue
		}
		pageFlush(table.Pager, i, PAGE_SIZE)
	}
	numAdditionalRows := table.numRows % ROW_PER_PAGE
	if numAdditionalRows > 0 {
		pageNum := table.numRows / ROW_PER_PAGE
		if table.Pager.pages[pageNum] != nil {
			pageFlush(table.Pager, pageNum, numAdditionalRows*ROW_SIZE)
		}
	}
	defer table.Pager.fs.Close()

}

// pagerFlush 这一页写入文件系统
func pageFlush(pager *Pager, pageNum, realNum int) error {
	if pager.pages[pageNum] == nil {
		return fmt.Errorf("pagerFlush null page")
	}
	offset, err := pager.fs.Seek(int64(pageNum*PAGE_SIZE), io.SeekStart)
	if err != nil {
		return fmt.Errorf("seek %v", err)
	}
	if offset == -1 {
		return fmt.Errorf("offset %v", offset)
	}
	buf := *(*[]byte)(pager.pages[pageNum])
	bytesWritten, err := pager.fs.WriteAt(buf[:realNum], offset)
	if err != nil {
		return fmt.Errorf("write %v", err)
	}
	// 捞取byte数组到这一页中
	fmt.Println("already wittern", bytesWritten)
	return nil
}

func main() {
	table := dbOpen("./1.db")
	reader := bufio.NewReader(os.Stdin)
	for {
		printPrompt()
		input, err := readInput(reader)
		if err != nil {
			fmt.Println("read err", err)
		}
		if len(input) != 0 && input[0] == '.' {
			switch doMetaCommand(input, table) {
			case metaCommandSuccess:
				continue
			case metaCommandUnRecognizedCommand:
				fmt.Println("Unrecognized command", input)
				continue
			}
		}
		var statement = Statement{Row: &Row{}}
		switch prepareStatement(input, &statement) {
		case PREPARE_SUCCESS:
		case PREPARE_UNRECOGNIZED_STATEMENT:
			fmt.Println("Unrecognized keyword at start of ", input)
			continue
		case PREPARE_NEGATIVE_ID:
			fmt.Println("ID must be positive.")
			continue
		default:
			fmt.Println("invalid input ", input)
			continue
		}
		switch executeStatement(&statement, table) {
		case EXECUTE_SUCCESS:
			fmt.Println("Executed.")
		case EXECUTE_TABLE_FULL:
			fmt.Println("Err:Table full.")
		}
	}
}

func getPage(pager *Pager, pageNum int) unsafe.Pointer {
	if pageNum > TABLE_MAX_PAGES {
		fmt.Printf("Tried to fetch page number out of bounds. %d > %d", pageNum, TABLE_MAX_PAGES)
		os.Exit(0)
	}
	if pager.pages[pageNum] == nil {
		page := make([]byte, PAGE_SIZE)
		numPages := int(pager.fileLength / PAGE_SIZE)
		if pager.fileLength%PAGE_SIZE == 0 {
			numPages += 1
		}
		if pageNum <= numPages {
			offset := numPages * PAGE_SIZE
			curNum, err := pager.fs.Seek(int64(offset), io.SeekStart)
			if err != nil {
				panic(err)
			}
			if _, err = pager.fs.ReadAt(page, curNum); err != nil && err != io.EOF {
				panic(err)
			}
		}
		pager.pages[pageNum] = unsafe.Pointer(&page)
	}
	return pager.pages[pageNum]
}

func cursorValue(cursor *Cursor) unsafe.Pointer {
	num := cursor.numRows
	pageNum := num / ROW_PER_PAGE
	page := getPage(cursor.table.Pager, pageNum)
	rowOffset := num % ROW_PER_PAGE
	byteOffset := rowOffset * ROW_SIZE
	buf := *(*[]byte)(page)
	point := unsafe.Pointer(&buf[byteOffset])
	return point
}

func Uint32ToBytes(id int32) []byte {
	x := uint32(id)
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, x)
	return buf.Bytes()
}

func BytesToInt32(buf []byte) int32 {
	bufs := bytes.NewBuffer(buf)
	var x int32
	binary.Read(bufs, binary.BigEndian, &x)
	return x
}
func serializeRow(source *Row, destination unsafe.Pointer) {
	id := Uint32ToBytes(source.ID)
	q := (*[ROW_SIZE]byte)(destination)
	copy(q[ID_OFFSET:ID_SIZE], id)
	copy(q[USERNAME_OFFSET:USERNAME_OFFSET+USERNAME_SIZE], source.UserName)
	copy(q[EMAIL_SIZE:ROW_SIZE], source.Email)
}

func deserializeRow(source unsafe.Pointer, destination *Row) {
	id := make([]byte, ID_SIZE, ID_SIZE)
	sourceBuf := (*[ROW_SIZE]byte)(source)
	copy(id, sourceBuf[ID_OFFSET:ID_SIZE])
	destination.ID = BytesToInt32(id)
	username := make([]byte, USERNAME_SIZE, USERNAME_SIZE)
	copy(username, sourceBuf[USERNAME_OFFSET:USERNAME_OFFSET+USERNAME_SIZE])
	destination.UserName = string(username)
	email := make([]byte, EMAIL_SIZE, EMAIL_SIZE)
	copy(email, sourceBuf[EMAIL_OFFSET:EMAIL_OFFSET+EMAIL_SIZE])
	destination.Email = string(email)
}

func prepareStatement(input string, statement *Statement) PrepareType {
	if len(input) > 6 && input[:6] == "insert" {
		statement.statementType = STATEMENT_INSERT
		inputs := strings.Split(input, " ")
		if len(inputs) < 1 {
			return PREPARE_UNRECOGNIZED_STATEMENT
		}
		id, err := strconv.ParseInt(inputs[1], 10, 32)
		if err != nil {
			return PREPARE_UNRECOGNIZED_STATEMENT
		}
		if id < 0 {
			return PREPARE_NEGATIVE_ID
		}
		statement.Row.ID = int32(id)
		statement.Row.UserName = inputs[2]
		statement.Row.Email = inputs[3]
		return PREPARE_SUCCESS
	}
	if len(input) >= 6 && input[:6] == "select" {
		statement.statementType = STATEMENT_SELECT
		return PREPARE_SUCCESS
	}
	return PREPARE_UNRECOGNIZED_STATEMENT
}

func executeStatement(statement *Statement, table *Table) ExecuteResultType {
	switch statement.statementType {
	case STATEMENT_SELECT:
		return executeSelect(statement, table)
	case STATEMENT_INSERT:
		return executeInsert(statement, table)
	}
	return EXECUTE_SUCCESS
}

func executeSelect(statement *Statement, table *Table) ExecuteResultType {
	var row Row
	cursor := tableStart(table)
	for !cursor.tableEnd {
		deserializeRow(cursorValue(cursor), &row)
		printRow(&row)
		cursorAdvance(cursor)
	}
	return EXECUTE_SUCCESS
}

func executeInsert(statement *Statement, table *Table) ExecuteResultType {
	if table.numRows >= TABLE_MAX_ROWS {
		return EXECUTE_TABLE_FULL
	}
	row := statement.Row
	cursor := tableEnd(table)
	serializeRow(row, cursorValue(cursor))
	table.numRows++
	return EXECUTE_SUCCESS
}
func doMetaCommand(input string, table *Table) metaCommandType {
	if input == ".exit" {
		dbClose(table)
		os.Exit(1)
		return metaCommandSuccess
	}
	if input == ".btree" {
		fmt.Printf("Tree:\n")
		//printLeafNode(getPage(table.pager, 0))
		return metaCommandSuccess
	}
	if input == ".constants" {
		fmt.Printf("Constants:\n")
		//printConstants()
		return metaCommandSuccess
	}
	return metaCommandUnRecognizedCommand
}

func readInput(reader *bufio.Reader) (string, error) {
	buf := bufio.NewReader(os.Stdin)
	data, err := buf.ReadBytes('\n')
	data = data[:len(data)-2]
	return string(data), err
}
func printPrompt() {
	fmt.Print("db>")
}

func printRow(row *Row) {
	fmt.Printf("(%d,%s,%s)\n", row.ID, row.UserName, row.Email)
}