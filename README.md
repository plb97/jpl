# jpl
Exemple (expérimental) d'utilisation du langage Go en astronomie à partir des éphémérides DE432 du JPL.

Source d'informations ftp://ssd.jpl.nasa.gov/pub/eph/planets

## Lectures intéressantes

NASA et Caltech Jet Propulsion Laboratory (JPL) 

*   ftp://ssd.jpl.nasa.gov/pub/eph/planets/README.txt
*   https://ssd.jpl.nasa.gov/?planet_eph_export
*   ftp://ssd.jpl.nasa.gov/pub/eph/planets/fortran

Description des groupes (1010, 1030, 1040, 1041, 1050 et 1070)

*   https://eqbridges.wordpress.com/2010/02/15/understanding-jpl-ephemerides-data-pt-2/

## Utilisation des tests

Avant tout, il faut creer la base MySQL.

Pour cela, un conteneur *Docker* est très facile à créer en choisissant pour l'utilisateur *root* un mot de passe à la place de *\<password\>* :

    docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=<password> -e MYSQL_DATABASE=test -e MYSQL_USER=test -e MYSQL_PASSWORD=test --name mysql_test mysql/mysql-server:5.7
   
Puis il faut initialiser la base de données à l'aide de la fonction *main* du package *main*.

Cela prend du temps, beaucoup de temps (plusieurs heures), donc patience...

**Remarque**: Le choix d'utiliser une base de données dans un conteneur n'est probablement pas le meilleur pour les performances.
La table des coefficients est très grosse (401792 lignes) et donc partitionnée (par numéro de planète) pour être utilisable.
