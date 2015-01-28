#!/usr/bin/env bash
#

# Copyright 2015 Crunchy Data Solutions, Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# option defaults
yum_config=/etc/yum.conf

#--------------------

#target=$(mktemp -d --tmpdir $(basename $0).XXXXXX)
target=/tmp/myimage
mkdir $target

set -x


mkdir -m 755 "$target"/dev
mknod -m 600 "$target"/dev/console c 5 1
mknod -m 600 "$target"/dev/initctl p
mknod -m 666 "$target"/dev/full c 1 7
mknod -m 666 "$target"/dev/null c 1 3
mknod -m 666 "$target"/dev/ptmx c 5 2
mknod -m 666 "$target"/dev/random c 1 8
mknod -m 666 "$target"/dev/tty c 5 0
mknod -m 666 "$target"/dev/tty0 c 4 0
mknod -m 666 "$target"/dev/urandom c 1 9
mknod -m 666 "$target"/dev/zero c 1 5


yum -c "$yum_config" --releasever=6Server --installroot="$target" \
--setopt=tsflags=nodocs --setopt=group_package_types=mandatory \
-y groupinstall Core

cp /etc/yum.repos.d/* $target/etc/yum.repos.d/

yum -c "$yum_config" --installroot="$target" -y --nogpgcheck install uuid-devel-1.6.1-10.el6.x86_64.rpm uuid-1.6.1-10.el6.x86_64.rpm

yum -c "$yum_config" --installroot="$target" -y install  docker-io postgresql93 postgresql93-server postgresql93-contrib which unzip openssh-server openssh-clients nginx libxslt pgpool-II-93


yum -c "$yum_config" --installroot="$target" -y clean all

cat > "$target"/etc/sysconfig/network <<EOF
NETWORKING=yes
HOSTNAME=localhost.localdomain
EOF

# effectively: febootstrap-minimize --keep-zoneinfo --keep-rpmdb
# --keep-services "$target".  Stolen from mkimage-rinse.sh
#  locales
rm -rf "$target"/usr/{{lib,share}/locale,{lib,lib64}/gconv,bin/localedef,sbin/build-locale-archive}
#  docs
rm -rf "$target"/usr/share/{man,doc,info,gnome/help}
#  cracklib
rm -rf "$target"/usr/share/cracklib
#  i18n
rm -rf "$target"/usr/share/i18n
#  sln
rm -rf "$target"/sbin/sln
#  ldconfig
rm -rf "$target"/etc/ld.so.cache
rm -rf "$target"/var/cache/ldconfig/*

version=
if [ -r "$target"/etc/redhat-release ]; then
    version="$(sed 's/^[^0-9\]*\([0-9.]\+\).*$/\1/' "$target"/etc/redhat-release)"
fi

if [ -z "$version" ]; then
    echo >&2 "warning: cannot autodetect OS version, using '$name' as tag"
    version=$name
fi

tar --numeric-owner -c -C "$target" . | docker import - crunchy:rhel65
docker run -i -t crunchy:rhel65  echo success

