buffer=$(printf '01234567890%.0s' {1..100})

{
    while :; do
        for (( i = 0; i < 10; i++ )); do
            echo "put x: 1 y: 2 height: 10 width: 10 fg: #f00 bg: #220 text: \"${buffer:$i:100}\""
        done
    done
} >&3
