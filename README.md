druns
=====

Web application to manage clients, developed for Evandro Pelissoni Levy.

install
-------

1. Install the following packages:

```
% sudo apt-get install postgresql rsyslog
```

2. Enable rsyslog to listen on TCP or UDP ports

```
% sudo vim /etc/rsyslog.conf
% sudo service rsyslog restart
```

3. Generate and install the project debian package:

```
% sudo apt-get install ruby-dev gcc
% sudo gem install fpm
% cd <project>/deploy
% ./gendeb.sh <version> <release>
% sudo dpkg -i druns_<version>-<release>_amd64.deb
```

4. Create the database in postgres:

```
% sudo -u postgres createdb druns
% sudo -u postgres psql druns < /usr/druns/db/structure.sql
```

5. Initialize the system:

```
% /usr/druns/bin/bootstrap -config /usr/druns/etc/webserver.conf -email <email> -name "<name>" -password "<password>"
```

6. Add periodics to crontab

```
% crontab -e

*/10 * * * * /usr/druns/bin/scheduler /usr/druns/etc/webserver.conf
```