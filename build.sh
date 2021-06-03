#!/bin/bash

Name=sync_eth
Version=$1

#################
if [ -z "$Version" ]; then
    echo "Version is empty"
    exit 25
fi
##########################

Image=${Name}:${Version}

##################
function docker_check(){
    printf "progress:[%-50s]%d%% ${Name} container is run?\r" $1 $2
    docker ps -a | grep ${Name} &> /dev/null
    if [ $? -eq 0 ]
    then
        printf "progress:[%-50s]%d%% ${Name} need stop container\r" $1 $2
    else
        printf "progress:[%-50s]%d%% ${Name} container stopped\r" $1 $2
    fi
}
##########################

##################
function image_check(){
    printf "progress:[%-50s]%d%% ${Image} image is exist?\r" $1 $2
    docker images | grep ${Name} |grep ${Version} &> /dev/null
    if [ $? -eq 0 ]
    then
        printf "progress:[%-50s]%d%% ${Image} del image...\r" $1 $2
        docker rmi $(docker images | grep ${Name} |grep ${Version} | awk '{print $3}') &> /dev/null
        printf "progress:[%-50s]%d%% ${Image} image deleted,start build\r" $1 $2
    else
        printf "progress:[%-50s]%d%% ${Image} images not exist,start build\r" $1 $2
    fi
}
##########################

##################
function image_build(){
    docker build . -t ${Image} &> /dev/null
    if [ $? -eq 0 ]
    then
        printf "progress:[%-50s]%d%% ${Image} build success\r" $1 $2
    else
        printf "progress:[%-50s]%d%% ${Image} build failed\r" $0 $1
        exit 26
    fi
}
##########################

##################
function image_clean(){
    printf "progress:[%-50s]%d%% clean err image===\r" $1 $2
    docker rmi `docker images | grep  "<none>" | awk '{print $3}'` &> /dev/null
    printf "progress:[%-50s]%d%% clean image success===\r" $1 $2
}
##########################


##################
function main(){
    array=("docker_check" "image_check" "image_build" "image_clean")
    length=${#array[@]}
    let one=100/${length}
    b=''
    for ((i=0; i<$length; i++))
    do
        let jingdu=($i+1)*$one
        ${array[$i]} $b $jingdu
        sleep 0.1
        a=`printf %.s# {1..25}`
        b=$a$b
    done
    echo
}

Note=`printf %.s# {1..75}`

case "$2" in
  docker_check)
        docker_check ${Note} 100
        echo
        ;;
  image_check)
        image_check ${Note} 100
        echo
        ;;
  image_build)
        image_build ${Note} 100
        echo
        ;;
  image_clean)
        image_clean ${Note} 100
        echo
        ;;
  help)
        echo $"Usage: $0 $1 {docker_check|image_check|image_build|image_clean}"
        exit 1
          ;;
  *)
        main
        ;;
esac
exit $RETVAL
##########################
