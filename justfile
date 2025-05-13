start:
    #!/usr/bin/env bash
    if [ ! -d logs ]; then mkdir logs; fi;
    go run src/* &> "logs/$(date --rfc-3339=seconds | sed 's/ /T/')"
