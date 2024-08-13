# Setting up and Running the API

To set up and run this Gin Go REST API, follow these steps:

1. **Clone the repository:**
    ```bash
    git clone https://github.com/eulmlk/Go-Learning-Path.git
    ```

2. **Navigate to the project directory:**
    ```bash
    cd task_manager_mongodb
    ```

3. **Install the required dependencies:**
    ```bash
    go mod download
    ```

4. **Set up MongoDB:**
    - **Install MongoDB:**
      - Follow the instructions to install MongoDB on your system from the [official MongoDB documentation](https://docs.mongodb.com/manual/installation/).

    - **Start MongoDB:**
      - On Windows, you can start MongoDB by opening the MongoDB shell using `mongod.exe`.
      - On Linux/MacOS, you can run the following command to start MongoDB:
        ```
        sudo systemctl start mongod
        ```

    - **Use MongoDB Compass for Database Management:**
      - **Download and Install MongoDB Compass:**
        - Download MongoDB Compass from the [official MongoDB Compass download page](https://www.mongodb.com/try/download/compass).
      - **Connect to MongoDB:**
        - Open MongoDB Compass and connect to your local MongoDB server using the default connection string:
          ```
          mongodb://localhost:27017
          ```
      - **Create a Database and Collection:**
        - In MongoDB Compass, click on "Create Database," name your database `task_manager`, and create a collection named `tasks`.

    - **Set environment variables for MongoDB connection:**
      - Create a `.env` file in the root directory of your project and add the following variables:
        ```
        MONGODB_URI=mongodb://localhost:27017
        MONGODB_DB=task_manager
        ```

5. **Build the application:**
    ```bash
    go build -o app
    ```

6. **Run the application:**
    ```bash
    ./app
    ```

7. The API will be accessible at `http://localhost:8080`.

> **Note:** Please make sure you have Go installed and properly configured on your system before following these steps.

# API Documentation

For more details about the API endpoints and how to use them, please refer to the [API documentation](https://documenter.getpostman.com/view/33183582/2sA3rxpsfh).

To test the API using Postman, you will need to have Postman installed. You can import the Postman collection by clicking the "Run in Postman" button on the documentation page.