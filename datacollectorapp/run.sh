if [ -f /hostetc/iscsi/initiatorname.iscsi ]; then
  echo "iscsi=`grep "InitiatorName=" /hostetc/iscsi/initiatorname.iscsi | awk -F  "=" '{print $2}'`"
fi
if [ -f /hostproc/modules ]; then
  echo "sdc=`grep "scini" /hostproc/modules | awk -F  " " '{print $2}'`"
fi
if [ -f /hostusr/local/sbin/mount.nfs ] || [ -f /hostusr/local/bin/mount.nfs ] || [ -f /hostusr/sbin/mount.nfs ] || [ -f /hostusr/bin/mount.nfs ]; then
  echo "nfs=true"
fi
