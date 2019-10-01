read -p 'Enter the number of Processes: ' num
runProcesses=false
while true; do
    read -p "Do you wish to run processes?" yn
    case $yn in
        [Yy]* ) runProcesses=true; break;;
        [Nn]* ) break;;
        * ) echo "Please answer yes or no.";;
    esac
done

for ((a = 0 ; a < $num; a++ ));
do
    mkdir "process$a"
    cp process.go "process$a/process$a.go"
    echo -e "PORT=700$a" >> "process$a/.env"
    if $runProcesses ; then
        echo "starting process $a on port 700$a"
        # go run "process$a/process$a.go" &
    fi
    

done

