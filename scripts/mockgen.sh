#!/bin/bash

# если возникла проблема с запускам срипта после редактирования под виндой:
#   /bin/bash^M: bad interpreter: No such file or directory
# запустить:
#   sed -i -e 's/\r$//' scripts/mockgen.sh
# команда заменит все символы переноса Windows на Linux

if ! [ -v "${GOPATH}" ];
then
    GOPATH=$(go env GOPATH)
    export GOPATH
fi

#
FILES=(interfaces/converter/converter.exe.go 
interfaces/converter/converter.type.go 
interfaces/converter/converter.go 
interfaces/storage/storage.go 
interfaces/storage/storage.key.go 
interfaces/storage/storage.obj.go 
interfaces/logger/logger.go 
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
