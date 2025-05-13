run: build
    #!/usr/bin/env bash
    if [ ! -d logs ]; then mkdir logs; fi;
    log_name="logs/$(date --rfc-3339=seconds | sed 's/ /T/')"
    touch $log_name
    rm -rf logs/current
    ln -s "$(pwd)/$log_name" logs/current
    echo Starting
    ./goserv &> $log_name

build:
    go build -o goserv src/*

pibuild:
    GOOS=linux GOARCH=arm64 go build -o piserv src/*

clean:
    rm -rf goserv piserv logs/
