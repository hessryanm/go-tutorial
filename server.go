package main

import(
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
)

type todo struct {
	Item string
}

func handleError(c *fiber.Ctx, errStr string) error {
	errStr = "An error occurred: " + errStr
	
	log.Fatalln(errStr)
	return c.SendString(errStr)
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	var res todo
	var todos []todo

	rows, err := db.Query("SELECT * FROM todos")
	defer rows.Close()

	if err != nil {
		return handleError(c, err.Error())
	}

	for rows.Next() {
		rows.Scan(&res.Item)
		todos = append(todos, res)
	}

	return c.JSON(todos)
}

func postHandler(c *fiber.Ctx, db *sql.DB) error {
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

	return c.Redirect("/")
}

func putHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("Not Implemented")
}

func deleteHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("Not Implemented")
}

func main() {
	dbConnStr := "postgresql://hess@localhost/postgres?sslmode=disable"
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalln(err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
        return indexHandler(c, db)
    })
	app.Post("/", func(c *fiber.Ctx) error {
        return postHandler(c, db)
    })
	app.Put("/update", func(c *fiber.Ctx) error {
        return putHandler(c, db)
    })
	app.Delete("/delete", func(c *fiber.Ctx) error {
        return deleteHandler(c, db)
    })

	log.Fatalln(app.Listen(":3000"))
}