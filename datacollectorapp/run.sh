if [ ! -f /hostetc/iscsi/initiatorname.iscsi ]; then
  echo "InitiatorName=Unavailable"
else
  grep "InitiatorName=" /hostetc/iscsi/initiatorname.iscsi | awk -F  ":" '{print $0}'
fi