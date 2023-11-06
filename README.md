# Chuck - Container Image Update Checker

Monitor and notify about new container image versions.

```
$ chuck run
2023/11/06 23:51:05 INFO fetched containers count=2
2023/11/06 23:51:05 INFO parsed image registry=docker.io name=library/nginx tag=1.25
2023/11/06 23:51:05 INFO parsed image registry=docker.io name=library/nginx tag=1.23
Container /funny_jennings ({docker.io library/nginx 1.25}) can be updated: [1.25.1 1.25.2 1.25.3]
Container /nice_mendeleev ({docker.io library/nginx 1.23}) can be updated: [1.23.1 1.23.2 1.23.3 1.23.4 1.24 1.24.0 1.25 1.25.0 1.25.1 1.25.2 1.25.3]
```

## Installing

### From source
To install Chuck from source, run the following command:
```bash
go install .
```

or, if you prefer building the binary first:
```bash
go build .
```

## Usage
### List local containers
```bash
chuck run
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
