#!/bin/bash

if ! [ -v "${GOPATH}" ];
then
    GOPATH=$(go env GOPATH)
    export GOPATH
fi

#
FILES=(pkg/core/models/converter/converter.exe.go 
pkg/core/models/converter/converter.type.go 
pkg/core/models/converter/converter.go 
pkg/core/models/storage/storage.go 
pkg/core/models/storage/storage.key.go 
pkg/core/models/storage/storage.obj.go 
pkg/core/models/logger/logger.go 
internal/service/thumb/thumb.go 
internal/pkg/img/convert/convert.go 
internal/pkg/img/convert/checks/source.go )

#
MOCKGEN=$GOPATH/bin/mockgen
MOCKS_DESTINATION=internal/test/mocks

#
echo "Generating mocks..."
rm -rf $MOCKS_DESTINATION
for file in ${FILES[@]}; do
        echo $file | sed "s/internal/intr/g" ;
        $MOCKGEN -source=$file -destination=$MOCKS_DESTINATION/`echo $file | sed "s/internal/intr/g"`;
    done;
