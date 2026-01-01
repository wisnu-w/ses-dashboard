package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Database connection
	dsn := "host=localhost port=5432 user=ses_user password=password123! dbname=ses_dashboard sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	fmt.Println("✅ Database connection successful")

	// Check if users table exists
	var tableExists bool
	err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&tableExists)
	if err != nil {
		log.Fatal("Failed to check table:", err)
	}
	fmt.Printf("✅ Users table exists: %v\n", tableExists)

	// Check admin user
	var id int
	var username, password, email, role string
	var active bool
	err = db.QueryRow("SELECT id, username, password, email, role, active FROM users WHERE username = 'admin'").Scan(&id, &username, &password, &email, &role, &active)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("❌ Admin user not found")
		} else {
			log.Fatal("Failed to query admin user:", err)
		}
		return
	}

	fmt.Printf("✅ Admin user found:\n")
	fmt.Printf("   ID: %d\n", id)
	fmt.Printf("   Username: %s\n", username)
	fmt.Printf("   Email: %s\n", email)
	fmt.Printf("   Role: %s\n", role)
	fmt.Printf("   Active: %v\n", active)
	fmt.Printf("   Password Hash: %s\n", password)

	// Test password verification
	testPasswords := []string{"secret", "admin123", "admin", "password"}
	for _, testPwd := range testPasswords {
		err = bcrypt.CompareHashAndPassword([]byte(password), []byte(testPwd))
		if err == nil {
			fmt.Printf("✅ Password '%s' matches!\n", testPwd)
		} else {
			fmt.Printf("❌ Password '%s' does not match\n", testPwd)
		}
	}
}
