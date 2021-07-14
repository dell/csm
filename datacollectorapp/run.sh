if [ -f /hostetc/iscsi/initiatorname.iscsi ]; then
  echo "iscsi=`grep "InitiatorName=" /hostetc/iscsi/initiatorname.iscsi | awk -F  "=" '{print $2}'`"
fi
if [ -f /hostproc/modules ]; then
  echo "sdc=`grep "scini" /hostproc/modules | awk -F  " " '{print $2}'`"
fi
