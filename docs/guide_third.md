# Tutorial : first steps with OTO (part 3)

## OTO instance

The OTO instance is a struct with every required client such as :
- postgres client
- temporal client
- schema for parameters

In order to get one, you must indicates the .env file your already created in the previous tutorial.
Lets says we have the following .env :
```
POSTGRES_DB=oto-storage
POSTGRES_USER=admin
POSTGRES_PASSWORD=1234
POSTGRES_PORT=5432
POSTGRES_SEED=postgres

TEMPORAL_HOST=:7233
TEMPORAL_NAMESPACE=default
```

In your go project, you can now call an OTO instance :
```go
instance, err := oto.NewInstanceOto(".env")
if err != nil {
    return err
}
```

## OTO executables

OTO is handling the execution of executables such as tool like openSSL, nmap or your own bash script.

The concept of OTO is to set every parameters of executables in order to create command.

## Install tools

OTO is handling the execution of executables such as tool like openSSL, nmap or your own bash script.
In order to launch your first job you need to first, install the executable you want and second, register it in the database.

In the folder **data/** you have executables ready for OTO (openSSL, nmap and masscan).
You either have the choice between install one of those tools or your own.

In the [main.go](../main.go) file, you have several functions :
- fill_database
- test_temporal

You can call **fill_database** function in order to register executables like nmap, masscan or openSSL. 


