package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID          int64    `form:"id"`
	CreatedAt   string
	UpdatedAt   string
	Title       string   `form:"title"`
	Complete    bool     `form:"complete"`
}

var db *sql.DB

func main() {
	dbSetup()
	defer db.Close()
	dbTablesSetup()

	r := gin.Default()
	r.LoadHTMLGlob("views/*")
	r.Static("/public", "./public")

	r.GET("/", getIndex)
	r.POST("/todos", postAddTodo)
	r.PUT("/todos/:id", putTodoById)
	r.DELETE("/todos/:id", deleteTodoById)

	r.Run("127.0.0.1:3000")
}

func getIndex(c *gin.Context) {
	todos, err := getAllTodos()

	if err != nil {
		c.HTML(http.StatusOK, "error.html", gin.H{
			"message": err.Error(),
		})
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"todos": todos,
		"count": len(todos),
	})
}

func postAddTodo(c *gin.Context) {
	createNewTodoQuery, err := loadQuery("createNewTodo.sql")
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load query"})
		return
	}

	var todo Todo;

	err = c.Bind(&todo);
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed bind todo"})
		return
	}

	result, err := db.Exec(createNewTodoQuery, false, todo.Title)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	todo.Complete = false
	todoId, _ := result.LastInsertId()
	todo.ID = todoId

	c.HTML(http.StatusCreated, "todo.html", todo)
}

func putTodoById(c *gin.Context) {
	todoId := c.Param("id")

	var todo Todo;

	err := c.Bind(&todo);
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed parse todo data"})
		return
	}

	updateTodoByIdQuery, err := loadQuery("updateTodoById.sql")
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load query"})
		return
	}

	_, err = db.Exec(updateTodoByIdQuery, !todo.Complete, todo.Title, todoId)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}

	todo.Complete = !todo.Complete

	c.HTML(http.StatusCreated, "todo.html", todo)
}

func deleteTodoById(c *gin.Context) {
	todoId := c.Param("id")

	deleteTodoByIdQuery, err := loadQuery("deleteTodoById.sql")
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load query"})
		return
	}

	_, err = db.Exec(deleteTodoByIdQuery, todoId)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}
}

func loadQuery(filename string) (string, error) {
	query, err := os.ReadFile("./sql/" + filename)
	if err != nil {
		return "", err
	}
	return string(query), nil
}

func getAllTodos() ([]Todo, error) {
	getAllTodosQuery, err := loadQuery("getAllTodos.sql")
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(getAllTodosQuery)
	if err != nil {
		return nil, err
	}

	var todos []Todo

	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.CreatedAt, &todo.UpdatedAt, &todo.Complete, &todo.Title)
		todos = append(todos, todo)
	}

	return todos, nil
}

func dbSetup() {
	var err error	
	db, err = sql.Open("sqlite3", "todos.db")
	if err != nil {
		panic(err)
	}
}

func dbTablesSetup() {
	createTodosTableQuery, err := loadQuery("createTable.sql")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createTodosTableQuery)
	if err != nil {
		panic(err)
	}
}
