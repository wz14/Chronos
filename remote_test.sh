docker-compose up aac_build || exit
./build/start remoteGen || exit
for ((i=0;i<=3;i++))
do
  mkdir $i
  cp config_$i.yaml ./$i/config.yaml || exit
  cp ./build/start ./$i/start || exit
  cd $i || exit
  ./start remote 2>/dev/null &
  cd ..
done

