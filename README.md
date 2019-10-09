# Users
Users Microservice

Note: Currently this service intializes with an empty users database

Using the API:

List all users: `GET: api/v1/customers`

Register a new user: `POST: /api/v1/customers`

    *Request Body: username, password, email, firstName, lastName*
    
    
Check service health: `GET: /api/v1/health`

Prometheus Metrics: `GET: /api/v1/metrics`
