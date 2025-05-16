# Persona2

Git and Configuration Profile Manager

[Documentation française](README.fr.md)

## Installation

After installation, initialize the configuration:
```bash
persona init
```

## Usage

### Initialize Configuration
To initialize the configuration:
```bash
persona init
```

This command creates the configuration file `~/.persona.json` if it doesn't exist yet. It also creates the directory `~/.persona/profiles` where all profiles will be stored.

### Add a Profile
To add a new profile:
```bash
persona add --url=<repository-url>
```

The profile name can be extracted from the URL if not provided. The repository will be cloned into `~/.persona/profiles/<profile-name>`.

### List Profiles
To list all available profiles:
```bash
persona list
```

This command shows all available profiles and their current status (active or inactive).

### Switch Profile
To switch to a different profile:
```bash
persona switch --profile=<profile-name>
```

This command creates symbolic links to the configuration files of the selected profile. If files already exist in the target location, the command will prompt for confirmation before replacing them.

### Remove Profile
To remove a profile:
```bash
persona remove --profile=<profile-name>
```

This command removes the profile and its symbolic links. It will prompt for confirmation before proceeding.

### Update Profile
To update a profile from its repository:
```bash
persona update --profile=<profile-name>
```

To update all profiles:
```bash
persona2 update
```

This command fetches the latest changes from the repository and updates the local files.

### Self-Update
To update Persona2 to the latest version:
```bash
persona self-update
```

This command downloads and installs the latest version of Persona2 from GitHub.

## Configuration

The configuration file is located at `~/.persona.json` and contains:
- List of all profiles
- Current active profile
- Profile aliases
- Configuration settings

## Profile Structure

Each profile is stored in `~/.persona/profiles/<profile-name>` and should contain:
- Configuration files
- Dotfiles
- Any other configuration resources

When switching profiles, Persona2 creates symbolic links from the profile files to their appropriate locations in your home directory.

## Support

Available languages:
- English
- French
- German

The application supports internationalization and can be used in multiple languages.
