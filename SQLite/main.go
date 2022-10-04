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

// 定义处理结果状态
type metaCommandType int32

const (
	metaCommandSuccess metaCommandType = iota
	metaCommandUnRecognizedCommand
)

// 定义语句类型
type StatementType int32

const (
	STATEMENT_SELECT StatementType = iota
	STATEMENT_INSERT
)

type Statement struct {
	statementType StatementType
	Row           *Row
}

// 定义预处理结果类型
type PrepareType int32

const (
	PREPARE_SUCCESS PrepareType = iota
	PREPARE_UNRECOGNIZED_STATEMENT
	PREPARE_NEGATIVE_ID
)

// 定义执行结果类型
type ExecuteResultType int32

const (
	EXECUTE_SUCCESS ExecuteResultType = iota
	EXECUTE_TABLE_FULL
	EXECUTE_DUPLICATE_KEY
)

// 定义表格数据大小和偏移量
const (
	ID_SIZE         = 4
	USERNAME_SIZE   = 32
	EMAIL_SIZE      = 255
	ID_OFFSET       = 0
	USERNAME_OFFSET = ID_SIZE + ID_OFFSET
	EMAIL_OFFSET    = USERNAME_OFFSET + USERNAME_SIZE
	ROW_SIZE        = ID_SIZE + USERNAME_SIZE + EMAIL_SIZE
)

// 定义表格结构
type Row struct {
	ID       int32
	UserName string
	Email    string
}

// 定义页表大小
const (
	PAGE_SIZE       = 4094
	TABLE_MAX_PAGES = 100
)

// 存储数据的页表项
type Pager struct {
	fs         *os.File
	fileLength int64
	pages      []unsafe.Pointer
	numPages   int
}

// 表格结构
type Table struct {
	rootPageNum int
	Pager       *Pager
}

type Cursor struct {
	table    *Table
	pageNum  int
	cellNum  int
	tableEnd bool
}

// 打印常量
func printConstants() {
	fmt.Printf("ROW_SIZE: %d\n", ROW_SIZE)
	fmt.Printf("COMMON_NODE_HEADER_SIZE: %d\n", COMMON_NODE_HEADER_SIZE)
	fmt.Printf("LEAF_NODE_HEADER_SIZE: %d\n", LEAF_NODE_HEADER_SIZE)
	fmt.Printf("LEAF_NODE_CELL_SIZE: %d\n", LEAF_NODE_CELL_SIZE)
	fmt.Printf("LEAF_NODE_SPACE_FOR_CELLS: %d\n", LEAF_NODE_SPACE_FOR_CELLS)
	fmt.Printf("LEAF_NODE_MAX_CELLS: %d\n", LEAF_NODE_MAX_CELLS)
}

// 获取存储对应页面的指针
func getPage(pager *Pager, pageNum int) unsafe.Pointer {
	if pageNum > TABLE_MAX_PAGES {
		fmt.Printf("Tried to fetch page number out of bounds. %d > %d", pageNum, TABLE_MAX_PAGES)
		os.Exit(0)
	}
	//如果页面不存在
	if pager.pages[pageNum] == nil {
		//使用byte存储值
		page := make([]byte, PAGE_SIZE)
		//获取当前数据库中的数据页面数量
		numPages := int(pager.fileLength / PAGE_SIZE)
		if pager.fileLength%PAGE_SIZE == 0 {
			numPages += 1
		}
		//当需要读取的页面在数据库中存在时
		if pageNum <= numPages {
			offset := pageNum * PAGE_SIZE
			//将文件的读标偏移
			curNum, err := pager.fs.Seek(int64(offset), io.SeekStart)
			if err != nil {
				panic(err)
			}
			//将文件中的数据读取到页面中
			if _, err = pager.fs.ReadAt(page, curNum); err != nil && err != io.EOF {
				panic(err)
			}
		}
		//更新页面指针和页面数量
		pager.pages[pageNum] = unsafe.Pointer(&page[0])
		if pageNum >= pager.numPages {
			pager.numPages = pageNum + 1
		}
	}
	return pager.pages[pageNum]
}

