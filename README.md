# Example

### Generate short url
`curl -H "Content-Type: application/json" -X POST -d '{"url":"http://www.google.com"}' http://localhost:8080/encode/`

Response: `{"success":true,"response":"http://localhost:8080/1"}`

### Info
Visit in browser `http://localhost:8080/info/1`



# Installation
### gulp
```bash
npm install
npm i -g gulp-cli

gulp watch
```

### RUN
```
go run main.go
```
