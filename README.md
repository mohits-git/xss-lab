# XSS-Lab

XSS-Lab is a simple blog website built with Go. There are many vulnerabilities left in the website to show case XSS attacks and their preventions.

### Draw.io diagram
[XSS-Lab Diagram](https://drive.google.com/file/d/1EpMI7AovdUgfQc01HOTXiwcdAmHypDSo/view?usp=sharing)

## Technologies Used
- Go (Golang)
- HTML/CSS
- JavaScript
- SQLite (for database)
- Goose (for db migrations)
### go dependencies
- sqlite3 driver [github.com/ncruces/go-sqlite3/](github.com/ncruces/go-sqlite3/)
- jwt [github.com/golang-jwt/jwt](github.com/golang-jwt/jwt)
- bcrypt [golang.org/x/crypto/bcrypt](golang.org/x/crypto/bcrypt)

## Project Setup
### Prerequisites
Goose
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```
SQLite3
```bash
sudo apt install sqlite3
```

1. Clone the repository:
```bash
git clone <repository-url>
cd xss-lab
```
2. Install dependencies:
```bash
go mod tidy
```
3. Set your environment variables:
```bash
export DB_URL=sqlite3.db
export JWT_SECRET=your_jwt_secret
export PORT=8080
```
4. Run database migrations:
```bash
goose sqlite3 sqlite3.db up --dir sql/schema
```
5. Start the server:
```bash
go run main.go
```
6. Open your browser and go to `http://localhost:8080`
7. Register a new user or log in with an existing user.

## Directory Structure
```
xss-lab/
├── main.go                # Main application entry point
├── templates/             # HTML templates
├── static/                # Static files (CSS, JS, images)
├── sql/schema             # SQL schema files for database migrations
├── sql/migrations         # SQL migration files
├── internal/              # Internal packages
│   ├── database/          # Database connection and queries
│   ├── auth/              # Authentication - JWT, passwords hashing and auth headers
├── api_config.go          # API configuration and routes
├── handler_X.go           # HTTP handlers for routes
├── middleware/            # Middleware for protected routes
├── main.go                # Main application logic
├── go.mod                 # Go module file
├── go.sum                 # Go module dependencies
└── README.md              # Project documentation
```

## Features
- User registration and login
- Create blog posts (logged in users only)
- View all blog posts
- Search for blog posts by title
- View blog post and comments
- Comment on a blog post (logged in users only)

## Usage
- Register a new user by clicking on the "Register" link.
- Log in with your credentials.
- Create a new blog post by clicking on the "Create Post" link (logged in users only).
- View all blog posts on the homepage.
- Search for blog posts by title using the search bar.
- Click on a blog post to view its details and comments.
- Comment on a blog post by filling out the comment form (logged in users only).

## Security Considerations
- The application uses JWT for user authentication.
- Passwords are hashed using bcrypt before storing them in the database.
- NO Input validation and sanitization are implemented to prevent XSS attacks.
- We are using unsafe `text/template` package to render user input directly in the HTML without sanitization.
- The application does not implement any Content Security Policy (CSP) headers, which could help mitigate XSS attacks.

## XSS Vulnerabilities
The application intentionally contains XSS vulnerabilities to demonstrate how they can be exploited.
- The `Create Post` and `Comment` features do not sanitize user input, allowing for script injection.
- The `View Post` feature displays user-generated content without sanitization, making it vulnerable to XSS attacks.
- The `Search` feature does not sanitize the search input, allowing for potential XSS attacks in the search results. The search input is sent via query param and is rendered directly in the HTML without sanitization.
- The application stores the JWT cookie in local storage, which can be accessed by JavaScript, making it vulnerable to XSS attacks if an attacker can inject malicious scripts.
- The application does not implement any Content Security Policy (CSP) headers, which could help mitigate XSS attacks.

## Preventing XSS Attacks
To prevent XSS attacks, the following measures can be implemented:
- Use a library to sanitize user input before rendering it in the HTML.
- Use the `html/template` package instead of `text/template` to automatically escape HTML characters.
- Implement input validation and sanitization for all user inputs.
- Use Content Security Policy (CSP) headers to restrict the sources of scripts and other resources.

## API Endpoints
| Method | Endpoint                    | Description                              |
|--------|-----------------------------|------------------------------------------|
| POST   | /api/register               | Register a new user                      |
| POST   | /api/login                  | Log in a user                            |
| GET    | /api/blogs?query=title      | Search for blog posts by title           |
| GET    | /api/blogs/count            | Get the count of all blog posts          |
| GET    | /api/users/{id}/blogs       | Get all blogs by a user                  |
| POST   | /api/blogs                  | Create a new blog post (protected)       |
| POST   | /api/comments/{blog_id}     | Add a comment to a blog post (protected) |
| PUT    | /api/blogs/{id}             | Update a blog post (protected)           |
| DELETE | /api/blogs/{id}             | Delete a blog post (protected)           |


## Pages
| Page                         | Description                                     |
|------------------------------|-------------------------------------------------|
| Home (/)                     | Displays all blog posts                         |
| Register (/register)         | User registration page                          |
| Login (/login)               | User login page                                 |
| All Blogs (/blogs)           | Displays all blogs + Search + Create a new post |
| View Post (/blogs/{id})      | Displays a single blog and its comments         |


## XSS Attacks examples

- Reflected XSS attack via search query parameter:
  - Search for a blog post with a script tag in the title.
  - Example: `http://localhost:8080/blogs?query=<script>alert('XSS')</script>`

- Stored XSS attack via blog post creation or commenting on a famous blog post:
    - Create a blog post with a script tag in the content.
    - Example: `<script>alert('XSS')</script>`

    - Injecting a XSS Worm into the applicaion:
        - Create a blog post with a script tag that creates a new blog post with a script tag in the content.
        - Example:
        ```html
            <script>
                const token = localStorage.getItem('authToken')
                const formData = new FormData();
                formData.append('title', 'Make Money in Minutes');
                formData.append('content', "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Reiciendis suscipit possimus dolore. Quibusdam mollitia id accusamus consequatur ea molestiae eum vitae suscipit, voluptatibus rem eveniet, tempora necessitatibus aliquam voluptate alias? Lorem ipsum dolor sit amet, consectetur adipisicing elit. Reiciendis suscipit possimus dolore. Quibusdam mollitia id accusamus consequatur ea molestiae eum vitae suscipit, voluptatibus rem eveniet, tempora necessitatibus aliquam voluptate alias? Lorem ipsum dolor sit amet, consectetur adipisicing elit. Reiciendis suscipit possimus dolore. Quibusdam mollitia id accusamus consequatur ea molestiae eum vitae suscipit, voluptatibus rem eveniet, tempora necessitatibus aliquam voluptate alias?");
                fetch('/api/blogs', {
                	method: 'POST',
                	body: formData,
                	headers: {
                		'Authorization': `Bearer ${token}`,
                        }
                });
            </script>
        ```
