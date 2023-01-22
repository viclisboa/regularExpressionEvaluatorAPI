# regularExpressionEvaluatorAPI

### Tests
 To run the tests use the following command
```shell
 go test -cover ./...
```

### Running the application
```shell
 docker-compose up
```

The application is using basic auth, the user and password are testeUser and testePassword. The application is running on port 808, to access you should use http://localhost:8080/expressions (example url used to recover all expressions in database)