set -xe

# Hustoj basic file system
useradd -m -u 1536 judge
mkdir -p /home/judge/etc
mkdir -p /home/judge/run0
mkdir -p /home/judge/data/1000
echo "1 2" >  /home/judge/data/1000/1.in
echo "3" >  /home/judge/data/1000/1.out
mkdir -p /home/judge/log
mkdir -p /home/judge/backup
mkdir -p /var/log/hustoj
mv /trunk /home/judge/src
chmod -R 700 /home/judge/etc
chmod -R 700 /home/judge/backup
chmod -R 700 /home/judge/src/web/include/
chown judge /home/judge/run0
chmod 750 /home/judge/run0
chown -R www-data:www-data /home/judge/data
chown -R www-data:www-data /home/judge/src/web 

# Judge daemon and client
make      -C /home/judge/src/core/judged
make      -C /home/judge/src/core/judge_client
make exes -C /home/judge/src/core/sim/sim_3_01
cp /home/judge/src/core/judged/judged                /usr/bin/judged
cp /home/judge/src/core/judge_client/judge_client    /usr/bin/judge_client 
cp /home/judge/src/core/sim/sim_3_01/sim_c.exe       /usr/bin/sim_c
cp /home/judge/src/core/sim/sim_3_01/sim_c++.exe     /usr/bin/sim_cc
cp /home/judge/src/core/sim/sim_3_01/sim_java.exe    /usr/bin/sim_java
cp /home/judge/src/core/sim/sim.sh                   /usr/bin/sim.sh
cp /home/judge/src/install/hustoj                    /etc/init.d/hustoj
chmod +x /home/judge/src/install/ans2out
chmod +x /usr/bin/judged
chmod +X /usr/bin/judge_client
chmod +x /usr/bin/sim_c
chmod +X /usr/bin/sim_cc
chmod +x /usr/bin/sim_java
chmod +x /usr/bin/sim.sh

# Adjust system configuration
cp /home/judge/src/install/java0.policy  /home/judge/etc/
cp /home/judge/src/install/judge.conf    /home/judge/etc/

echo "OJ_USE_PTRACE=0" >> /home/judge/etc/judge.conf
sed -i "s#OJ_HOST_NAME[[:space:]]*=[[:space:]]*127.0.0.1#OJ_HOST_NAME=mysql#g"    /home/judge/etc/judge.conf
sed -i "s#OJ_USER_NAME[[:space:]]*=[[:space:]]*root#OJ_USER_NAME=root#g"    /home/judge/etc/judge.conf
sed -i "s#OJ_PASSWORD[[:space:]]*=[[:space:]]*root#OJ_PASSWORD=testtest#g"      /home/judge/etc/judge.conf
sed -i "s#OJ_COMPILE_CHROOT[[:space:]]*=[[:space:]]*1#OJ_COMPILE_CHROOT=0#g"     /home/judge/etc/judge.conf
sed -i "s#OJ_RUNNING[[:space:]]*=[[:space:]]*1#OJ_RUNNING=1#g"                /home/judge/etc/judge.conf
sed -i "s#OJ_SHM_RUN[[:space:]]*=[[:space:]]*1#OJ_SHM_RUN=0#g"                   /home/judge/etc/judge.conf
sed -i "s#DB_HOST[[:space:]]*=[[:space:]]*\"localhost\"#DB_HOST=\"mysql\"#g"                  /home/judge/src/web/include/db_info.inc.php
sed -i "s#DB_USER[[:space:]]*=[[:space:]]*\"root\"#DB_USER=\"root\"#g"                  /home/judge/src/web/include/db_info.inc.php
sed -i "s#DB_PASS[[:space:]]*=[[:space:]]*\"root\"#DB_PASS=\"testtest\"#g"                  /home/judge/src/web/include/db_info.inc.php

PHP_VER=`apt-cache search php-fpm|grep -e '[[:digit:]]\.[[:digit:]]' -o`
if [ "$PHP_VER" = "" ] ; then PHP_VER="8.1"; fi

	echo "modify the default site"
	sed -i "s#root /var/www/html;#root /home/judge/src/web;#g" /etc/nginx/sites-enabled/default
	sed -i "s:index index.html:index index.php:g" /etc/nginx/sites-enabled/default
	sed -i "s:#location ~ \\\.php\\$:location ~ \\\.php\\$:g" /etc/nginx/sites-enabled/default
	sed -i "s:#\tinclude snippets:\tinclude snippets:g" /etc/nginx/sites-enabled/default
	sed -i "s|#\tfastcgi_pass unix|\tfastcgi_pass unix|g" /etc/nginx/sites-enabled/default
	sed -i "s:}#added by hustoj::g" /etc/nginx/sites-enabled/default
	sed -i "s:php7.4:php$PHP_VER:g" /etc/nginx/sites-enabled/default
	sed -i "s|# deny access to .htaccess files|}#added by hustoj\n\n\n\t# deny access to .htaccess files|g" /etc/nginx/sites-enabled/default

sed -i "s/post_max_size = 8M/post_max_size = 500M/g" /etc/php/8.1/fpm/php.ini     
sed -i "s/upload_max_filesize = 2M/upload_max_filesize = 500M/g" /etc/php/8.1/fpm/php.ini 

if grep "client_max_body_size" /etc/nginx/nginx.conf ; then 
    echo "client_max_body_size already added" ;
else
    sed -i "s:include /etc/nginx/mime.types;:client_max_body_size    280m;\n\tinclude /etc/nginx/mime.types;:g" /etc/nginx/nginx.conf
fi

chmod 755 /home/judge
