# Dashboard Système

---

## 🚀 Résumé rapide

```bash
# Construire l'image Docker
docker build -f build/Dockerfile -t monitoring .

# Lancer le conteneur (API + UI)
docker run -p 8080:8080 -p 9090:9090 monitoring

# Mettre à jour l'image Docker (rebuild + relancer)
docker build -f build/Dockerfile -t monitoring:latest .
docker run -p 8080:8080 -p 9090:9090 monitoring:latest
```

---

## 🔧 Nom de l’application

**Dashboard Système**

---

## 🎯 Objectif principal

Une application **web minimaliste et portable** pour :

1. 📊 Consulter en temps réel l’état du système (CPU, mémoire, réseau, disque…)
2. 🔍 Scanner les ports ouverts sur une adresse IP donnée
3. 🌐 Afficher une interface web locale intégrée au binaire, sans dépendance externe

---

## ⚙️ Composants principaux

### 1. API HTTP (`localhost:9090`)

Deux endpoints JSON REST :

* **`/api/status`**
  Retourne l’état complet du système :

  * Utilisation CPU, modèle et nombre de cœurs
  * Mémoire vive et swap
  * Usage disque (`/`)
  * Statistiques réseau (I/O)
  * Charge moyenne du système
  * Informations hôte (nom, uptime…)
  * Nombre de connexions TCP actives

* **`/api/scan?ip=<adresse_ip>`**
  Scan rapide des ports TCP (20 à 1024) sur l’IP donnée, retourne les ports ouverts.

---

### 2. Interface Web (`localhost:8080`)

Serveur HTTP qui héberge un dashboard HTML/CSS/JS moderne avec :

