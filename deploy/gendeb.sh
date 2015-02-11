#!/bin/sh

usage() {
  echo "Usage: $1 <version> <release>"
}

pack_name="druns"
version="$1"
release="$2"
vendor="Rafael Dantas Justo"
maintainer="Rafael Dantas Justo <adm@rafael.net.br>"
url="http://rafael.net.br"
license="All rights reserved"
description="Web application that schedule clients in a weekday agenda"

if [ -z "$version" ]; then
  echo "Version not defined!"
  usage $0
  exit 1
fi

if [ -z "$release" ]; then
  echo "Release not defined!"
  usage $0
  exit 1
fi

install_path=/usr/druns
tmp_dir=/tmp/druns
project_root=$tmp_dir$install_path

workspace=`echo $GOPATH | cut -d: -f1`
workspace=$workspace/src/github.com/rafaeljusto/druns

# recompiling everything
current_dir=`pwd`
cd $workspace
go build
cd $workspace/utils/bootstrap
go build
cd $workspace/utils/password
go build
cd $current_dir

if [ -f $pack_name*.deb ]; then
  # remove old deb
  rm $pack_name*.deb
fi

if [ -d $tmp_dir ]; then
  rm -rf $tmp_dir
fi

mkdir -p $tmp_dir$install_path/bin $tmp_dir$install_path/web $tmp_dir$install_path/db
mv $workspace/druns $workspace/utils/bootstrap/bootstrap $workspace/utils/password/password $project_root/bin/
cp -r $workspace/web/templates $project_root/web/
cp -r $workspace/web/assets $project_root/web/
cp -r $workspace/etc $project_root/
cp $workspace/core/db/structure.sql $tmp_dir$install_path/db/

fpm -s dir -t deb \
  --exclude=.git -n $pack_name -v "$version" --iteration "$release" --vendor "$vendor" \
  --maintainer "$maintainer" --url $url --license "$license" --description "$description" \
  --deb-upstart $workspace/deploy/druns.upstart \
  --deb-user root --deb-group root \
  --prefix / -C $tmp_dir usr/druns
