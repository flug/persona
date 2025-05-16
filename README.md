# Persona2

Git and Configuration Profile Manager

[Documentation française](README.fr.md)

## Installation

To install Persona2, use the following command:
```bash
go install github.com/yourusername/persona2/cmd/persona2@latest
```

After installation, initialize the configuration:
```bash
persona2 init
```

## Usage

### Initialize Configuration
To initialize the configuration:
```bash
persona2 init
```

This command creates the configuration file `~/.persona2.json` if it doesn't exist yet.

### Add a Profile
To add a new profile:
```bash
persona2 add --profile=<profile-name> --url=<repository-url>
```

### List Profiles
To list all available profiles:
```bash
persona2 list
```

### Switch Profile
To switch to a different profile:
```bash
persona2 switch --profile=<profile-name>
```

### Remove Profile
To remove a profile:
```bash
persona2 remove --profile=<profile-name>
```

### Update Profile
To update a profile from its repository:
```bash
persona2 update --profile=<profile-name>
```

### Self-Update
To update Persona2 to the latest version:
```bash
persona2 self-update
```
