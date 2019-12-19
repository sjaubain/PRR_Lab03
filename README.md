# Programmation Répartie - Laboratoire III

Auteurs: Jobin Simon, Teklehaimanot Robel

## Énoncé du problème

Réalisez un programme en langage GO qui implémente l'algorithme d'élection de Chang et Roberts avec pannes possibles des processus que nous avons vu en classe. Pour ce faire, nous avons des processus qui, de temps en temps, interrogent le dernier site élu, et si celui-ci n'est plus opérationnel, une élection est démarrée. Un processus en panne doit pouvoir se réinsérer au sein de l'anneau lors de sa reprise.

## Fonctionnement

Le programme est décomposé en trois parties : la partie *network*, s'occupant de l'envoi et de la réception des messages en UDP, la partie *algoCR* étant le noeud principal s'occupant de gérer les élections ainsi que la partie *main* qui lance le programme, initialise le réseau et la partie algorithmique et offre la possibilité à l'utilisateur d'entrer les commandes `E` pour lancer une élection ou `G` pour obtenir la valeur courante de l'élu.

Une fois une élection lancée, l'annonce circule dans l'anneau et un mécanisme d'acquittement permet à un site de transmettre au suivant si son voisin directe ne répond pas dans un délai choisi (timeout).

## Configuration

Un fichier Json contenant le nombre de sites ainsi que leurs adresses et leurs aptitudes permet de configurer comme on le souhaite l'architecture des differents processus. Par défaut, toutes les addresses IP sont *localhost*, puisqu'on fait tourner tous les sites localement. Il est donc possible de changer ce fichier de configuration et d'ajouter d'autres sites (local ou remote) selon leur répartition.

Pour ce laboratoire, nous avons décidé d'attribuer les identifiants des sites avec des valeurs comprises entre *0* et *n-1*, où *n* est le nombre de sites.

Nous avons par conséquent choisi de construire tous les payloads envoyés de la manière suivante :

* **byte 1** : opcode `R` (résultat), `A` (annonce) ou `F` (fin de l'élection)

Pour les messages annonces :

* **du byte 2 jusqu'à la fin** : tuples *id-apt* des sites, séparés par un `;`

Pour les messages résultats

* **du byte 2 jusqu'à la fin** : liste des sites ayant reçu le résultat, séparés par `,`

Notons que le `F` est utilisé lorsque l'on désire informer les autres sites que l'élu a été déterminé.

## Utilisation

Tout d'abord, cloner le repository https://github.com/sjaubain/PRR_Lab03 quelque part depuis le GOPATH. Pour lancer les différents sites, il faut ouvrir un terminal où se situe le fichier `main.go` puis lancer la commande suivante pour construire l'exécutable :

```bash
go build main.go
```

Ouvrir ensuite autant de terminal qu'il y a de sites configurés (par défaut 3) et les lancer avec leur identifiant en argument. Le numéro du site doit être entre *0* et *n-1*, où *n* est le nombre de sites. Par exemple :

```bash
./site 0
```

Une fois que tous sites seront lancés et interconnectés, il suffit comme cité précédemment d'entrer les commandes `E` ou `G` pour lancer une élection ou consulter la valeur de l'élu

## Amélioration possible

On pourrait imaginer comme amélioration un ping régulier sur le site élu pour voir si celui-ci est en panne, auquel cas on relancerait une nouvelle élection.
