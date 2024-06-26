package database 

import (
	"log"
	"database/sql"

	_"github.com/lib/pq"
)


var DB *sql.DB  // global database connection variable 

// struct that represents the user registration model 
type UserRegistrationRequest struct {
	ID			int		`json:"id"` 
	Username	string	`json:"username"`
	Email		string	`json:"email"`
	Password	string	`json:"password"`
	Password2	string 	`json:"password2"`
	// IsAdmin		bool	`json:"isAdmin"`
	// IsActive	bool	`json:"isActive"` 
}

// struct that represents user login model 
type UserLoginRequest struct{
	ID		  int	  `json:"id"`
	LoginId   string  `json:"login_id"` // this allows user login with either username or email
	Password  string  `json:"password"`    
}

// struct that represents the user model 
type User struct {
	ID			int		`json:"id"`
	Username	string	`json:"username"`
	Email		string	`json:"email"`
	Password	string	`json:"password"`
	IsAdmin		bool	`json:"isAdmin"`
	IsActive	bool	`json:"isActive"` 
}

// struct that represents the profile model 
type Profile struct {
	ID			    int		`json:"id"`
	FullName	    string	`json:"fullname"`
	Username	    string	`json:"username"`
	Email		    string 	`json:"email"`
	Bio			    string 	`json:"bio"`
	ProfilePicture  string  `json:"profile_pic"`	
	// embedding User struct into profile struct
	UserObj		    User  // user struct instance on profile 
}

// initialize database connection 
func InitDB(connectionString string) error {
	var err error

	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping() // test database connection 
	if err != nil {
		log.Println("Unable to connect to database")
		log.Fatal(err)
	}

	log.Println("Connected to database") // print message if connection is successful

	// defer DB.Close() // defer database connection 
	return nil 

}


// migration function
func CreateMigrations() error {
	// Create user table
	createUserTable := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(150) UNIQUE NOT NULL,
			password TEXT NOT NULL,
			is_admin BOOLEAN DEFAULT false,
			is_active BOOLEAN DEFAULT true
		);
	`
	_, err := DB.Exec(createUserTable)
	if err != nil {
		log.Fatal("Unable to create database table 'user':", err.Error())
		return err
	}
	log.Println("Migrations successfully created for users")

	// Create profile table
	createProfileTable := `
		CREATE TABLE IF NOT EXISTS profile (
			id SERIAL PRIMARY KEY,
			fullname VARCHAR(60),
			username VARCHAR(50) NOT NULL,
			email VARCHAR(150) NOT NULL,
			bio TEXT,
			profile_picture BYTEA,
			user_id INT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`
	_, err = DB.Exec(createProfileTable)
	if err != nil {
		log.Fatal("Unable to create database table 'profile':", err.Error())
		return err
	}
	log.Println("Migrations successfully created for profile")

	// Create blacklisted_token table
	createBlacklistedTokenTable := `
		CREATE TABLE IF NOT EXISTS blacklisted_token (
			token  VARCHAR(300)
		);
	`
	_, err = DB.Exec(createBlacklistedTokenTable)
	if err != nil {
		log.Fatal("Unable to create database table for 'blacklisted_token':", err.Error())
		return err
	}
	log.Println("Migrations successfully created for blacklisted_token")
	
	return nil
}