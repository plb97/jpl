# JPL

Exemple (expérimental) d'utilisation du langage Go en astronomie à partir des éphémérides DE432 du JPL.

[Source d'informations](ftp://ssd.jpl.nasa.gov/pub/eph/planets)
[Source de données](ftp://ssd.jpl.nasa.gov/pub/eph/planets/ascii/de432)

## Lectures intéressantes

NASA et Caltech Jet Propulsion Laboratory (JPL)

* [<https://arxiv.org/pdf/1507.04291>](https://arxiv.org/pdf/1507.04291)
* [<http://ipnpr.jpl.nasa.gov/progress_report/42-196/196C.pdf>](http://ipnpr.jpl.nasa.gov/progress_report/42-196/196C.pdf)
* [<https://ssd.jpl.nasa.gov/?planet_eph_export>](https://ssd.jpl.nasa.gov/?planet_eph_export)
* [<https://naif.jpl.nasa.gov/pub/naif/toolkit_docs/Tutorials/pdf/individual_docs/>](https://naif.jpl.nasa.gov/pub/naif/toolkit_docs/Tutorials/pdf/individual_docs/)
* [<https://naif.jpl.nasa.gov/pub/naif/toolkit_docs/C/req/spk.html>](https://naif.jpl.nasa.gov/pub/naif/toolkit_docs/C/req/spk.html)
* [<ftp://ssd.jpl.nasa.gov/pub/eph/planets/README.txt>](ftp://ssd.jpl.nasa.gov/pub/eph/planets/README.txt)
* [<ftp://ssd.jpl.nasa.gov/pub/eph/planets/fortran>](ftp://ssd.jpl.nasa.gov/pub/eph/planets/fortran)

Description des groupes (1010, 1030, 1040, 1041, 1050 et 1070)

* [<https://eqbridges.wordpress.com/2010/02/15/understanding-jpl-ephemerides-data-pt-2/>](https://eqbridges.wordpress.com/2010/02/15/understanding-jpl-ephemerides-data-pt-2/)

## Utilisation des tests

Avant tout, il faut creer la base MySQL.

Pour cela, un conteneur *Docker* est très facile à créer en choisissant pour l'utilisateur *root* un mot de passe à la place de *\<password\>* :

    docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=<password> -e MYSQL_DATABASE=test -e MYSQL_USER=test -e MYSQL_PASSWORD=test --name mysql_test mysql/mysql-server:5.7

Puis il faut initialiser la base de données à l'aide de la fonction *main* du package *main*.

Cela prend du temps, beaucoup de temps (plusieurs heures sur un portable), donc patience...

**Remarque**: Le choix d'utiliser une base de données dans un conteneur n'est probablement pas le meilleur pour les performances.
La table des coefficients est très grosse (401792 lignes) et donc partitionnée (par numéro de planète) pour être utilisable.

## Remerciements

A la *NASA* et au *JPL* ansi qu'au projet *github.com/go-sql-driver/mysql*
