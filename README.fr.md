# Persona2

Gestionnaire de profils Git et de configuration

## Installation

Pour installer Persona2, utilisez la commande suivante :
```bash
go install github.com/yourusername/persona2/cmd/persona2@latest
```

Après l'installation, initialisez la configuration :
```bash
persona2 init
```

## Utilisation

### Initialiser la configuration
Pour initialiser la configuration :
```bash
persona2 init
```

Cette commande crée le fichier de configuration `~/.persona2.json` si celui-ci n'existe pas déjà.

### Ajouter un profil
Pour ajouter un nouveau profil :
```bash
persona2 add --profile=<nom-profil> --url=<url-repository>
```

### Lister les profils
Pour lister tous les profils disponibles :
```bash
persona2 list
```

### Changer de profil
Pour changer de profil :
```bash
persona2 switch --profile=<nom-profil>
```

### Supprimer un profil
Pour supprimer un profil :
```bash
persona2 remove --profile=<nom-profil>
```

### Mettre à jour un profil
Pour mettre à jour un profil depuis son repository :
```bash
persona2 update --profile=<nom-profil>
```

### Mettre à jour Persona2
Pour mettre à jour Persona2 vers la dernière version :
```bash
persona2 self-update
```
