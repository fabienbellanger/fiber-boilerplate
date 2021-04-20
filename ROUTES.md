# fiber-boilerplate
A simple boilerplate for [Fiber](https://github.com/gofiber/fiber)


## Routes

### Web

- **[GET] `/health-check`**: Check server
    ```bash
    http GET localhost:3000/health-check
    ```

- **[GET] `/metrics`**: Prometheus metrics
    ```bash
    http GET localhost:3000/metrics
    ```

### API

- **[POST] `/api/login`**: Authentication
    ```bash
    http POST localhost:3000/api/login username=test@gmail.com password=0000
    ```
    Response:
    ```json
    {
        "id": "2a40080f-6077-4273-9075-1c5503ac95eb",
        "username": "test@gmail.com",
        "lastname": "Test",
        "firstname": "Toto",
        "created_at": "2021-03-08T20:43:28.345Z",
        "updated_at": "2021-03-08T20:43:28.345Z",
        "token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTAzLTA4VDIwOjQzOjI4LjM0NVoiLCJleHAiOjE2MTYxMDAyMTUsImZpcnN0bmFtZSI6IkZhYmllbiIsImlhdCI6MTYxNTIzNjIxNSwiaWQiOjEsImxhc3RuYW1lIjoiQmVsbGFuZ2VyIiwibmJmIjoxNjE1MjM2MjE1LCJ1c2VybmFtZSI6InZhbGVudGlsQGdtYWlsLmNvbSJ9.RL_1C2tYqqkXowEi8Np-y3IH1qQLl8UVdFNWswcBcIOYB6W4T-L_RAkZeVK04wtsY4Hih2JE1KPcYqXnxj2FWg",
        "expires_at": "2021-03-18T21:43:35.641Z"
    }
    ```

- **[POST] `/api/register`**: User creation
    ```bash
    http POST localhost:3000/api/register lastname=Test firstname=Toto username=test@gmail.com password=0000
    ```
    Response:
    ```json
    {
        "id": "cb13cc29-13bb-4b84-bf30-17da00ec7400",
        "username": "test@gmail.com",
        "lastname": "Test",
        "firstname": "Toto",
        "created_at": "2021-03-09T21:05:35.564747+01:00",
        "updated_at": "2021-03-09T21:05:35.564747+01:00"
    }
    ```

- **[GET] `/v1/users`**: Users list
    ```bash
    http GET localhost:3000/api/v1/users "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTAzLTA4VDIwOjQzOjI4LjM0NVoiLCJleHAiOjE2MTYxMDAyMTUsImZpcnN0bmFtZSI6IkZhYmllbiIsImlhdCI6MTYxNTIzNjIxNSwiaWQiOjEsImxhc3RuYW1lIjoiQmVsbGFuZ2VyIiwibmJmIjoxNjE1MjM2MjE1LCJ1c2VybmFtZSI6InZhbGVudGlsQGdtYWlsLmNvbSJ9.RL_1C2tYqqkXowEi8Np-y3IH1qQLl8UVdFNWswcBcIOYB6W4T-L_RAkZeVK04wtsY4Hih2JE1KPcYqXnxj2FWg"
    ```
    Response:
    ```json
    [
        {
            "id": "2a40080f-6077-4273-9075-1c5503ac95ed",
            "username": "test@gmail.com",
            "lastname": "Test",
            "firstname": "Toto",
            "created_at": "2021-03-08T20:43:28.345Z",
            "updated_at": "2021-03-08T20:43:28.345Z"
        },
        {
            "id": "2a40080f-6077-4273-9075-1c5503ac95eb",
            "username": "test1@gmail.com",
            "lastname": "Test",
            "firstname": "Toto",
            "created_at": "2021-03-08T20:45:51.16Z",
            "updated_at": "2021-03-08T20:45:51.16Z"
        }
    ]
    ```

- **[GET] `/v1/users/{id}`**: Get user information
    ```bash
    http GET localhost:3000/api/v1/users/2a40080f-6077-4273-9075-1c5503ac95eb "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTAzLTA4VDIwOjQzOjI4LjM0NVoiLCJleHAiOjE2MTYxMDAyMTUsImZpcnN0bmFtZSI6IkZhYmllbiIsImlhdCI6MTYxNTIzNjIxNSwiaWQiOjEsImxhc3RuYW1lIjoiQmVsbGFuZ2VyIiwibmJmIjoxNjE1MjM2MjE1LCJ1c2VybmFtZSI6InZhbGVudGlsQGdtYWlsLmNvbSJ9.RL_1C2tYqqkXowEi8Np-y3IH1qQLl8UVdFNWswcBcIOYB6W4T-L_RAkZeVK04wtsY4Hih2JE1KPcYqXnxj2FWg"
    ```
    Response:
    ```json
    {
        "id": "2a40080f-6077-4273-9075-1c5503ac95eb",
        "username": "test@gmail.com",
        "lastname": "Test",
        "firstname": "Toto",
        "created_at": "2021-03-08T20:43:28.345Z",
        "updated_at": "2021-03-08T20:43:28.345Z"
    }
    ```

- **[DELETE] `/v1/users/{id}`**: Delete user
    ```bash
    http DELETE localhost:3000/api/v1/users/2a40080f-6077-4273-9075-1c5503ac95eb "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTAzLTA4VDIwOjQzOjI4LjM0NVoiLCJleHAiOjE2MTYxMDAyMTUsImZpcnN0bmFtZSI6IkZhYmllbiIsImlhdCI6MTYxNTIzNjIxNSwiaWQiOjEsImxhc3RuYW1lIjoiQmVsbGFuZ2VyIiwibmJmIjoxNjE1MjM2MjE1LCJ1c2VybmFtZSI6InZhbGVudGlsQGdtYWlsLmNvbSJ9.RL_1C2tYqqkXowEi8Np-y3IH1qQLl8UVdFNWswcBcIOYB6W4T-L_RAkZeVK04wtsY4Hih2JE1KPcYqXnxj2FWg"
    ```
  Response code `204`

- **[PUT] `/v1/users/{id}`**: Update user information
    ```bash
    http PUT localhost:3000/api/v1/users/2a40080f-6077-4273-9075-1c5503ac95eb "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTAzLTA4VDIwOjQzOjI4LjM0NVoiLCJleHAiOjE2MTYxMDAyMTUsImZpcnN0bmFtZSI6IkZhYmllbiIsImlhdCI6MTYxNTIzNjIxNSwiaWQiOjEsImxhc3RuYW1lIjoiQmVsbGFuZ2VyIiwibmJmIjoxNjE1MjM2MjE1LCJ1c2VybmFtZSI6InZhbGVudGlsQGdtYWlsLmNvbSJ9.RL_1C2tYqqkXowEi8Np-y3IH1qQLl8UVdFNWswcBcIOYB6W4T-L_RAkZeVK04wtsY4Hih2JE1KPcYqXnxj2FWg"  lastname=Test firstname=Tutu username=test3@gmail.com password=2222
    ```
  Response:
    ```json
    {
        "id": "2a40080f-6077-4273-9075-1c5503ac95eb",
        "username": "test3@gmail.com",
        "lastname": "Test",
        "firstname": "Tutu",
        "created_at": "2021-03-08T20:43:28.345Z",
        "updated_at": "2021-03-12T20:43:28.345Z"
    }
    ```
