# task: download image & unzip into project/tmp dir

CWD=$(dirname $(realpath ${0}))
TMP_DIR=${CWD}/../tmp
IMAGE_URL=${IMAGE_URL:-"https://downloads.raspberrypi.com/raspios_oldstable_armhf/images/raspios_oldstable_armhf-2024-03-12/2024-03-12-raspios-bullseye-armhf.img.xz"}
FILENAME_COMPRESSED=${TMP_DIR}/image.img.xz
FILENAME_IMAGE=${TMP_DIR}/image.img

mkdir -p ${TMP_DIR}

if [ ! -f ${FILENAME_COMPRESSED} ]; then
  curl -L --output ${FILENAME_COMPRESSED} ${IMAGE_URL}
fi

if [ ! -f ${FILENAME_IMAGE} ]; then
  gunzip -c ${FILENAME_COMPRESSED} > ${FILENAME_IMAGE}
fi
