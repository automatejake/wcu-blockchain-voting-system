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
    mkdir "processes$a"
    cp process1.go "processes$a/process$a.go"
    echo -e "PORT=700$a" >> "processes$a/.env"
    if $runProcesses ; then
        echo "starting process $a on port 700$a"
        go run "processes$a/process$a.go" &
    fi
    
    
done

