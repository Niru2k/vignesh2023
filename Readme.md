

## General view about the Job Site Project

This is a simple job site project built using the Echo web framework in Golang. 

The project provides RESTful APIs for user authentication, job posting, and user comments on job posts. 

The backend includes functionalities to sign up, login, post job details, view job postings, post and view user comments,

update and delete job postings and comments.

## Technology Stack

The Job Site Real Time Exercise Task project is built using the following technologies:

- **Golang**: The backend is written in Go (Golang), a statically typed, compiled language.

- **Echo**: The Echo web framework is used to create RESTful APIs and handle HTTP requests.

- **JWT**: JSON Web Tokens are used for secure user authentication and authorization.

- **bcrypt**: Passwords are stored securely in hashed form using the bcrypt hashing algorithm.

- **postgres**: PostgreSQL is an advanced, enterprise class open source relational database that supports both SQL and JSON  querying. 
                It is a highly stable database management system, which has contributed to its high levels of resilience,and correctness. 
   

## Setup

The application will be accessible at `http://localhost:8080`.

 ##  Project explanation

-> In this Job site Real time task it has signup and login api. While login up it generates a JWT token. With the help of this token we can do        
   authorization for user and admin.

-> In job posting API only admins are allowed to post the job details, and also only admin have the authorization to do update and delete the job post.

-> Only user and admin are authorized to view the job post details, get job post details using job id, get job post details using specific company name. 

-> Only user are given authorization to post comment on job post using job id, update comment and delete comment.

-> Only user and admin are authorized to view the job post comments, get comment using comment id.



The following API endpoints are available in the Job Site project:

- **POST /signup**: Register a new user account.

- **POST /login**: Log in with registered credentials and receive a JWT token.

- **POST /jobposting**: Post job details (Admin Authentication required).

- **GET /jobpostings**: Get all company job posting details (Admin and user authentication required ).

- **GET /jobpostings/:id**: Get job posting details by job ID (Admin and user authentication required).

- **PUT /jobpostings/:id**: Update job posting details by job ID (Admin Authentication required).

- **DELETE /jobpostings/:id**: Delete a job posting by job ID (Admin Authentication required).

- **GET /jobpostings/company/:companyname**: Get all job postings from a specific company ( Admin and user authentication required).

- **POST /usercomments**: Post a comment on a job post (User authentication required).

- **GET /usercomments**: Get all user comments (Admin and user authentication required).

- **GET /usercomments/:id**: Get a specific user comment by comment ID (Admin and user authentication required).

- **PUT /usercomments/:id**: Update a user comment by comment ID (User authentication required).

- **DELETE /usercomments/:id**: Delete a user comment by comment ID (User authentication required).

