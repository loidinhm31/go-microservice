## Modify module for import module in local environment
````
go mod edit -replace example.com/greetings=../greetings
````

## Connect Mongo Compass

````
mongodb://admin:password@localhost:27017/logs?authSource=admin&readPreference=primary&directConnection=true&ssl=false
````