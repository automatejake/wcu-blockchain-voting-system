read -p 'Enter the number of Processes: ' num

# touch "process1.go"

for ((a = 0 ; a < $num; a++ ));
do
    mkdir "processes$a"
    cp process1.go "processes$a/process$a.go"
    echo -e "PORT=700$a" >> "processes$a/.env"
    go run "processes$a/process$a.go" &
done

