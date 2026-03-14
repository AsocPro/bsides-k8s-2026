#!/bin/bash
# Enable systemd cgroup delegation for Podman + Dagger + k3s
#
# Problem: Podman doesn't delegate cpuset (and other) cgroup controllers
# to containers by default. k3s requires cpuset, so it fails with:
#   "failed to find cpuset cgroup (v2)"
#
# Solution: Tell systemd to delegate all cgroup controllers to user sessions.
# This is also recommended by Podman docs for rootless containers.
#
# After running this script, restart the Dagger engine:
#   podman rm -f dagger-engine-v0.20.1
#   dagger call --help  # triggers engine restart
#
# Usage: sudo ./setup-cgroup-delegation.sh

set -euo pipefail

if [ "$(id -u)" -ne 0 ]; then
  echo "Error: must run as root (sudo $0)" >&2
  exit 1
fi

CONF_DIR="/etc/systemd/system/user@.service.d"
CONF_FILE="${CONF_DIR}/delegate.conf"

echo "Creating ${CONF_FILE} ..."
mkdir -p "${CONF_DIR}"
cat > "${CONF_FILE}" <<EOF
[Service]
Delegate=yes
EOF

echo "Reloading systemd daemon ..."
systemctl daemon-reload

echo ""
echo "Done. Cgroup delegation is now enabled for user sessions."
echo ""
echo "Next steps:"
echo "  1. Restart the Dagger engine:"
echo "       podman rm -f \$(podman ps -q --filter name=dagger-engine)"
echo "       dagger version  # triggers engine restart"
echo ""
echo "  2. Verify cpuset is delegated:"
echo "       podman exec \$(podman ps -q --filter name=dagger-engine) cat /sys/fs/cgroup/cgroup.controllers"
echo "       # Should include: cpuset cpu io memory hugetlb pids"
