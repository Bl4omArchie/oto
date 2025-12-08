# Tutorial : first steps with OTO (part 1)

# What is OTO ?

OTO is a go service that orchestrate executables. An executable can be a tool, a custom bash script etc.

For each executables, you have a set of parameters.

$$ 
E = \{a, b, c, d, e \}
$$

From those parameters you can create commands.

$$
EC1 = \{a, b\}
$$

From those commands you can defines jobs which are commands with pre-filled values.

$$
JC1 = \{a: 192.168.0.1, b: true\}
$$

Then, you can orchestrate those jobs with temporal by using the scheduling feature.

```
JC1 -> JC2 -> JC1
    -> JC3 
```

# Generate an RSA keypair with openSSL

Lets generate an RSA keypair with **openSSL** as an example.

### 1. Define an executable

You have an executable called openSSL version 3.5.3 defined as follow :
```go
var openssl Executable {
    Tag : "openssl - 3.5.3"
    Name: "openssl"
    Version: "3.5.3"
    Path: "/usr/bin/openssl"
    Descriptio: "toolkit for general-purpose cryptography and secure communication"
}
```
We want to execute this command :
```bash
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out key.pem
```

### 2. Define parameters

To do so, we want to ingest the following parameters : **genpkey** **-algorithm**, **-pkeyopt**, **rsa_keygen_bits:2048**, and **-out**.

We defined them as follow :
```go
var param1 Parameter = {Flag: "genpkey", Description: "generate keypair", ExecutableTag: "openssl - 3.5.3", RequiresRoot: false, RequiresValue: false, ValueType: None}

var param2 Parameter = {Flag: "-algorithm", Description: "select a cryptosystem", ExecutableTag: "openssl - 3.5.3", RequiresRoot: false, RequiresValue: true, ValueType: String}

var param3 Parameter = {Flag: "-pkeyopt", Description: "define keypair options (i.e : size in bits)", ExecutableTag: "openssl - 3.5.3", RequiresRoot: false, RequiresValue: true, ValueType: String}

var param4 Parameter = {Flag: "-out", Description: "filepath for key storage", ExecutableTag: "openssl - 3.5.3", RequiresRoot: false, RequiresValue: true, ValueType: FilePath}
```

### 3. Define a command

Now we take our parameters and build our command :
```go
var cmd Command = {
    Name: "GenRSA"
    Description: "Generate an rsa keypair"
    ExecutableTag: "openssl - 3.5.3"
    RequiresRoot: false
    Parameters: {param1, param2, param3, param4}   
}
```
Okay great, now we have a command.

But I also want to define different usages for one command with different setting options, like several key size.

This is where the **jobs** are coming in.

### 4. Define a job

A job is a command with predefined values. A job is the actual command OTO will run later.
```go
var job Job = {
    Name: "GenRSA-2048"
    Description: "Generate an rsa keypair"
    CommandId: "GenRSA" 
    Parameters: map[Parameter]any{param1: "", param2: "RSA", param3: "rsa_keygen_bits:2048", param4: "/path/for/my/keypair.pem"}
}
```

Here as you can see, we define the algorithm we want, the key sizes with rsa_keygen_bits:2048 and the path where to store the key.

### 5. Run the job

Currently, the running system for jobs isn't 100% done because Temporal isn't completly integrated. But theoritically, you would have to call RunJobWorkflow() function like this :

```go
output, err := RunJobDemo(ctx context.Context, "GenRSA-2048")
```

Then in output, you'll find Stdout and Stderr.


Next part of the guide [here](2-services.md)
