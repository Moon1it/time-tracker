# Time Tracker Application

## Description

This project implements a time tracker with a RESTful API, which enables user to manage and track time of completing different tasks. The project includes filtering, pagination, and CRUD operations for user data.

## Test Assignment

This project was developed as a part of a test assignment for Effective Mobile for the position of Junior Golang Developer. The goal of the assignment was to demonstrate the ability to create a RESTful API for time tracking, including user management, filtering, pagination, and integration with an external API. The project was completed within one week and includes the following features:

- User CRUD operations
- Task time tracking
- Filtering and pagination
- Integration with an external API
- PostgreSQL database for data storage

## Getting Started

Follow these instructions to set up and run the project locally.

### Prerequisites

Ensure you have the following installed:
- Docker
- Docker Compose

### Setup

1. **Create .env File**
   - Copy the contents of the `app.example.env` file.
   - Create a new file named `.env` in the root directory.
   - Paste the copied contents into the `.env` file.
  
2. **Build and Run the Project**
   - Open a terminal and navigate to the root directory of the project.
   - Run the following command to build and start the project:
     ```bash
     docker-compose up
     ```
   - Docker Compose will handle building and starting the services defined in the `docker-compose.yml` file.

3. **Access the Documentation**
   - When the project is running, open your web browser and navigate to:
     ```
     http://localhost:8000/swagger/index.html
     ```
   - Explore the Swagger documentation to familiarize yourself with the functionality of the application.
