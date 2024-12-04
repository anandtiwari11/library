# Library Management System

A library management system built using **Go** and **GORM** ORM, designed to manage users, books, and their borrowing processes. This API enables the following functionalities:

- Manage users (create, retrieve)
- Manage books (add, retrieve, borrow)
- Manage authors (create, list books by author)

## Features

- **User Management**: Create and retrieve user details.
- **Book Management**: Add new books, retrieve books, and borrow books.
- **Author Management**: Add new authors and view books written by a specific author.
- **Borrowing System**: Allows users to borrow books, ensuring books are available before borrowing.

## Database Structure

- **Users**: Information about users including name, email, subscription end date, and borrowed books.
- **Books**: Information about books including title, author, availability status, and publication details.
- **Authors**: Information about authors including name and books written.

### Database Models

- **User**: Contains user details like name, email, subscription end date, and a list of borrowed books.
- **Book**: Contains book details like title, availability, and its author.
- **Author**: Contains details about authors including name and associated books.

## Technologies Used

- **Go**: The programming language used for backend development.
- **GORM**: ORM for interacting with the SQLite database.
- **SQLite**: The database used to store the data.
- **Gin**: Web framework used for building the RESTful API.
  
## Setup and Installation

### Prerequisites

- Go 1.18+
- SQLite (optional: can be installed by Go)
- A GitHub account to clone the repository

### Installation Steps

1. **Clone the repository**:

   ```bash
   git clone https://github.com/anandtiwari11/library-management-system.git
   cd library-management-system
   

### Notes:

- Feel free to modify the **installation** and **setup** sections based on your actual setup and environment.
- Ensure your database models and routes match the endpoints mentioned in the `README.md`.

