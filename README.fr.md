# Persona

Gestionnaire de profils Git et de configuration

## Installation

Pour installer Persona, utilisez la commande suivante :
```bash
go install github.com/yourusername/persona/cmd/persona@latest
```

Après l'installation, initialisez la configuration :
```bash
persona init
```

## Utilisation

### Initialiser la configuration
Pour initialiser la configuration :
```bash
persona init
```

Cette commande crée le fichier de configuration `~/.persona.json` si celui-ci n'existe pas déjà. Elle crée également le dossier `~/.persona/profiles` où seront stockés tous les profils.

### Ajouter un profil
Pour ajouter un nouveau profil :
```bash
persona add --url=<url-repository>
```

Le nom du profil peut être extrait de l'URL si non fourni. Le repository sera cloné dans `~/.persona/profiles/<nom-profil>`.

### Lister les profils
Pour lister tous les profils disponibles :
```bash
persona list
```

Cette commande affiche tous les profils disponibles et leur statut actuel (actif ou inactif).

### Changer de profil
Pour changer de profil :
```bash
persona switch --profile=<nom-profil>
```

Cette commande crée des liens symboliques vers les fichiers de configuration du profil sélectionné. Si des fichiers existent déjà dans l'emplacement cible, la commande demandera une confirmation avant de les remplacer.

### Supprimer un profil
Pour supprimer un profil :
```bash
persona remove --profile=<nom-profil>
```

Cette commande supprime le profil et ses liens symboliques. Elle demandera une confirmation avant de procéder.

### Mettre à jour un profil
Pour mettre à jour un profil depuis son repository :
```bash
persona update --profile=<nom-profil>
```

Pour mettre à jour tous les profils :
```bash
persona update
```

Cette commande récupère les dernières modifications depuis le repository et met à jour les fichiers locaux.

### Mettre à jour Persona
Pour mettre à jour Persona vers la dernière version :
```bash
persona self-update
```

Cette commande télécharge et installe la dernière version de Persona depuis GitHub.

## Configuration

Le fichier de configuration est situé à `~/.persona.json` et contient :
- Liste de tous les profils
- Profil actif
- Alias des profils
- Paramètres de configuration

## Structure des profils

Chaque profil est stocké dans `~/.persona/profiles/<nom-profil>` et doit contenir :
- Fichiers de configuration
- Dotfiles
- Autres ressources de configuration

Lors du changement de profil, Persona crée des liens symboliques des fichiers du profil vers leurs emplacements appropriés dans votre dossier personnel.

## Support

Langues disponibles :
- Anglais
- Français
- Allemand

L'application supporte l'internationalisation et peut être utilisée dans plusieurs langues.
