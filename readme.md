
# Env Manager 0.0.22

Manages multiple `.env` files configurations for your project.
Configurations are stored encrypted.

## Install
Run `install.sh` from your system

## Add a configuration
In order to store a new configuration file you need to add an header to your file to
specify the `identifier` of that configuration

_example_
```
#- identifier: LOCAL
STAGE=LOCAL
PEM_KEY=~/Downloads/MyKey.pem
SECRET_KEY=~/Downloads/MyKey.pem
```

You also need to generate a secret to encrypt your configurations.
```
openssl rand -hex 32 > .secret
```

Currently this is the only way to pass the secret.
Eventaully the secret will be read from environment and from cli argument.

```
env-manager add -f .env
```

This command will save the current `.env` enviroment using the `identifier` specified
in the headers

In order to restore it

```
env-manager get -i <header identifer>
```