// 根据level打印
func indent(level int) {
	for i := 0; i < level; i++ {
		fmt.Print(" ")
	}
}

func printTree(pager *Pager, pageNum int, indentationLevel int) {
	node := getPage(pager, pageNum)
	var numKeys, child int
	switch getNodeType(node) {
	case NODE_LEAF:
		numKeys, _ = leafNodeNumCells(node)
		indent(indentationLevel)
		fmt.Printf("- leaf (size %d)\n", numKeys)
		for i := 0; i < numKeys; i++ {
			indent(indentationLevel + 1)
			fmt.Printf("- %d\n", leafNodeKey(node, i))
		}
	case NODE_INTERNAL:
		numKeys, _ = internalNodeNumKeys(node)
		indent(indentationLevel)
		fmt.Printf("- internal (size %d)\n", numKeys)
		for i := 0; i < numKeys; i++ {
			child, _ = internalNodeChild(node, i)
			printTree(pager, child, indentationLevel+1)

			indent(indentationLevel + 1)
			t, _ := internalNodeKey(node, i)
			fmt.Printf("- key %d\n", t)
		}
		child, _ = internalNodeRightChild(node)
		printTree(pager, child, indentationLevel+1)
	}
}

// 将输入存储至页面
func serializeRow(source *Row, destination unsafe.Pointer) {
	id := Uint32ToBytes(source.ID)
	q := (*[ROW_SIZE]byte)(destination)
	copy(q[ID_OFFSET:ID_SIZE], id)
	copy(q[USERNAME_OFFSET:USERNAME_OFFSET+USERNAME_SIZE], source.UserName)
	copy(q[EMAIL_OFFSET:ROW_SIZE], source.Email)
}

// 将页面中的数据输出到对应结构体
func deserializeRow(source unsafe.Pointer, destination *Row) {
	id := make([]byte, ID_SIZE, ID_SIZE)
	sourceBuf := (*[ROW_SIZE]byte)(source)
	copy(id, sourceBuf[ID_OFFSET:ID_SIZE])

	destination.ID = BytesToInt32(id)
	username := make([]byte, USERNAME_SIZE, USERNAME_SIZE)
	copy(username, sourceBuf[USERNAME_OFFSET:USERNAME_OFFSET+USERNAME_SIZE])
	destination.UserName = string(username)

	email := make([]byte, EMAIL_SIZE, EMAIL_SIZE)
	copy(email, sourceBuf[EMAIL_OFFSET:ROW_SIZE])
	destination.Email = string(email)
}

// 从对应游标中查找对应的值
func cursorValue(cursor *Cursor) unsafe.Pointer {
	pageNum := cursor.pageNum
	page := getPage(cursor.table.Pager, pageNum)
	return leafNodeValue(page, cursor.cellNum)
}

// 将游标移向下一个节点
func cursorAdvance(cursor *Cursor) {
	pageNum := cursor.pageNum
	node := getPage(cursor.table.Pager, pageNum)
	// 将cellNum查找的值加1，移向下一个cell
	cursor.cellNum += 1
	// 如果超出当前的node的最大cells，则移向下一个叶子节点
	if cells, _ := leafNodeNumCells(node); cursor.cellNum >= cells {
		//cursor.tableEnd = true
		nextPageNum, _ := leafNodeNextLeaf(node)
		if nextPageNum == 0 {
			cursor.tableEnd = true
		} else {
			cursor.pageNum = nextPageNum
			cursor.cellNum = 0
		}
	}
}

// 打开页面
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
		numPages:   int(fileLength / PAGE_SIZE),
	}
	if fileLength%PAGE_SIZE != 0 {
		fmt.Println("Db file is not a whole number of pages. Corrupt file")
		os.Exit(0)
	}
	return pager, nil
}

