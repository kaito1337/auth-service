﻿# Authentication service

To initialize the project, you need to enter:

```

go run .\cmd\main\main.go

```

To initialize the PostgreSQL, you need to enter:

```

docker run -d --name postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=root postgres:15.3

```

To initialize the MongoDB, you need to enter:

```

docker run -d --name mongodb -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=root mongo:4.4.22

```
