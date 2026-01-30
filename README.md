📚 Library API (Pet Project)

Mini pet-project for practicing Go backend development.
A simple HTTP API for managing a library — without frontend, tested via Postman.

The project allows you to add books, remove them, take and return books, and view available or taken ones.
All data is stored in memory (no database).

🚀 Features

Add a book

Delete a book

Get all books

Get available books

Get taken books

Take a book

Return a book

⭐ Download books list as a .txt file

🛠 Tech Stack

Go

net/http (or Gin)

JSON

Postman

📦 Book Model

Example:

{
"id": 1,
"title": "Clean Code",
"author": "Robert Martin",
"taken": false,
"takenBy": ""
}

📡 API Endpoints
➕ Add book

POST /books

{
"title": "Go in Action",
"author": "William Kennedy"
}

❌ Delete book

DELETE /books/{id}

📚 Get all books

GET /books

✅ Get available books

GET /books/available

🔒 Get taken books

GET /books/taken

🙋 Take book

POST /books/{id}/take

{
"user": "Alex"
}

🔁 Return book

POST /books/{id}/return

⭐ Download books list

GET /books/download

Returns a .txt file with all books:

[1] Clean Code - Robert Martin - AVAILABLE
[2] Go in Action - William Kennedy - TAKEN by Alex

▶️ Run Project
go mod init library-api
go run .


or

go run main.go

🧠 Rules

Data is stored in memory

IDs are auto-generated

No database

Data resets after restart

API returns JSON

🎯 Goal

This project is made to practice:

HTTP handlers

REST API

JSON encode/decode

Go structures

Business logic

File download

API design

😄 Future Ideas

Save data to file

Add auth

Add filters

Add tests

Dockerize

💪 Author

Pet project for Go practice 🚀