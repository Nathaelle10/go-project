
# **Résumé de l'application**

CMD : 
// launch image
docker build -f build/Dockerfile -t monitoring .g
docker run -p 8080:8080 -p 9090:9090 monitoring

// update image docker
docker build -f build/Dockerfile -t monitoring:latest .
docker run -p 8080:8080 -p 9090:9090 monitoring:latest



---

## 🔧 Nom :

**Dashboard Système**

---

## 🎯 Objectif principal :

Créer une application **web minimaliste et portable** qui permet de :

1. 📊 **Consulter en temps réel l'état du système** (CPU, mémoire, réseau, disque, etc.)
2. 🔍 **Scanner les ports ouverts** d'une adresse IP donnée
3. 🌐 **Afficher une interface graphique locale (dashboard web)** intégrée au binaire, sans dépendance externe

---

## ⚙️ Composants principaux

### 1. 🖥️ API HTTP (`localhost:9090`)

Fournit 2 endpoints JSON REST :

#### ➤ `/api/status`

Retourne l’état du système :

* Utilisation CPU, modèle et nombre de cœurs
* Mémoire vive & swap
* Informations disque (`/`)
* Interfaces réseau (I/O)
* Charge moyenne du système
* Infos hôte (hostname, uptime, etc.)
* Nombre de connexions TCP actives

#### ➤ `/api/scan?ip=<adresse_ip>`

Fait un **scan rapide des ports TCP 20 à 1024** sur l’IP fournie.
Réponse JSON contenant les **ports ouverts**.

---

### 2. 🧭 Interface Web (`localhost:8080`)

Serveur HTTP qui héberge un dashboard HTML/CSS/JS avec :

* Une **UI moderne** grâce à [Shoelace UI](https://shoelace.style/) (ou CSS local si offline)
* Des **graphiques dynamiques** via Chart.js
* Un **formulaire pour scanner une IP**
* Des **cartes d'information** (CPU, mémoire, etc.)

👉 Tous les fichiers `/www` (HTML, JS, CSS) sont **intégrés directement dans le binaire** grâce à `embed`.

---

## 📦 Packaging

L'application peut être :

* ✅ Compilée en **binaire autonome multiplateforme** (Linux, Windows, macOS)
* 🐳 Intégrée dans un **Docker `FROM scratch`** (aucune dépendance système)
* 📁 Poussée sur GitHub avec arborescence propre
* ☁️ Facilement déployée ou distribuée

---

## ✅ Fonctionnalités clés

| Fonction                     | Détail                                                        |
| ---------------------------- | ------------------------------------------------------------- |
| 💻 API système               | Expose les métriques matérielles du serveur local             |
| 🔍 Scan de ports             | Permet de tester si des ports sont ouverts sur une IP donnée  |
| 🌍 UI embarquée              | Interface graphique intégrée directement dans le binaire      |
| 🚀 Binaire statique portable | Aucun besoin de dépendance ou serveur web externe             |
| 🌐 Middleware CORS           | Autorise appels cross-origin pour UI + API                    |
| 🔄 Multi-serveur             | Un serveur API (9090) + un serveur UI (8080), lancés ensemble |

---

## 📁 Structure du code (extrait simplifié)

```
Makefile 
go.mod 
go.sum 
build/ 
  ├── Dockerfile 
  ├── app-linux 
  └──app-windows.exe 
src/ 
  ├── main.go          // Lance les deux serveurs 
  │
  ├── api/             // Tous les handlers API 
  │ ├── status.go      // /api/status 
  │ ├── scan.go        // /api/scan 
  │ ├── middleware.go  // CORS 
  │ ├── routes.go      // Serveur API 
  │ └── server.go      // Serveur web statique 
  │ 
  └── web/ 
    ├── embed.go       // embed des fichiers www 
    └── www/           // Fichiers front embarqués 
    ├── index.html 
    └── assets/ 
        ├── css/style.css 
        └── js/dashboard.js

```


```
Makefile 
build/ 
  ├── Dockerfile 
  ├── app-linux 
  └──app-windows.exe 
src/ 
  ├── main.go          // Lance les deux serveurs 
  ├── go.mod          
  ├── go.sum          
  ├── status.go        // /api/status 
  ├── scan.go          // /api/scan 
  ├── middleware.go    // CORS 
  ├── routes.go        // Serveur API 
  ├── server.go        // Serveur web statique  
  └── www/             // Fichiers front embarqués 
        ├── index.html 
        └── assets/ 
            ├── css/style.css 
            └── js/dashboard.js

```

---

## 🛠️ Technologies utilisées

| Langage         | Utilisation                                        |
| --------------- | -------------------------------------------------- |
| **Go**          | Back-end, serveur HTTP, scanner de ports           |
| **HTML/CSS/JS** | UI dashboard                                       |
| **Chart.js**    | Graphiques système dynamiques (CPU, RAM, etc.)     |
| **Shoelace UI** | Composants visuels UI modernes                     |
| **gopsutil**    | Récupération des infos système                     |
| **embed**       | Intégration des fichiers statiques dans le binaire |

---



Oui, tout à fait ! Tu peux **passer un message de commit en argument** à une commande `make`, c’est très pratique pour automatiser Git avec des messages dynamiques.

---

## ✅ Étape 1 : Ajoute une règle `git-push` dans ton `Makefile`

Voici la version basique :

```makefile
git-push:
	git add .
	git commit -m "$(m)"
	git push
```

---

## ✅ Étape 2 : Utilisation dans le terminal

Tu passes le message avec `m="ton message"` comme ceci :

```bash
make git-push m="c'est mon message de commit"
```

---

## ✅ Bonus : Ajoute de la sécurité & du style (facultatif)

Tu peux améliorer la commande avec un check si le message est vide :

```makefile
git-push:
	@if [ -z "$(m)" ]; then \
		echo "❌ Merci de fournir un message de commit avec m=\"ton message\""; \
		exit 1; \
	fi
	@git add .
	@git commit -m "$(m)"
	@git push
```

---

## ✅ Exemple complet

Ajoute à ton `Makefile` :

```makefile
git-push:
	@if [ -z "$(m)" ]; then \
		echo "❌ Merci de fournir un message de commit avec m=\"ton message\""; \
		exit 1; \
	fi
	@echo "📦 Commit et push avec message : '$(m)'"
	@git add .
	@git commit -m "$(m)"
	@git push
```

Et utilise-le comme ça :

```bash
make git-push m="✨ Nouvelle version avec Docker et Makefile"
```

---

Souhaite-tu chaîner ça avec `make auto` pour faire `build + git-push` dans une seule commande ?
