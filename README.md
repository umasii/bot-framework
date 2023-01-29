# CicadaAIO

Go task framework used for Cicada AIO.

Finally decided to open source this, as it's a little shame for this code to go to waste. Apologies for some of the mess, the code was essentially in-between versions and hastily patched to be deployable for the LV drop. 

Uses goroutines and custom networking solutions (TLS and HTTP/2 client not open sourced, there are plenty already out there) to automate purchasing on a variety of websites. Highly extensible to new sites/features, and capable of running thousands of concurrent tasks while still remaining highly performant.

Used to have a UI/GRPC integration as well, but this was removed prior to open sourcing. Not meant to be run as some imports are in private repositories, although this can be easily fixed. I've removed all features/packages not relevant to the task framework.

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

Credits to https://github.com/g4cko, https://github.com/trrai, as well as https://github.com/Epacity for the LV module code.

Also big shoutout to https://github.com/amcode21 for teaching me what I know.
