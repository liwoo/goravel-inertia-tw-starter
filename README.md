# Goravel Blog Application

A modern web application built with Goravel (Go) and React, featuring JWT authentication, dark mode, and a responsive UI.

## Getting Started

This guide will help you set up and run the application locally for development.

### Prerequisites

- Go 1.18 or higher
- Node.js 16 or higher
- NPM or Yarn
- MySQL or PostgreSQL database

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd blog
   ```

2. Install Go dependencies:
   ```bash
   go mod download
   ```

3. Install JavaScript dependencies:
   ```bash
   npm install
   # or
   yarn install
   ```

4. Configure your environment:
   - Copy `.env.example` to `.env`
   - Update database credentials and other settings in `.env`

5. Run database migrations:
   ```bash
   go run . artisan migrate
   ```

6. Generate application key:
   ```bash
   go run . artisan key:generate
   ```

### Creating an Admin User

To create an admin user, run:

```bash
go run . artisan user:create
```

Follow the prompts to enter email, name, and password.

### Running the Application

For development, you'll need to run two processes in separate terminal windows:

#### Terminal 1: Backend Server (with Air for hot-reload)

```bash
# If you don't have Air installed:
go install github.com/cosmtrek/air@latest

# Run the backend server with hot-reload
air
```

#### Terminal 2: Frontend Assets (with Vite)

```bash
# Run Vite development server
npm run dev
# or
yarn dev
```

The application should now be running at `http://localhost:8000`

## Project Structure

### Backend Structure

```
├── app/                  # Application code
│   ├── http/             # HTTP layer (controllers, middleware)
│   │   ├── controllers/  # Request handlers
│   │   └── middleware/   # HTTP middleware
│   ├── models/           # Database models
│   └── providers/        # Service providers
├── bootstrap/            # Application bootstrap code
├── config/               # Configuration files
├── database/             # Database migrations and seeds
├── public/               # Public assets
├── resources/            # Frontend resources
│   ├── css/              # CSS files
│   └── js/               # JavaScript/TypeScript files
├── routes/               # Route definitions
│   ├── api.go            # API routes
│   └── web.go            # Web routes
└── storage/              # Storage for logs, cache, etc.
```

### Frontend Structure

```
resources/
├── css/
│   └── app.css           # Global CSS with theme variables
└── js/
    ├── app.tsx           # Main React entry point
    ├── components/       # Reusable UI components
    │   ├── ThemeToggle.tsx       # Dark/light mode toggle
    │   └── app-sidebar.tsx       # Sidebar with theme toggle
    ├── context/
    │   └── ThemeContext.tsx      # Theme context provider
    ├── layouts/          # Page layouts
    │   ├── Admin.tsx     # Layout for authenticated users
    │   └── Auth.tsx      # Layout for login/register pages
    └── pages/            # Page components
        ├── auth/         # Authentication pages
        └── dashboard/    # Dashboard pages
```

## Key Features

### Authentication

- JWT-based authentication
- Login, registration, and logout functionality
- Protected routes with middleware

### Layouts

The application uses two main layouts:

1. **Auth Layout** (`resources/js/layouts/Auth.tsx`)
   - Used for login, registration, and password reset pages
   - Features a split-screen design with logo and theme toggle
   - Displays a form in the center of the screen

2. **Admin Layout** (`resources/js/layouts/Admin.tsx`)
   - Used for authenticated user pages
   - Features a sidebar navigation with collapsible menu
   - Includes user dropdown menu and theme toggle

### Routing

Routes are defined in two main files:

1. **Web Routes** (`routes/web.go`)
   - Public routes (login, register)
   - Protected routes that require authentication
   - Middleware groups for JWT authentication

2. **API Routes** (`routes/api.go`)
   - API endpoints for AJAX requests
   - Protected by JWT authentication

### Dark Mode

The application supports dark mode with:

- Theme toggle in both layouts
- CSS variables for theming in `resources/css/app.css`
- Theme state persisted in localStorage

## Goravel Fundamentals

Goravel is a web application framework for Go that follows Laravel's elegant syntax and architecture. Here are the core features used in this project:

### Routing

```go
// Define a simple route
facades.Route().Get("/", homeController.Index)

// Route with middleware
facades.Route().Middleware("auth:api").Get("/dashboard", dashboardController.Index)

// Route groups
facades.Route().Group(func(router route.Router) {
    router.Get("/users", userController.Index)
    router.Post("/users", userController.Store)
})
```

### ORM and Database

```go
// Query with relationships
user, err := facades.Orm().Query().With("Posts").First(&models.User{}, "id = ?", 1)

// Create a record
user := models.User{Name: "Goravel", Email: "goravel@example.com"}
facades.Orm().Query().Create(&user)

// Update a record
facades.Orm().Query().Model(&models.User{}).Where("id = ?", 1).Update("name", "New Name")
```

### Authentication

```go
// Authenticate a user
token, err := facades.Auth(ctx).Login(credentials)

// Get authenticated user
user, err := facades.Auth(ctx).User()

// Logout
err := facades.Auth(ctx).Logout()
```

### Middleware

```go
// Define middleware
func JwtAuth() http.Middleware {
    return func(ctx http.Context) {
        // Authentication logic
        ctx.Request().Next()
    }
}

// Apply middleware to routes
facades.Route().Middleware("jwt.auth").Get("/dashboard", controller.Dashboard)
```

### Artisan Commands

```go
// Define a command
type UserCreate struct {
    console.Command
}

// Register the command
facades.Schedule().Command("user:create").EveryMinute()

// Run a command
go run . artisan user:create
```

### Logging

```go
// Log messages at different levels
facades.Log().Debug("Debug message")
facades.Log().Info("Info message")
facades.Log().Warning("Warning message")
facades.Log().Error("Error message")
```

### Configuration

```go
// Access configuration values
appName := facades.Config().GetString("app.name")
dbConnection := facades.Config().GetString("database.default")
```

### Validation

```go
// Validate request data
validator := validation.Make(map[string]any{
    "name": "Goravel",
    "email": "goravel@example.com",
})

validator.Rules(map[string]string{
    "name": "required|max:255",
    "email": "required|email",
})

if validator.Fails() {
    errors := validator.Errors()
    // Handle validation errors
}
```

## Troubleshooting

### Common Issues

1. **Database Connection Issues**
   - Verify database credentials in `.env`
   - Ensure database server is running

2. **JWT Authentication Issues**
   - Check that `JWT_SECRET` is set in `.env`
   - Verify token expiration settings

3. **Frontend Build Issues**
   - Clear node_modules and reinstall: `rm -rf node_modules && npm install`
   - Check for JavaScript errors in browser console

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Commit your changes: `git commit -m 'Add some feature'`
4. Push to the branch: `git push origin feature-name`
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
