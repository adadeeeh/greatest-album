# greatest-album

# Steps

1. Download this [csv](https://www.kaggle.com/notgibs/500-greatest-albums-of-all-time-rolling-stone/version/1) and convert it to JSON.
2. Create **.env** file in **app** folder. Write the **MONGOURI** environment. Example if not using docker-compose:
   ```
   MONGOURI=mongodb://admin:password@localhost:27017
   ```
   Example if using docker-compose:
   ```
   MONGOURI=mongodb://admin:password@name_of_the_mongodb_service
   ```
3. Initialize Go module `go mod init greatest-album`
4. Install required packages `go get -u github.com/gin-gonic/gin go.mongodb.org/mongo-driver/mongo github.com/joho/godotenv github.com/go-playground/validator/v10`

   **github.com/gin-gonic/gin** is a framework for building web applications.

   **go.mongodb.org/mongo-driver/mongo** is a driver for connecting to MongoDB.

   **github.com/joho/godotenv** is a library for managing environment variables.

   **github.com/go-playground/validator/v10** is a library for validating structs and fields.

5. Create a docker image named **webapp:multistage** by running command `docker build -t webapp:multistage .`
6. Run `docker-compose up -d`
7. Via Mongo Express, create new collection named **album** inside **greatest-album** database. Import file or add data via new document and paste the whole JSON.
8. Create collection named **account** inside **greatest-album** database. And create one document with key Username and Password.
9. Start Gin `go run main.go`

# API

1. Get all album ![](images/get-all-album.png)
2. Get specific album ![](images/get-album-by-ranking-number.png)
3. Create album ![](images/create-album.png)
4. Delete specific album ![](images/delete-album-by-ranking-number.png)

# References

1. https://pkg.go.dev/github.com/gin-gonic/gin
2. https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo
3. https://docs.mongodb.com/drivers/go/current/fundamentals/crud/write-operations/upsert/
4. https://dev.to/hackmamba/build-a-rest-api-with-golang-and-mongodb-gin-gonic-version-269m
5. https://jonathanmh.com/go-gin-http-basic-auth/
