// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	_ "github.com/lib/pq"
// )

// var DB *sql.DB

// // User struct maps to the users table
// type User struct {
// 	UserID   int    `json:"userId"`
// 	Name     string `json:"name"`
// 	Email    string `json:"email"`
// 	IsActive bool   `json:"isActive"`
// }

// func main() {
// 	fmt.Println("Hello, World!")

// 	// Initialize database connection
// 	if err := CreateDbObject(); err != nil {
// 		log.Fatalf("Database connection failed: %v", err)
// 	}

// 	// Set up Fiber app
// 	app := fiber.New()

// 	// Routes
// 	app.Get("/abc", func(c *fiber.Ctx) error {
// 		return c.SendString("Syntax Squad is awesome.")
// 	})

// 	app.Get("/users", GetAllUsers)

// 	app.Post("/signin", SignInWeb)

// 	// Start server
// 	log.Fatal(app.Listen(":3000"))
// }

// // PostgreSQL connection string
// func GetPsqlInfo() string {
// 	return fmt.Sprintf(
// 		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
// 		"dpg-d1ibun6r433s73abdee0-a.oregon-postgres.render.com",
// 		"5432",
// 		"database_nejo_user",
// 		"kqXECzS8xa574BPtVeOBwkuGavTYVU3O",
// 		"database_nejo",
// 	)
// }

// // CreateDbObject sets up the database connection
// func CreateDbObject() error {
// 	var err error
// 	DB, err = sql.Open("postgres", GetPsqlInfo())
// 	if err != nil {
// 		return fmt.Errorf("error opening DB: %w", err)
// 	}

// 	if err = DB.Ping(); err != nil {
// 		return fmt.Errorf("error connecting to DB: %w", err)
// 	}

// 	fmt.Println("Connected to database successfully")

// 	// Set connection pooling options
// 	DB.SetMaxOpenConns(25)
// 	DB.SetMaxIdleConns(25)
// 	DB.SetConnMaxIdleTime(10 * time.Minute)
// 	DB.SetConnMaxLifetime(1 * time.Hour)

// 	return nil
// }

// // GetAllUsers handles GET /users
// func GetAllUsers(c *fiber.Ctx) error {
// 	rows, err := DB.Query("SELECT userId, name, email, isActive FROM users")
// 	if err != nil {
// 		log.Printf("Query failed: %v", err)
// 		return c.Status(500).SendString("Internal Server Error")
// 	}
// 	defer rows.Close()

// 	var users []User

// 	for rows.Next() {
// 		var u User
// 		if err := rows.Scan(&u.UserID, &u.Name, &u.Email, &u.IsActive); err != nil {
// 			log.Printf("Scan failed: %v", err)
// 			return c.Status(500).SendString("Error reading data")
// 		}
// 		users = append(users, u)
// 	}

// 	if err = rows.Err(); err != nil {
// 		log.Printf("Row iteration error: %v", err)
// 		return c.Status(500).SendString("Error retrieving data")
// 	}

// 	return c.JSON(users)
// }

// // signin
// type SignInRequest struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// func SignInWeb(c *fiber.Ctx) error {
// 	signInRequestObject := SignInRequest{}
// 	if err := c.BodyParser(&signInRequestObject); err != nil {
// 		c.Status(500)
// 		return nil
// 	}
// 	fmt.Println(signInRequestObject)
// 	c.JSON("Its WORKING")
// 	return nil
// }

// func FetchUserIDFromEmailID(email string) (int, error) {
// 	query := fmt.Sprintf(`SELECT userid FROM "user" WHERE email ="%s" `,email)

// 	rows, err := DB.Query(query)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer rows.Close()

// 		var u int
// 		err = rows.Scan(&u)
// 		if err != nil {
// 			return 0, err
// 		}

// 	return u, nil
// }

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

var DB *sql.DB

// User maps to the `users` table
type User struct {
	UserID   int    `json:"userId"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	IsActive bool   `json:"isActive"`
}

func main() {
	fmt.Println("Hello, Fiber + Postgres!")

	// 1️⃣ connect to Postgres
	if err := createDB(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// 2️⃣ start Fiber
	app := fiber.New()

	app.Get("/abc", func(c *fiber.Ctx) error {
		return c.SendString("Syntax Squad is awesome.")
	})
	app.Get("/users", getAllUsers)
	app.Post("/signin", signInWeb)

	log.Fatal(app.Listen(":3000"))
}

// assemble connection string
func pgConnString() string {
	host := getEnv("PGHOST", "dpg-d1ibun6r433s73abdee0-a.oregon-postgres.render.com")
	port := getEnv("PGPORT", "5432")
	user := getEnv("PGUSER", "database_nejo_user")
	pass := getEnv("PGPASSWORD", "kqXECzS8xa574BPtVeOBwkuGavTYVU3O")
	db := getEnv("PGDATABASE", "database_nejo")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, pass, db,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func createDB() error {
	var err error
	DB, err = sql.Open("postgres", pgConnString())
	if err != nil {
		return fmt.Errorf("open DB: %w", err)
	}
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("ping DB: %w", err)
	}

	log.Println("✅ Connected to database")

	// pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxIdleTime(10 * time.Minute)
	DB.SetConnMaxLifetime(1 * time.Hour)
	return nil
}

// GET /users
func getAllUsers(c *fiber.Ctx) error {
	rows, err := DB.Query(`SELECT userid, name, email, isactive FROM users`)
	if err != nil {
		log.Printf("query: %v", err)
		return c.Status(500).SendString("Internal Server Error")
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.UserID, &u.Name, &u.Email, &u.IsActive); err != nil {
			log.Printf("scan: %v", err)
			return c.Status(500).SendString("Error reading data")
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		log.Printf("rows err: %v", err)
		return c.Status(500).SendString("Error retrieving data")
	}
	return c.JSON(users)
}

// POST /signin

func signInWeb(c *fiber.Ctx) error {
	var req SignInRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).SendString("Bad request")
	}
	fmt.Printf("sign‑in attempt: %+v\n", req)
	return c.JSON(fiber.Map{"message": "It works!"})
}

//func signInWeb(c*fiber.txt)
