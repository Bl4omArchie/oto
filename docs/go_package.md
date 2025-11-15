# Tutorial - OTO with go package

## Get an OTO instance

This tutorial shows how to use OTO with the provided Go package.

- First, you need to get OTO with the go get command.
At the root of your project, enter the following :
```sh
go get github.com/Bl4omArchie/oto
```

- Now you want to import the package in to your go file.
Use the following import : 
```go
import (
    "github.com/Bl4omArchie/oto/pkg"
)
```

- Then you want to start a new OTO instance that will launch every needed services.
In your code do :
``` go
oto, err := oto.NewInstanceOto("db/storage.db")
if err != nil {
    fmt.Println(err)
}
```

Here you can modify the path of the database. OTO is currently supporting only sqlite database.
