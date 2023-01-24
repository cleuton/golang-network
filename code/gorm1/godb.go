package main

import (
    "fmt"
    "os"
    "network.golang/curso-gorm/db"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {

    // Run postgres: docker run -p 5432:5432 --name <name> -e POSTGRES_PASSWORD=<password> -d postgres

    // Get credentials from Environment Variables:
    databaseName := os.Getenv("DATABASENAME") 
    userName := os.Getenv("USERNAME")
    password := os.Getenv("PASSWORD")
    host := os.Getenv("DATABASEHOST")

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=America/Sao_Paulo",
                      host, userName, password, databaseName)
    fmt.Println(dsn)

    // Open connection to a postgresql database:
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
      DisableForeignKeyConstraintWhenMigrating: true,
    })
    if err != nil {
        panic("failed to connect database")
    }

    fmt.Println("OK")


    // Let's run the migrations to create tables:

    db.Migrator().AutoMigrate(&model.Member{}, &model.Note{})

    // Let's create the constrants:
    db.Migrator().CreateConstraint(&model.Member{}, "Notes")
    db.Migrator().CreateConstraint(&model.Member{}, "Connections")
   
    // Let's add some data:
    user1 := &model.Member{Name: "Paul", Email: "paul@testmail"}
    db.Create(&user1)
    user2 := &model.Member{Name: "John", Email: "john@test"}
    db.Create(&user2)

    // Connect John to Paul: 
    db.Model(&user2).Association("Connections").Append(user1)
    db.Save(user2)

    // Paul wrote a note:
    db.Model(&user1).Association("Notes").Append(&model.Note{Author: *user1, Text: "This is a note"})
    db.Save(user1)

    // Now get user Paul from database:
    paul := &model.Member{}
    db.Preload("Notes").Preload("Connections").First(&paul,1)
    fmt.Printf("\nPaul has a note with text: %s\n",paul.Notes[0].Text)
    fmt.Printf("\nPaul has %d connections\n", len(paul.Connections))

    // Now get user John, who is connected to Paul:
    john := &model.Member{}
    db.Preload("Notes").Preload("Connections").First(&john,2)
    fmt.Printf("\nJohn has %d notes\n", len(john.Notes))
    fmt.Printf("\nJohn has a connection with %s\n",john.Connections[0].Name)    


}
