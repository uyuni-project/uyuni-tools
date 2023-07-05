#!/bin/bash
#TODO this should be copied in /usr/local/bin/uyuni-tools

set -e
UYUNI_SHORTFORM=${UYUNI_FQDN%%.*}
#TODO parse this folder from option.json
for folder in /var/lib/pgsql \
              /var/cache \
              /var/spacewalk \
              /var/log \
              /srv/salt \
              /srv/www/htdocs/pub \
              /srv/www/cobbler \
              /srv/www/os-images \
              /srv/tftpboot \
              /srv/formula_metadata \
              /srv/pillar \
              /srv/susemanager \
              /srv/spacewalk \
              /root \
              /etc/apache2 \
              /etc/rhn \
              /etc/systemd/system/multi-user.target.wants \
              /etc/salt \
              /etc/tomcat \
              /etc/cobbler \
              /etc/sysconfig;
do
  rsync -e "ssh -A -o StrictHostKeyChecking=no" --rsync-path='sudo rsync' -avz $REMOTE_USER@$UYUNI_FQDN:$folder/ $folder;
done;
rm -f /srv/www/htdocs/pub/RHN-ORG-TRUSTED-SSL-CERT;
ln -s /etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT /srv/www/htdocs/pub/RHN-ORG-TRUSTED-SSL-CERT;

##SET NEW ADDRESS FOR REMOTE DEBUGGING. || true is required to avoid issue if file is missing
#sed "s/address=[^:]*:/address=$UYUNI_FQDN:/" -i /etc/rhn/taskomatic.conf || true;
#sed "s/address=[^:]*:/address=$UYUNI_FQDN:/" -i /etc/tomcat/conf.d/remote_debug.conf || true;
#sed "s/address=[^:]*:/address=$UYUNI_FQDN:/" -i /usr/lib/systemd/system/taskomatic.service.d/override.conf || true;


#SETUP POSTGRES
rhn-ssl-tool --gen-ca --no-rpm --set-common-name=$UYUNI_FQDN --set-country=$CERT_COUNTRY --set-state=$CERT_STATE --set-city=$CERT_CITY --set-org=$CERT_O --set-org-unit=$CERT_OU --set-email=$MANAGER_ADMIN_EMAIL --password=$CERT_PASS --force

cp /root/ssl-build/$UYUNI_SHORTFORM/server.crt /etc/pki/tls/certs/spacewalk.crt
cp /root/ssl-build/$UYUNI_SHORTFORM/server.key /etc/pki/tls/private/pg-spacewalk.key
chown postgres /etc/pki/tls/private/pg-spacewalk.key

#SETUP APACHE
cp /root/ssl-build/$UYUNI_SHORTFORM/server.key /etc/pki/tls/private/spacewalk.key

touch /root/.MANAGER_SETUP_COMPLETE
echo "DONE"

