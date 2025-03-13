set -xe

if [ ! -d /volume ]; then
mkdir /volume;
fi

if [ ! -d /volume/backup ]; then
cp -rp /home/judge/backup  /volume/backup; 
fi 
if [ ! -d /volume/data ]; then  
cp -rp /home/judge/data /volume/data;  
fi 
if [ ! -d /volume/etc ]; then   
cp -rp /home/judge/etc /volume/etc;
fi 
if [ ! -d /volume/web ]; then   
cp -rp /home/judge/src/web /volume/web;
fi 

chmod -R 700 /volume/etc
chmod -R 700 /volume/backup 
chmod -R 700 /volume/web/include/   
chown -R www-data:www-data /volume/data 
chown -R www-data:www-data /volume/web  
chown -R www-data:www-data /var/log/hustoj
rm -rf /home/judge/backup   
rm -rf /home/judge/data 
rm -rf /home/judge/etc  
rm -rf /home/judge/src/web  
ln -s /volume/backup /home/judge/backup 
ln -s /volume/data   /home/judge/data   
ln -s /volume/etc    /home/judge/etc
ln -s /volume/web    /home/judge/src/web

RUNNING=`cat /home/judge/etc/judge.conf | grep OJ_RUNNING`
RUNNING=${RUNNING:11}
for i in `seq 1 $RUNNING`; do
    mkdir -p    /home/judge/run`expr ${i} - 1`;
    chown judge /home/judge/run`expr ${i} - 1`;
done 

ln -sf /dev/stdout /var/log/nginx/access.log
ln -sf /dev/stderr /var/log/nginx/error.log

service php8.1-fpm start  
service hustoj     start  
nginx -g "daemon off;"