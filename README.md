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