* UI basée sur [Shoelace UI](https://shoelace.style/)
* Graphiques dynamiques via Chart.js
* Formulaire pour scanner une IP
* Cartes d’information système (CPU, mémoire, disque…)

Les fichiers statiques `/www` sont **servis directement depuis le dossier `/www` dans l’image**.

---

## 📦 Packaging & déploiement

* 🏗️ Application compilée en binaire autonome multiplateforme (Linux, Windows, macOS)
* 🐳 Conteneur Docker minimal `FROM scratch` sans dépendances externes
* 📁 Projet organisé avec une arborescence claire
* ☁️ Facile à déployer localement ou sur un serveur

---

## ✅ Fonctionnalités clés

| Fonction            | Description                                    |
| ------------------- | ---------------------------------------------- |
| 💻 API système      | Expose métriques hardware du serveur local     |
| 🔍 Scan de ports    | Teste les ports ouverts sur une IP donnée      |
| 🌍 UI embarquée     | Interface graphique intégrée dans le binaire   |
| 🚀 Binaire statique | Pas de dépendance ni serveur externe           |
| 🌐 Middleware CORS  | Autorise les appels cross-origin               |
| 🔄 Multi-serveur    | Serveur API (9090) + UI (8080) lancés ensemble |

---

## 📁 Structure du projet (simplifiée)

```plaintext
Makefile 
build/ 
  ├── Dockerfile 
  ├── app-linux 
  └── app-windows.exe 
src/ 
  ├── main.go          // Lance les deux serveurs (API + UI)
  ├── go.mod          
  ├── go.sum          
  ├── status.go        // /api/status 
  ├── scan.go          // /api/scan 
  ├── middleware.go    // CORS 
  ├── routes.go        // Routes API 
  ├── server.go        // Serveur web statique (fichiers dans /www)
  └── www/             // Fichiers front embarqués (servis en statique)
        ├── index.html 
        └── assets/ 
            ├── css/style.css 
            └── js/dashboard.js
```

---

## 🛠️ Technologies utilisées

| Technologie     | Usage                         |
| --------------- | ----------------------------- |
| **Go**          | Serveur HTTP, API, scan ports |
| **HTML/CSS/JS** | Interface utilisateur         |
| **Chart.js**    | Graphiques dynamiques         |
| **Shoelace UI** | Composants UI modernes        |
| **gopsutil**    | Récupération d’infos système  |

---

## 📚 Documentation API REST

L’API est disponible sur le port **9090**.

---

### 1. GET `/api/status`

Retourne un JSON avec les métriques système :

* Utilisation CPU (modèle, nombre de cœurs, % d’utilisation)
* Mémoire (RAM & swap)
* Disque (espace utilisé / total)
* Interfaces réseau (trafic entré/sorti)
* Charge moyenne du système
* Infos hôte (hostname, uptime)
* Nombre de connexions TCP actives

**Exemple de requête**

```bash
curl http://localhost:9090/api/status
```

**Réponse**

```json
{
  "cpu": {
    "model": "Intel(R) Core(TM) i7",
    "cores": 8,
    "usage_percent": 23.5
  },
  "memory": {
    "total": 16777216,
    "used": 7340032,
    "free": 9437184
  },
  "disk": {
    "total": 256000000000,
    "used": 100000000000,
    "free": 156000000000
  },
  "network": {
    "interfaces": [
      {
        "name": "eth0",
        "bytes_sent": 1234567,
        "bytes_recv": 7654321
      }
    ]
  },
  "load_avg": [0.12, 0.15, 0.10],
  "host": {
    "hostname": "monserveur",
    "uptime_seconds": 3600
  },
  "tcp_connections": 42
}
```

---

### 2. GET `/api/scan?ip=<adresse_ip>`

Scanne les ports TCP **20 à 1024** sur l’adresse IP donnée et retourne les ports ouverts.

**Exemple**

```bash
curl "http://localhost:9090/api/scan?ip=192.168.1.1"
```

**Réponse**

```json
{
  "open_ports": [22, 80, 443]
}
```

---

## 🌐 Interface Web (port 8080)

* Accessible via [http://localhost:8080](http://localhost:8080)
* UI moderne basée sur **Shoelace UI**
* Graphiques dynamiques avec **Chart.js**
* Formulaire pour scanner une IP et afficher les résultats
* Cartes d’information système (CPU, mémoire, disque, réseau, hôte)

---

## 🚀 Instructions de build et lancement

Le projet utilise un **Makefile** pour simplifier les commandes.

### 1. Compiler localement (Linux x64)

```bash
make build-monitoring
```

Le binaire sera dans `build/app`

---

### 2. Construire l’image Docker

```bash
make docker-build
```

Crée l’image `monitoring` selon `build/Dockerfile`

---

### 3. Lancer le container Docker

```bash
make run
```

Expose les ports `8080` (front) et `9090` (API).

---

### 4. Lancer sans Docker (Go requis)

```bash
make run-test
```

---

### 5. Nettoyer les builds

```bash
make clean
```

---

### 6. Builds multiplateformes

```bash
make build-linux
make build-windows
make build-macos
make build-all
```

---

### 7. Pipeline complet : clean + build + docker + run

```bash
make auto
```

---

### 8. Git commit + push avec message dynamique

```makefile
git-push:
	@if [ -z "$(m)" ]; then \
		echo "❌ Merci de fournir un message avec m=\"message\""; \
		exit 1; \
	fi
	@git add .
	@git commit -m "$(m)"
	@git push
```

Usage terminal :

```bash
make git-push m="✨ Mise à jour avec Docker"
```

---

# Fin de la documentation — bonne continuation ! 🚀

---
---
---


# 📈 Plan d’améliorations pour ton app Monitoring

---

## A. **Améliorations Frontend**

 **Ajouter des composants dynamiques**

   * Graphiques interactifs (avec Chart.js)
   * Tableau de bord avec sections repliables, filtres, recherche
 **Responsive design** pour support mobile/tablettes
 **Améliorer l’UX :**

   * Indicateurs visuels de chargement
   * Feedback erreurs / succès lors des scans
 **Possibilité d’intégrer websocket** pour mettre à jour en temps réel les données

---

## B. **Améliorations Backend / Go**

1. **Structurer clairement les services** (API, scans, status, etc.)
2. **Utiliser des goroutines pour :**

   * Scanner plusieurs IPs en parallèle
   * Rafraîchir les données système périodiquement sans bloquer l’API
3. **Ajouter de nouvelles métriques et infos système** :

   * Infos sur les processus
   * Logs système
   * Surveillance des ressources disque / température / batterie etc.
4. **Gestion des erreurs et logs** plus robustes
5. **Ajouter un système de configuration (fichier ou flags CLI)**
6. **Sécuriser l’API (authentification / permissions)**
7. **Tests unitaires et d’intégration**

---

## C. **DevOps & CI/CD**

1. **Automatiser les builds et tests** dans GitHub Actions ou GitLab CI
2. **Déployer sur un serveur ou cloud** (ex: Docker Hub + Kubernetes ou simple VPS)
3. **Surveillance de la santé de l’app** (logs, alertes)

---

# 🎯 Proposition de roadmap pour avancer étape par étape

| Étape | Objectif                                  | Détails                                                                          |
| ----- | ----------------------------------------- | -------------------------------------------------------------------------------- |
| 1     | Nettoyer et améliorer le front            | Refaire les pages CPU / Host / Network, ajouter styles |
| 2     | Optimiser le backend avec goroutines      | Scan IP parallèle, rafraîchissement asynchrone                                   |
| 3     | Ajouter des métriques supplémentaires     | Logs, processus, température, ressources                                         |
| 4     | Ajouter authentification & sécurité API   | Token JWT simple, clés API, gestion des accès                                    |
| 5     | Ajouter tests et CI/CD                    | Tests unitaires + pipeline GitHub Actions                                        |
| 6     | Automatiser le déploiement Docker         | Mise en place déploiement continu                                                |
| 7     | Ajouter WebSocket ou SSE                  | Push en temps réel du monitoring au front                                        |
| 8     | Responsive, mobile friendly et UX avancée | Interface plus ergonomique                                                       |

---