// 打开数据库存储文件
func dbOpen(filename string) *Table {
	table := &Table{}
	table.Pager, _ = pageOpen(filename)
	table.rootPageNum = 0
	if table.Pager.numPages == 0 {
		rootNode := getPage(table.Pager, 0)
		initLeafNode(rootNode)
		setNodeRoot(rootNode, true)
	}
	return table
}

// pagerFlush 这一页写入文件系统
func pageFlush(pager *Pager, pageNum int) error {
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
	originBuf := make([]byte, PAGE_SIZE)
	buf := (*[PAGE_SIZE]byte)(pager.pages[pageNum])
	copy(originBuf[:PAGE_SIZE], buf[:PAGE_SIZE])
	bytesWritten, err := pager.fs.WriteAt(originBuf, offset)
	if err != nil {
		return fmt.Errorf("write %v", err)
	}
	// 捞取byte数组到这一页中
	fmt.Println("already wittern", bytesWritten)
	return nil
}

// 关闭数据库，将数据存储在数据库
func dbClose(table *Table) {
	for i := 0; i < table.Pager.numPages; i++ {
		if table.Pager.pages[i] == nil {
			continue
		}
		pageFlush(table.Pager, i)
	}
	defer table.Pager.fs.Close()
}

// 执行基础的命令
func doMetaCommand(input string, table *Table) metaCommandType {
	if input == ".exit" {
		dbClose(table)
		os.Exit(1)
		return metaCommandSuccess
	}
	if input == ".btree" {
		fmt.Printf("Tree:\n")
		printTree(table.Pager, 0, 0)
		return metaCommandSuccess
	}
	if input == ".constants" {
		fmt.Printf("Constants:\n")
		printConstants()
		return metaCommandSuccess
	}
	return metaCommandUnRecognizedCommand
}

// 对基础语句进行处理
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

func executeInsert(statement *Statement, table *Table) ExecuteResultType {
	row := statement.Row
	cursor := tableFind(table, int(row.ID))
	node := getPage(table.Pager, cursor.pageNum)
	numCells, _ := leafNodeNumCells(node)
	if cursor.cellNum < numCells {
		p := leafNodeKey(node, cursor.cellNum)
		indexKey := BytesToInt32((*[LEAF_NODE_KEY_SIZE]byte)(p)[:])
		if indexKey == row.ID {
			return EXECUTE_DUPLICATE_KEY
		}
	}
	leafNodeInsert(cursor, int(row.ID), row)
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

func executeStatement(statement *Statement, table *Table) ExecuteResultType {
	switch statement.statementType {
	case STATEMENT_SELECT:
		return executeSelect(statement, table)
	case STATEMENT_INSERT:
		return executeInsert(statement, table)
	}
	return EXECUTE_SUCCESS
}

func main() {
	run()
}
func run() {
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
		case EXECUTE_DUPLICATE_KEY:
			fmt.Println("Error: Duplicate key.")
		case EXECUTE_TABLE_FULL:
			fmt.Println("Err:Table full.")
		}
	}
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

var idx int = 0

func readInput_(reader *bufio.Reader) (string, error) {
	idx++
	s := fmt.Sprintf("insert %d user%d person1@example.com", idx, idx)
	if idx == 13 {
		s = fmt.Sprintf("select")
		s = fmt.Sprintf(".exit")
	}
	return s, nil
}

func readInput(reader *bufio.Reader) (string, error) {
	buf := bufio.NewReader(os.Stdin)
	data, err := buf.ReadBytes('\n')
	data = data[:len(data)-2]
	for len(data) > 0 && data[0] == ' ' {
		data = data[1:]
	}
	for len(data) > 0 && data[len(data)-1] == ' ' {
		data = data[:len(data)-1]
	}
	res := []byte{}
	for i := 0; i < len(data); i++ {
		if i > 0 && data[i] == ' ' && data[i-1] == ' ' {
			continue
		}
		res = append(res, data[i])
	}
	return string(res), err
}
func printPrompt() {
	fmt.Print("db>")
}
func printRow(row *Row) {
	fmt.Printf("(%d,%s,%s)\n", row.ID, row.UserName, row.Email)
}
