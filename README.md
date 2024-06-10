# INSAkari Backend

## Description

Partie backend relié à l'application INSAkari. Gestion des comptes utilisateurs et des parties multijoueurs. 

## REST API

### Authentification

#### POST /auth/signup

Création d'un compte utilisateur

##### Paramètres

- ``username``: Nom d'utilisateur (unique) 
- ``email``: Adresse mail (unique)
- ``password``: Mot de passe répondant aux critères de sécurité :
  - Au moins 8 caractères
  - Au moins 1 minuscule
  - Au moins 1 majuscule
  - Au moins 1 caractère spécial
  - Au moins 1 chiffre

##### Réponse

- Succès

Code 201

- Erreur

Code 400 (Bad Request) ou 409 (Conflict)
```json
{
  "message": "Username not available" // message d'erreur
}
```

#### POST /auth/login

Connexion à un compte utilisateur

##### Paramètres
Soit :

- ``username``: Nom d'utilisateur
- ``password``: Mot de passe du compte associé au nom d'utilisateur

ou :
- ``token``: Token donné par l'API

##### Réponse

- Succès

Code 200
```json
{
  "token": "c9f12e7d-4aef-45cb-aab7-917e7f8eb613", // Token de connexion
  "user": { // Données utilisateur
      "id": 0,
	  "username": "admin",
	  "email": "admin@insakari.fr",
	  "score": 500
  }
}
```

- Erreur

Code 400 (Bad Request) ou 403 (Forbidden)
```json
{
  "message": "Invalid username or password" // message d'erreur
}
```

#### POST /auth/logout

Déconnexion d'un compte utilisateur

##### En-tête

- ``INSAkari-Connect-Token``: Token de connexion donné par l'API

##### Réponse

- Succès

Code 200

- Erreur

Code 401 (Unauthorized)
```json
{
  "message": "Invalid session" // message d'erreur
}
```

#### POST /auth/signout

Suppression d'un compte utilisateur

##### En-tête

- ``INSAkari-Connect-Token``: Token de connexion donné par l'API

##### Paramètres

- ``password``: Mot de passe du compte de l'utilisateur

##### Réponse

- Succès

Code 200

- Erreur

Code 400 (Bad Request) ou 401 (Unauthorized) ou 403 (Forbidden)
```json
{
  "message": "Invalid password" // message d'erreur
}
```

#### GET /auth/user

Récupération des informations de l'utilsateur

##### En-tête

- ``INSAkari-Connect-Token``: Token de connexion donné par l'API

##### Réponse

- Succès

Code 200
```json
{
  "user": {
      "id": 0,
	  "username": "admin",
	  "email": "admin@insakari.fr",
	  "score": 500
  }
}
```

- Erreur

Code 403 (Forbidden)
```json
{
  "message": "Incorrect token" // message d'erreur
}
```

#### POST /auth/user

Modification des informations de l'utilsateur

##### En-tête

- ``INSAkari-Connect-Token``: Token de connexion donné par l'API

##### Paramètres

- ``email``: Nouvelle adresse mail du compte de l'utilisateur
- ``username``: Nouveau nom d'utilisateur du compte de l'utilisateur
- ``password`` (facultatif): Nouveau mot de passe du compte de l'utilisateur

##### Réponse

- Succès

Code 200

- Erreur

Code 400 (Bad Request) ou 409 (Conflict) ou 401 (Unauthorized)
```json
{
  "message": "Username not available" // message d'erreur
}
```