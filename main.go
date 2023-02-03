package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Person struct {
	Id         int
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/gotest")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/persons", func(c *gin.Context) {
		var (
			person  Person
			persons []Person
		)
		rows, err := db.Query("select id, first_name, last_name from person;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&person.Id, &person.First_Name, &person.Last_Name)
			persons = append(persons, person)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"result": persons,
			"count":  len(persons),
		})
	})

	r.GET("/person/:id", func(c *gin.Context) {
		var (
			person Person
			result gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("select id, first_name, last_name from person where id = ?;", id)
		err = row.Scan(&person.Id, &person.First_Name, &person.Last_Name)
		if err != nil {
			// If no results send null
			result = gin.H{
				"result": nil,
				"count":  0,
			}
		} else {
			result = gin.H{
				"result": person,
				"count":  1,
			}
			fmt.Println("get person by ID")
		}

		c.JSON(http.StatusOK, result)
	})

	r.POST("/person", func(c *gin.Context) {
		var p Person
		c.BindJSON(&p)

		stmt, err := db.Prepare("insert into person (first_name, last_name) values(?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(p.First_Name, p.Last_Name)

		if err != nil {
			fmt.Print(err.Error())
		}

		output := p.First_Name + " " + p.Last_Name + " created"
		fmt.Println(output)

		c.JSON(http.StatusOK, gin.H{
			"message": output,
		})
	})

	r.PUT("/person", func(c *gin.Context) {
		var p Person
		c.BindJSON(&p)

		stmt, err := db.Prepare("update person set first_name= ?, last_name= ? where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(p.First_Name, p.Last_Name, p.Id)

		if err != nil {
			fmt.Print(err.Error())
		}

		output := "Successful update id  " + strconv.Itoa(p.Id)
		fmt.Println(output)

		c.JSON(http.StatusOK, gin.H{
			"message": output,
		})
	})

	r.DELETE("/person", func(c *gin.Context) {
		var p Person
		c.BindJSON(&p)

		stmt, err := db.Prepare("delete from person where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(p.Id)
		if err != nil {
			fmt.Print(err.Error())
		}

		output := "Delete person id  " + strconv.Itoa(p.Id)
		fmt.Println(output)

		c.JSON(http.StatusOK, gin.H{
			"message": output,
		})
	})

	r.Run(":8081")
}
