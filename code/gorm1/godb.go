package main

import (
    "fmt"
    "time"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

type Pessoa struct {
    Cpf string `gorm:"primary_key;type:varchar(11);column:cpf"`
    Name string `gorm:"type:varchar(50);column:nome"`
    BirthDate *time.Time `gorm:"column:data_nascimento"`
    Employee bool `gorm:"column:funcionario"`
}

func main() {

    // Open connection to a postgresql database running on ElephantQL:
    db, err := gorm.Open("postgres", "host=elmer.db.elephantsql.com port=5432 user=userdbname dbname=userdbname password=password")
    if err != nil {
        panic("failed to connect database")
    }
    defer db.Close()
    fmt.Println("OK")

    // Let's start a transaction:
    tx := db.Begin()

    // Let's create some persons:
    dt,_ := time.Parse("2006-01-02", "1979-08-18")
    person1 := &Pessoa{"111","Person#1",&dt,false}
    tx.Table("public.pessoa").Create(&person1)
    person2 := &Pessoa{"222","Person#2",&dt,false}
    tx.Table("public.pessoa").Create(&person2)
    person3 := &Pessoa{"333","Person#3",&dt,false}
    tx.Table("public.pessoa").Create(&person3)

    // Commiting the transaction:
    tx.Commit()

    // Select data from first record:
    var person Pessoa
    if result := db.Table("public.pessoa").First(&person); result.Error != nil {
      panic(result.Error)
    }
    fmt.Println("First person: ",person)

    // Select all: 
    var persons []Pessoa
    if result := db.Table("public.pessoa").Find(&persons); result.Error != nil {
      panic(result.Error)
    }
    fmt.Printf("Persons: %+v\n",persons)

    // Select / where:
    var person222 Pessoa
    if result := db.Table("public.pessoa").Where("cpf like ?", "222").First(&person222); result.Error != nil {
      panic(result.Error)
    }
    fmt.Println("Person identified by cpf: ",person222)



    // Update: 
    fmt.Println("Before updating: ",person3.BirthDate)
    if result := db.Table("public.pessoa").Where("cpf like ?", "333").First(&person3); result.Error != nil {
      panic(result.Error)
    }
    newdate, _ := time.Parse("2006-01-02", "1990-11-01")
    if result := db.Table("public.pessoa").Model(&person3).Update("data_nascimento", &newdate); result.Error != nil {
      panic(result.Error)
    }
    fmt.Println("Person data updated",person3.BirthDate)

    // Delete: 

    if result := db.Table("public.pessoa").Where("cpf like ?", "111").First(&person1); result.Error != nil {
      panic(result.Error)
    }
    fmt.Println(person1)
    if result := db.Table("public.pessoa").Delete(&person1); result.Error != nil {
      panic(result.Error)
    }
    fmt.Println("Person deleted")
 
    if result := db.Table("public.pessoa").Where("cpf like ?", "222").First(&person2); result.Error != nil {
      panic(result.Error)
    }
    if result := db.Table("public.pessoa").Delete(&person2); result.Error != nil {
      panic(result.Error)
    }
    fmt.Println("Person deleted")

    if result := db.Table("public.pessoa").Where("cpf like ?", "333").First(&person3); result.Error != nil {
      panic(result.Error)
    }
    if result := db.Table("public.pessoa").Delete(&person3); result.Error != nil {
      panic(result.Error)
    }
    fmt.Println("Person deleted")

}
