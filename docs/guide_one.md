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

Next part of the guide [here](guide_two.md)
