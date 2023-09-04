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
	ID int `json:"id"`
	CreatedAt string
	UpdatedAt string
	Title string `json:"title"`
	Complete bool `json:"complete"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "todos.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	createTodosTableQuery, err := loadQuery("createTable.sql");
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createTodosTableQuery)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.LoadHTMLGlob("views/*")
	r.Static("/public", "./public");

	r.GET("/", getIndex)
	r.POST("/todos", postAddTodo)
	r.PUT("/todos/:id", putTodoById)
	r.DELETE("/todos/:id", deleteTodoById)

	r.Run("127.0.0.1:3000")
}

func getIndex(c *gin.Context) {
	todos, err := getAllTodos()

	if err != nil {
		c.HTML(http.StatusOK, "index.html", nil)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"todos": todos,
	})
}

func postAddTodo(c *gin.Context) {
	createNewTodoQuery, err := loadQuery("createNewTodo.sql")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load query"})
		return
	}

	todo := c.PostForm("title");

	_, err = db.Exec(createNewTodoQuery, false, todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	todoHTML := "<p class='todo'>" + todo + "</p>"

	c.String(http.StatusOK, todoHTML);
}

func putTodoById(c *gin.Context) {
	todoId := c.Param("id")

	updateTodoByIdQuery, err := loadQuery("updateTodoById.sql")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load query"})
		return
	}

	result, err := db.Exec(updateTodoByIdQuery, true, "Update", todoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}
	
	fmt.Println(result)
}

func deleteTodoById(c *gin.Context) {
	todoId := c.Param("id")

	deleteTodoByIdQuery, err := loadQuery("deleteTodoById.sql")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load query"})
		return
	}

	result, err := db.Exec(deleteTodoByIdQuery, todoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}

	fmt.Println(result)
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
