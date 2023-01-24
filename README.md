# CicadaAIO

Finally decided to open source this, as it's a little shame for this code to go to waste. As this was basically in between versions/stripped of a lot of functionality prior to being open sourced, as well as quickly patched for the Louis Vuitton drop, some of the code is messy/nonfunctional. 

Either way, it's not intended to be ran as my net/http and utls forks are not (yet) open sourced.

All modules removed except Louis Vuitton and Walmart.

# Commands

```
go run main.go -r group_id
```

Group_ID should be set for a task group in `Data/tasks.json`.

On Linux:

```
go build . && ./framework -r "all"
```

Credits to https://github.com/g4cko, https://github.com/trrai, as well as https://github.com/Epacity for the module code.
