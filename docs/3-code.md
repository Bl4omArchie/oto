# Tutorial : first steps with OTO (part 3)



## Prerequisite

OTO is handling the execution of executables such as tool like openSSL, nmap or your own bash script.
In order to launch your first job you need to first, install the executable you want and second, register it in the database.

In the folder **data/** you have executables ready for OTO (openSSL, nmap and masscan).
You either have the choice between install one of those tools or your own.

In the [main.go](../main.go) file, you have several functions :
- fill_database
- test_temporal

In order to run the demo, you must installed **openssl** and fill the database. Of course you can modify the script and use your executable, parameters and create your own command.


## Run the demo

At the root of the project, there is a file called [main.go](../main.go). 
This main package call a function called **launch_demo()** that will ingest into the database the parameters of openssl, create a command and a job.

Execute the script  :
```go
go run main.go
```

Normally you should see a key.pem file in the project folder. The job has been successfully executed.




