package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

/**
* Create database test
* CREATE TABLE `pages` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `guid` varchar(256) NOT NULL,
  `title` varchar(256) DEFAULT NULL,
  `content` mediumtext,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `pages_guid_IDX` (`guid`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=latin1
*/

//DBHost,DBPort,DBUser,DBPassword,DBBase, PORT const
const (
	DBHost     = "127.0.0.1"
	DBPort     = "3306"
	DBUser     = "root"
	DBPassword = ""
	DBBase     = "test"
	PORT       = ":9000"
)

//Page struct
type Page struct {
	Title      string
	RawContent string
	Content    template.HTML
	Date       string
	GUID       string
}

var database *sql.DB

func main() {
	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUser, DBPassword, DBHost, DBPort, DBBase)
	fmt.Println(dbConn)
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Println("Could not connect!")
		log.Println(err)
	}

	database = db

	router := gin.Default()

	//GET Page with GUID
	router.GET("/api/page/:guid", func(c *gin.Context) {
		var (
			page   Page
			result gin.H
		)

		guid := c.Param("guid")
		row := database.QueryRow("SELECT title,content,date FROM pages WHERE guid=?", guid)
		err := row.Scan(&page.Title, &page.RawContent, &page.Date)
		page.Content = template.HTML(page.RawContent)
		if err != nil {
			result = gin.H{
				"result": nil,
				"count":  1,
			}
		} else {
			result = gin.H{
				"result": page,
				"count":  1,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	//GET All Page
	router.GET("/api/pages", func(c *gin.Context) {
		var (
			page  Page
			pages []Page
		)
		rows, err := db.Query("SELECT guid,title,content,date FROM pages;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&page.GUID, &page.Title, &page.RawContent, &page.Date)
			page.Content = template.HTML(page.RawContent)
			pages = append(pages, page)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"result": pages,
			"count":  len(pages),
		})
	})

	//Insert Page
	router.POST("/api/page", func(c *gin.Context) {
		var buffer bytes.Buffer
		guid := c.PostForm("guid")
		title := c.PostForm("title")
		content := c.PostForm("content")
		stmt, err := database.Prepare("INSERT INTO pages (guid,title,content) values(?,?,?);")
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = stmt.Exec(guid, title, content)
		if err != nil {
			fmt.Println(err.Error())
		}

		buffer.WriteString(title)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("%s Successfully created", name),
		})
	})

	//Update Page
	router.PUT("/api/page", func(c *gin.Context) {
		var buffer bytes.Buffer
		guid := c.Query("guid")
		newGUID := c.PostForm("guid")
		title := c.PostForm("title")
		content := c.PostForm("content")
		stmt, err := database.Prepare("UPDATE pages SET guid=?, title=?, content=? WHERE guid=?;")
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = stmt.Exec(newGUID, title, content, guid)
		if err != nil {
			fmt.Println(err.Error())
		}

		buffer.WriteString(title)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully updated to %s", name),
		})
	})

	router.DELETE("/api/page", func(c *gin.Context) {
		guid := c.Query("guid")
		stmt, err := database.Prepare("DELETE FROM pages WHERE guid=?")
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = stmt.Exec(guid)
		if err != nil {
			fmt.Println(err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully delete page: %s", guid),
		})
	})

	router.Run(PORT)
}
