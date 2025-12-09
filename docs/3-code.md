# Tutorial : first steps with OTO (part 3)


## Prerequisite

OTO is handling the execution of executables such as tool like openSSL, nmap or your own bash script.

In order to run the demo, you must install **openssl** and fill the database. Of course you can modify the script and use your executable, parameters and create your own command and jobs.


## Run the demo

At the root of the project, there is a file called [main.go](../main.go). 
This main package call the function **launch_demo()**. The demo is about generating an RSA keypair of 2024 bits. Steps :

1. Ingest executable openssl version 3.5.3
2. Ingest parameters
3. Create a command called **GenRSA**, with the created parameters
4. Take each parameters of the command, add values and store them into a job called **GenRSA-2048**.
5. Run the job and store the key into key.pem

Execute the script  :
```go
go run main.go
```

You now should see a key.pem file in the project folder. The job has been successfully executed !
