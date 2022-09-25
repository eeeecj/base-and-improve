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

func dbOpen(filepath string) (*Table, error) {
	return &Table{}, nil
}

func printPrompt() {
	fmt.Print("db>")
}
func readInput() (string, error) {
	buf := bufio.NewReader(os.Stdin)
	data, err := buf.ReadBytes('\n')
	return string(data), err
}

// run main 主函数，这样写方便单元测试
func run() {
	table, err := dbOpen("./db.txt")
	if err != nil {
		panic(err)
	}
	for {
		printPrompt()
		// 语句解析
		inputBuffer, err := readInput()
		if err != nil {
			fmt.Println("read err", err)
		}
		// 特殊操作
		if len(inputBuffer) != 0 && inputBuffer[0] == '.' {
			switch doMetaCommand(inputBuffer, table) {
			case metaCommandSuccess:
				continue
			case metaCommandUnRecognizedCommand:
				fmt.Println("Unrecognized command", inputBuffer)
				continue
			}
		}
		// 普通操作 code Generator
		statement := Statement{}
		switch prepareStatement(inputBuffer, &statement) {
		case prepareSuccess:
			break
		case prepareUnrecognizedStatement:
			fmt.Println("Unrecognized keyword at start of ", inputBuffer)
			continue
		default:
			fmt.Println("invalid unput ", inputBuffer)
			continue
		}
		res := executeStatement(&statement, table)
		if res == ExecuteSuccess {
			fmt.Println("Exected")
			continue
		}
		if res == ExecuteTableFull {
			fmt.Printf("Error: Table full.\n")
			break
		}
		if res == EXECUTE_DUPLICATE_KEY {
			fmt.Printf("Error: Duplicate key.\n")
			break
		}
	}
}

type metaCommandType int32

const (
	metaCommandSuccess metaCommandType = iota
	metaCommandUnRecognizedCommand
)

func doMetaCommand(input string, table *Table) metaCommandType {
	if input == ".exit" {
		dbClose(table)
		os.Exit(0)
		return metaCommandSuccess
	}
	if input == ".btree" {
		fmt.Printf("Tree:\n")
		printLeafNode(getPage(table.pager, 0))
		return metaCommandSuccess
	}
	if input == ".constants" {
		fmt.Printf("Constants:\n")
		printConstants()
		return metaCommandSuccess
	}
	return metaCommandUnRecognizedCommand
}

type Statement struct {
	statementType StatementType
	rowToInsert   *Row
}
type Row struct {
	ID       int32
	UserName string
	Email    string
}
type StatementType int32

const (
	statementSelect StatementType = iota
	statementInsert
)

type PrepareType int32

const (
	prepareUnrecognizedStatement PrepareType = iota
	prepareUnrecognizedSynaErr
	prepareSuccess
)

func prepareStatement(input string, statement *Statement) PrepareType {
	if len(input) >= 6 && input[0:6] == "insert" {
		statement.statementType = statementInsert
		inputs := strings.Split(input, " ")
		if len(inputs) <= 1 {
			return prepareUnrecognizedStatement
		}
		id, err := strconv.ParseInt(inputs[1], 10, 64)
		if err != nil {
			return prepareUnrecognizedSynaErr
		}
		statement.rowToInsert.ID = int32(id)
		statement.rowToInsert.UserName = inputs[2]
		statement.rowToInsert.Email = inputs[3]
		return prepareSuccess
	}
	if len(input) >= 6 && input[0:6] == "select" {
		statement.statementType = statementSelect
		return prepareSuccess
	}
	return prepareUnrecognizedStatement
}

type ExecuteResult int32

const (
	ExecuteSuccess ExecuteResult = iota
	ExecuteTableFull
	EXECUTE_DUPLICATE_KEY
)

// executeStatement 实行sql语句 ，解析器解析程statement，将最终成为我们的虚拟机
func executeStatement(statement *Statement, table *Table) executeResult {
	switch statement.statementType {
	case statementInsert:
		return executeInsert(statement, table)
	case statementSelect:
		return executeSelect(statement, table)
	default:
		fmt.Println("unknown statement")
	}
	return ExecuteSuccess
}

const (
	ROW_SIZE      = 291
	ID_SIZE       = 4
	USERNAME_SIZE = 32
	EMAIL_SIZE    = 255
)

func Uint32ToBytes(id int32) []byte {
	x := uint32(id)
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, x)
	return buf.Bytes()
}

func BytesToInt32(buf []byte) int32 {
	bufs := bytes.NewBuffer(buf)
	var x int32
	binary.Read(bufs, binary.BigEndian, x)
	return x
}

func getUseFulByte(buf []byte) string {
	bufs := bytes.NewBuffer(buf)
	var x string
	binary.Read(bufs, binary.BigEndian, bufs)
	return x
}

// 将row序列化到指针，为标准写入磁盘做准备
func serializeRow(row *Row, destionaton unsafe.Pointer) {
	ids := Uint32ToBytes(row.ID)
	q := (*[ROW_SIZE]byte)(destionaton)
	copy(q[0:ID_SIZE], ids)
	copy(q[ID_SIZE+1:ID_SIZE+USERNAME_SIZE], row.UserName)
	copy(q[ID_SIZE+USERNAME_SIZE+1:ROW_SIZE], row.Email)
}

