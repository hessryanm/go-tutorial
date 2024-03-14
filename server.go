package main

import(
	"database/sql"
	"time"

	"github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type todo struct {
	Item string
	Done bool
	ID int
}

func saveLog(message string, myLog chan string) {
	myLog <- message
}

func handleError(c *fiber.Ctx, errStr string) error {
	errStr = "An error occurred: " + errStr
	
	log.Fatal(errStr)
	return c.SendString(errStr)
}

func indexHandler(c *fiber.Ctx, db *sql.DB, myLog chan string) error {
	var res todo
	var todos []todo

	rows, err := db.Query("SELECT item, done, id FROM todos")
	defer rows.Close()

	if err != nil {
		return handleError(c, err.Error())
	}

	for rows.Next() {
		rows.Scan(&res.Item, &res.Done, &res.ID)
		todos = append(todos, res)
	}

	go saveLog("Got Index", myLog)
	return c.JSON(todos)
}

func postHandler(c *fiber.Ctx, db *sql.DB, myLog chan string) error {
	newTodo := todo{}
	if err := c.BodyParser(&newTodo); err != nil {
		return handleError(c, err.Error())
	}

	if newTodo.Item != "" {
		_, err := db.Exec("INSERT INTO todos VALUES ($1)", newTodo.Item)
		if err != nil {
			return handleError(c, err.Error())
		}
	}

	go saveLog("Created new item " + newTodo.Item, myLog)
	return c.JSON(newTodo)
}

func putHandler(c *fiber.Ctx, db *sql.DB, myLog chan string) error {
	m := c.Queries()
	newTodo := todo{m["updateItem"], false, 0}

	newItem, hasNewItem := m["newItem"]
	if hasNewItem {
		newTodo.Item = newItem
	}

	done, hasDone := m["done"]
	if hasDone && done == "true" {
		newTodo.Done = true
	}
	parameters := []any{newTodo.Item}

	queryStr := "UPDATE todos SET item = $1 "
	if hasDone {
		queryStr += ", done = $2 "
		parameters = append(parameters, newTodo.Done)
	}
	queryStr += "WHERE item = $3"
	parameters = append(parameters, m["updateItem"])

	_, err := db.Exec(queryStr, parameters...)
	if err != nil {
		return handleError(c, err.Error())
	}

	err = db.QueryRow("SELECT * FROM todos WHERE item = $1", newTodo.Item).Scan(&newTodo.Item, &newTodo.Done)

	go saveLog("Updated item " + newTodo.Item, myLog)
	return c.JSON(newTodo)
}

func deleteHandler(c *fiber.Ctx, db *sql.DB, myLog chan string) error {
	_, err := db.Exec("DELETE FROM todos WHERE item = $1", c.Query("item"))
	if err != nil {
		return handleError(c, err.Error())
	}

	go saveLog("Deleted item " + c.Query("item"), myLog)
	return c.SendString("Deleted item " + c.Query("item"))
}

func startLog(myLog chan string) {
	for message := range myLog {
		log.Debug(message)
	}
}

func startDbLog(l *pq.Listener) {
	select {
        case notification := <-l.Notify:
            log.Debug("DB log: " + notification.Extra)
        case <-time.After(90 * time.Second):
            go l.Ping()
            // Check if there's more work available, just in case it takes
            // a while for the Listener to notice connection loss and
            // reconnect.
            log.Debug("received no work for 90 seconds, checking for new work")
    }
}

func main() {
	dbConnStr := "postgresql://hess@localhost/postgres?sslmode=disable"
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal(err)
	}

	myLog := make(chan string)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
        return indexHandler(c, db, myLog)
    })
	app.Post("/", func(c *fiber.Ctx) error {
        return postHandler(c, db, myLog)
    })
	app.Put("/update", func(c *fiber.Ctx) error {
        return putHandler(c, db, myLog)
    })
	app.Delete("/delete", func(c *fiber.Ctx) error {
        return deleteHandler(c, db, myLog)
    })

    go startLog(myLog)

    minReconn := 10 * time.Second
    maxReconn := time.Minute
    reportProblem := func(ev pq.ListenerEventType, err error) {
        if err != nil {
            log.Panic(err.Error())
        }
    }

    listener := pq.NewListener(dbConnStr, minReconn, maxReconn, reportProblem)
    listenErr := listener.Listen("ithappened")
    if listenErr != nil {
        panic(listenErr)
    }

    go startDbLog(listener)

	log.Fatal(app.Listen(":3000"))
}