package handler

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	// Optional for AWS Lambda-style handlers, but Vercel doesn't always require this.
)

type ShortLink struct {
	Id  string
	Url string
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var linkmap = map[string]ShortLink{"example": {Id: "example", Url: "http://example.com"}}

// Handler function that Vercel looks for to handle requests
func Handler(w http.ResponseWriter, r *http.Request) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.GET("/:id", RedirectHandler)
	e.POST("/", IndexHandler)
	e.POST("/submit", SubmitHandler)

	e.ServeHTTP(w, r)
}

func GenerateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	var result []byte
	for i := 0; i < length; i++ {
		index := seededRand.Intn(len(charset))
		result = append(result, charset[index])
	}
	return string(result)
}

func RedirectHandler(c echo.Context) error {
	id := c.Param("id")
	link, found := linkmap[id]

	if !found {
		return c.String(http.StatusNotFound, "Link not found")
	}

	return c.Redirect(http.StatusMovedPermanently, link.Url)
}

func IndexHandler(c echo.Context) error {
	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Short Links</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f4f4f4;
				margin: 0;
				padding: 0;
				display: flex;
				flex-direction: column;
				align-items: center;
				justify-content: center;
				height: 100vh;
			}
			h1 {
				color: #333;
			}
			form {
				margin-bottom: 20px;
				background-color: #fff;
				padding: 20px;
				border-radius: 8px;
				box-shadow: 0 2px 10px rgba(0,0,0,0.1);
			}
			label {
				margin-right: 10px;
			}
			input[type="text"] {
				padding: 8px;
				border-radius: 4px;
				border: 1px solid #ccc;
				width: 250px;
			}
			input[type="submit"] {
				padding: 8px 16px;
				border: none;
				border-radius: 4px;
				background-color: #28a745;
				color: white;
				cursor: pointer;
				margin-left: 10px;
			}
			input[type="submit"]:hover {
				background-color: #218838;
			}
			ul {
				list-style-type: none;
				padding: 0;
			}
			ul li {
				margin: 5px 0;
			}
			ul li a {
				color: #007bff;
				text-decoration: none;
			}
			ul li a:hover {
				text-decoration: underline;
			}
		</style>
	</head>
	<body>
		<h1>Submit a new website</h1>
		<form action="/submit" method="POST">
			<label for="url">Website URL:</label>
			<input type="text" id="url" name="url" placeholder="https://example.com">
			<input type="submit" value="Submit">
		</form>
		<h2>Existing Links</h2>
		<ul>`
	for _, link := range linkmap {
		html += `<li><a href="/` + link.Id + `">` + link.Id + `</a></li>`
	}
	html += `
		</ul>
	</body>
	</html>`

	return c.HTML(http.StatusOK, html)
}

func SubmitHandler(c echo.Context) error {
	url := c.FormValue("url")

	if url == "" {
		return c.String(http.StatusBadRequest, "URL is required")
	}
	if !(len(url) >= 4 && (url[:4] == "http" || url[:5] == "https")) {
		url = "https://" + url
	}

	id := GenerateRandomString(8)

	linkmap[id] = ShortLink{Id: id, Url: url}

	return c.Redirect(http.StatusSeeOther, "/")
}