// deserializeRow 将文件内容序列化成数据库元数据
func deserializeRow(source unsafe.Pointer, rowDestination *Row) {
	ids := make([]byte, ID_SIZE, ID_SIZE)
	sourceByte := (*[ROW_SIZE]byte)(source)
	copy(ids[0:ID_SIZE], (*sourceByte)[0:ID_SIZE])
	rowDestination.ID = BytesToInt32(ids)
	userName := make([]byte, USERNAME_SIZE, USERNAME_SIZE)
	copy(userName[0:], (*sourceByte)[ID_SIZE+1:ID_SIZE+USERNAME_SIZE])
	realNameBytes := getUseFulByte(userName)
	rowDestination.UserName = (string)(realNameBytes)
	emailStoreByte := make([]byte, EMAIL_SIZE, EMAIL_SIZE)
	copy(emailStoreByte[0:], (*sourceByte)[1+ID_SIZE+USERNAME_SIZE:ROW_SIZE])
	emailByte := getUseFulByte(emailStoreByte)
	rowDestination.Email = (string)(emailByte)
}

// Pager 管理数据从磁盘到内存
type Pager struct {
	osfile     *os.File
	fileLength int64
	numPages   uint32
	pages      []unsafe.Pointer // 存储数据
}

// Table 数据库表
type Table struct {
	rootPageNum uint32
	pager       *Pager
}

// pagerFlush 这一页写入文件系统
func pagerFlush(pager *Pager, pageNum, realNum uint32) error {
	if pager.pages[pageNum] == nil {
		return fmt.Errorf("pagerFlush null page")
	}
	offset, err := pager.osfile.Seek(int64(pageNum*PageSize), io.SeekStart)
	if err != nil {
		return fmt.Errorf("seek %v", err)
	}
	if offset == -1 {
		return fmt.Errorf("offset %v", offset)
	}
	originByte := make([]byte, realNum)
	q := (*[PageSize]byte)(pager.pages[pageNum])
	copy(originByte[0:realNum], (*q)[0:realNum])
	// 写入到byte指针里面
	bytesWritten, err := pager.osfile.WriteAt(originByte, offset)
	if err != nil {
		return fmt.Errorf("write %v", err)
	}
	// 捞取byte数组到这一页中
	fmt.Println("already wittern", bytesWritten)
	return nil
}

func dbClose(table *Table) {
	for i := uint32(0); i < table.pager.numPages; i++ {
		if table.pager.pages[i] == nil {
			continue
		}
		pagerFlush(table.pager, i, PageSize)
	}
	defer table.pager.osfile.Close()
	// go语言自带gc
}

const (
	TABLE_MAX_PAGES = 1000000
	PageSize        = 1024 * 1024 * 8
)

func getPage(pager *Pager, pageNum uint32) unsafe.Pointer {
	if pageNum > TABLE_MAX_PAGES {
		fmt.Println("Tried to fetch page number out of bounds:", pageNum)
		os.Exit(0)
	}
	if pager.pages[pageNum] == nil {
		page := make([]byte, PageSize)
		numPage := uint32(pager.fileLength / PageSize) // 第几页
		if pager.fileLength%PageSize == 0 {
			numPage += 1
		}
		if pageNum <= numPage {
			curOffset := pageNum * PageSize
			// 偏移到下次可以读读未知
			curNum, err := pager.osfile.Seek(int64(curOffset), io.SeekStart)
			if err != nil {
				panic(err)
			}
			fmt.Println(curNum)
			// 读到偏移这一页到下一页，必须是真的有多少字符
			if _, err = pager.osfile.ReadAt(page, curNum); err != nil && err != io.EOF {
				panic(err)
			}
		}
		pager.pages[pageNum] = unsafe.Pointer(&page[0])
		if pageNum >= pager.numPages {
			pager.numPages = pageNum + 1
		}
	}
	return pager.pages[pageNum]
}

// 返回key的位置，如果key不存在，返回应该被插入的位置
func tableFind(table *Table, key uint32) *Cursor {
	rootPageNum := table.rootPageNum
	rootNode := getPage(table.pager, rootPageNum)
	// 没有找到匹配到
	if getNodeType(rootNode) == leafNode {
		return leafNodeFind(table, rootPageNum, key)
	} else {
		fmt.Printf("Need to implement searching an internal node\n")
		os.Exit(0)
	}
	return nil
}

func leafNodeFind(table *Table, pageNum uint32, key uint32) *Cursor {
	node := getPage(table.pager, pageNum)
	num_cells := *leafNodeNumCells(node)
	cur := &Cursor{
		table:   table,
		pageNum: pageNum,
	}
	// Binary search
	var min_index uint32
	var one_past_max_index = num_cells
	for one_past_max_index != min_index {
		index := (min_index + one_past_max_index) / 2
		key_at_index := *leafNodeKey(node, index)
		if key == key_at_index {
			cur.cellNum = index
			return cur
		}
		// 如果在小到一边，就将最大值变成当前索引
		if key < key_at_index {
			one_past_max_index = index
		} else {
			min_index = index + 1 // 选择左侧
		}
	}
	cur.cellNum = min_index
	return cur
}

// Cursor 光标
type Cursor struct {
	table      *Table
	pageNum    uint32 // 第几页
	cellNum    uint32 // 多少个数据单元
	endOfTable bool
}

func tableStart(table *Table) *Cursor {
	rootNode := getPage(table.pager, table.rootPageNum)
	numCells := *leafNodeNumCells(rootNode)
	return &Cursor{
		table:      table,
		pageNum:    table.rootPageNum,
		cellNum:    0,
		endOfTable: numCells == 0,
	}
}

func cursorAdvance(cursor *Cursor) {
	node := getPage(cursor.table.pager, cursor.pageNum)
	cursor.cellNum += 1
	if cursor.cellNum >= (*leafNodeNumCells(node)) {
		cursor.endOfTable = true
	}
}
