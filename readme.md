
# Env Manager 0.0.24

Manages multiple `.env` files configurations for your project.
Configurations are stored encrypted.

## Install
Download and install installation script (Only tested on linux)

```bash
 curl -sSL https://raw.githubusercontent.com/vertefra/env-manager/master/install.sh | bash
```

## Add a configuration
In order to store a new configuration file you need to add an `header` to your file.
Headers are at the top of the file and are metadata about the file.
The **identifier** identify the configuration when you want to restore it.

_example_
```
#- identifier: LOCAL
#- restore-as: .env
STAGE=LOCAL
PEM_KEY=~/Downloads/MyKey.pem
SECRET_KEY=~/Downloads/MyKey.pem
```

Another `header` is `restore-as` which is the name of the file that will be restored. If not provided the file will be restored as `.env`

## Secrets
In order to encrypt and decrypt the configurations you need to generate a secret.
Secret is read first from the environment variable `ENV_MANAGER_SECRET`. If no secret is found, it will try to look into a `.secret` file in the current directory.

> Valid secret is a 16, 24 or 32 bytes long string. If Generating from `hexadecimal` consider that 2 characters are 1 byte.

**Generate a secret in .secret file**
```bash
openssl rand -hex 16 > .secret
```

**Generate a secret in environment**
```bash
export ENV_MANAGER_SECRET=$(openssl rand -hex 16)
```

## Usage

Adding a new configuration
```
env-manager add -f .env
```

Where `.env` is the file you want to store and contains a valid `#- identifier: <config identifier>` header

In order to restore it

```
env-manager get -i <header identifer>
```

