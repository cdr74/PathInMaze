package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
)

// ---------------------------------------------------------------------------

const runTest bool = false
const TEST_FILE string = "test.data"
const DATA_FILE string = "actual.data"

var X_SIZE int
var Y_SIZE int

var START Position
var END Position

var maze [][]Node
var queue *Queue = NewQueue()

// ---------------------------------------------------------------------------

type Position struct {
	x int
	y int
}

type Node struct {
	Height      int
	Distance    int // distance from start
	Predecessor Position
}

// ---------------------------------------------------------------------------

type Queue struct {
	list *list.List
}

func NewQueue() *Queue {
	return &Queue{list: list.New()}
}

func (q *Queue) Enqueue(position Position) {
	q.list.PushBack(position)
}

func (q *Queue) Dequeue() Position {
	// fails if empty
	element := q.list.Front()
	q.list.Remove(element)
	return element.Value.(Position) // Adjust the type assertion here
}

func (q *Queue) IsEmpty() bool {
	return q.list.Len() == 0
}

func (q *Queue) Size() int {
	return q.list.Len()
}

// ---------------------------------------------------------------------------

func ReadDataFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		// log.Fates does exit the prog
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	return lines
}

func StringToIntArray(input string) []int {
	intArray := make([]int, 0, len(input))

	for _, char := range input {
		var num int
		switch char {
		case 'S':
			num = 96 // 'a' - 1; this is an ugly hack to pass info to InitializeMaze
		case 'E':
			num = 123 // 'z' + 1; this is an ugly hack to pass info to InitializeMaze
		default:
			// Get the ASCII value of the character
			num = int(char)
		}

		intArray = append(intArray, num)
	}

	return intArray
}

func InitializeMaze(input []string) {
	X_SIZE = len(input[0]) + 2
	Y_SIZE = len(input) + 2

	maze = make([][]Node, Y_SIZE)
	for i := range maze {
		maze[i] = make([]Node, X_SIZE)
	}

	// initialize with a value too hight to climb as a border
	for i := 0; i < Y_SIZE; i++ {
		for j := 0; j < X_SIZE; j++ {
			maze[i][j] = Node{
				Height:      9999,
				Distance:    9999,
				Predecessor: Position{x: -1, y: -1},
			}
		}
	}

	// now add input data leaving the border too high
	for i := 0; i < len(input); i++ {
		line := input[i]
		values := StringToIntArray(line)
		for j := 0; j < len(values); j++ {
			value := values[j]
			switch value {
			case 96:
				START = Position{y: i + 1, x: j + 1}
				maze[i+1][j+1].Height = 97
				maze[i+1][j+1].Distance = 0
				maze[i+1][j+1].Predecessor.x = i + 1 // self
				maze[i+1][j+1].Predecessor.y = j + 1 // self
			case 123:
				END = Position{y: i + 1, x: j + 1}
				maze[i+1][j+1].Height = 122
			default:
				maze[i+1][j+1].Height = value
			}
		}
	}
}

func PrintMaze() {
	fmt.Printf("Start(%d, %d)\n", START.x, START.y)
	fmt.Printf("End(%d, %d)\n", END.x, END.y)

	for i := 0; i < Y_SIZE; i++ {
		for j := 0; j < X_SIZE; j++ {
			fmt.Printf("(%d, %d): Height=%d, Distance=%d, Predecessor=%d\n", i, j, maze[i][j].Height, maze[i][j].Distance, maze[i][j].Predecessor)
		}
	}
}

// ---------------------------------------------------------------------------

// Rules on what nodes are worth visiting:
// - target has a height <= current + 1
// - target has a distance > current + 1 (already visited, onyl update if we found shorter path)
func AddToQueue(pos Position) {
	newDistance := maze[pos.y][pos.x].Distance + 1

	// check left
	new_x := pos.x - 1
	new_y := pos.y
	if maze[new_y][new_x].Height <= maze[pos.y][pos.x].Height+1 && maze[new_y][new_x].Distance > maze[pos.y][pos.x].Distance+1 {
		next := Position{x: new_x, y: new_y}
		queue.Enqueue(next)
		maze[new_y][new_x].Predecessor = Position{x: pos.x, y: pos.y}
		maze[new_y][new_x].Distance = newDistance
	}
	// check right
	new_x = pos.x + 1
	new_y = pos.y
	if maze[new_y][new_x].Height <= maze[pos.y][pos.x].Height+1 && maze[new_y][new_x].Distance > maze[pos.y][pos.x].Distance+1 {
		next := Position{x: new_x, y: new_y}
		queue.Enqueue(next)
		maze[new_y][new_x].Predecessor = Position{x: pos.x, y: pos.y}
		maze[new_y][new_x].Distance = newDistance
	}
	// check up
	new_x = pos.x
	new_y = pos.y - 1
	if maze[new_y][new_x].Height <= maze[pos.y][pos.x].Height+1 && maze[new_y][new_x].Distance > maze[pos.y][pos.x].Distance+1 {
		next := Position{x: new_x, y: new_y}
		queue.Enqueue(next)
		maze[new_y][new_x].Predecessor = Position{x: pos.x, y: pos.y}
		maze[new_y][new_x].Distance = newDistance
	}
	// check down
	new_x = pos.x
	new_y = pos.y + 1
	if maze[new_y][new_x].Height <= maze[pos.y][pos.x].Height+1 && maze[new_y][new_x].Distance > maze[pos.y][pos.x].Distance+1 {
		next := Position{x: new_x, y: new_y}
		queue.Enqueue(next)
		maze[new_y][new_x].Predecessor = Position{x: pos.x, y: pos.y}
		maze[new_y][new_x].Distance = newDistance
	}
}

func FindPath() int {
	result := -1
	AddToQueue(START)
	for !queue.IsEmpty() {
		var next Position = queue.Dequeue()
		if next.x == END.x && next.y == END.y {
			result = maze[next.y][next.x].Distance
		} else {
			AddToQueue(next)
		}
	}
	return result
}

// ---------------------------------------------------------------------------

func main() {
	var result int
	var input []string

	if runTest {
		input = ReadDataFile(TEST_FILE)
	} else {
		input = ReadDataFile(DATA_FILE)
	}

	InitializeMaze(input)
	PrintMaze()
	result = FindPath()
	PrintMaze()

	// -------------------------------------
	fmt.Println("Running as test:\t", runTest)
	fmt.Println("Result:\t\t\t", result)
}
