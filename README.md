# Implementing Authentication and Authorization with JWT for Task Management API
## Objective
  - The objective of this task is to enhance the Task Management API by adding authentication and authorization mechanisms using JSON Web Tokens (JWT). This enhancement will introduce the concept of users, login functionality, and protected routes to restrict access to certain endpoints based on user roles.

## Instructions
  - Implement user management endpoints for user registration and login, including:
    - POST /register: Create a new user account with a unique username and password.
    - POST /login: Authenticate users and generate JWT tokens upon successful login.
  - Generate JWT tokens with appropriate claims (e.g., user ID, username, expiration time) using a secure JWT library.
  - Implement middleware to validate JWT tokens for protected routes, ensuring that only authenticated users can access them.
  - Define user roles and restrict access to certain endpoints based on user roles using middleware.
  - Update existing API endpoints to enforce authentication and authorization requirements for protected routes.
  - Test the API endpoints using Postman or similar tools to verify that authentication and authorization are working correctly.
  - Verify that only authenticated users can access protected routes, and unauthorized access attempts are rejected with appropriate error responses.
  - Ensure that user credentials are securely stored using appropriate encryption and hashing techniques to protect against security threats.
  - Document the authentication and authorization process, including instructions for user registration, login, and usage of protected endpoints.
  - Update the API documentation to reflect changes related to authentication and authorization, including any modifications to request and response formats.

## Folder Structure
  - Follow the following folder structure for this task
    ```
    task_manager/
    ├── main.go
    ├── controllers/
    │   └── controller.go
    ├── models/
    │   ├── task.go
    │   └── user.go
    ├── data/
    │   ├── task_service.go
    │   └── user_service.go
    ├── middleware/
    │   └── auth_middleware.go
    ├── router/
    │   └── router.go
    ├── docs/
    │   └── api_documentation.md
    └── go.mod
    ```
    - **main.go**: Entry point of the application.
    - **controllers/controller.go**: Handles incoming HTTP requests and invokes the appropriate service methods for both tasks and user authentication.
    - **models/task.go**: Defines the Task struct.
    - **models/user.go**: Defines the User struct.
    - **data/task_service.go**: Contains business logic and data manipulation functions for tasks.
    - **data/user_service.go**: Contains business logic and data manipulation functions for users, including hashing passwords.
    - **middleware/auth_middleware.go**: Implements middleware to validate JWT tokens for authentication and authorization.
    - **router/router.go**: Sets up the routes and initializes the Gin router and defines the routing configuration for the API.
    - **docs/api_documentation.md**: Contains API documentation and other related documentation.
    - **go.mod**: Defines the module and its dependencies.
  - **Note**:
    - Authentication and authorization are critical components of web applications to ensure secure access to resources. Pay close attention to implementing these mechanisms securely and effectively.
