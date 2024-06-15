# task: download image & unzip into project/tmp dir

CWD=$(dirname $(realpath ${0}))
TMP_DIR=${CWD}/../tmp
FILENAME_IMAGE=${TMP_DIR}/image.img
SD_CARD=${SD_CARD:-""}

if [ ! -f ${FILENAME_IMAGE} ]; then
  echo "abort, ${FILENAME_IMAGE} is missing"
  exit 1
fi

if [ ! -b ${SD_CARD} ]; then
  echo "abort, env var: ${SD_CARD} does not point to an block device"
  exit 1
fi

if [ -x $(which umount) ]; then
  sudo umount ${SD_CARD}
fi
if [ -x $(which diskutil) ]; then
  diskutil unmountDisk ${SD_CARD}
fi

sudo tee ${SD_CARD} < ${FILENAME_IMAGE} > /dev/null
MOUNT_SD_PATH=$(df ${SD_CARD}* | awk 'END{print $NF}')
