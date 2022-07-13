#!/bin/bash

docker build -t lamassu-virtual-device . 

docker run -it \
    -e VDEV_EST_SERVER_URL=$VDEV_EST_SERVER_URL \
    -e VDEV_AWS_IOT_CORE_ENDPOINT=$VDEV_AWS_IOT_CORE_ENDPOINT \
    -v $VDEV_DEVICE_CERTIFICATES_DIR:/app/device-certificates/ \
    -v $VDEV_EST_SERVER_CERT:/app/downstream/tls.crt \
    lamassu-virtual-device