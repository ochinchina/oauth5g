[![Go Report Card](https://goreportcard.com/badge/github.com/ochinchina/oauth5g)](https://goreportcard.com/report/github.com/ochinchina/oauth5g)

# oauth5g

This project implements the oauth2 server & proxy defined in 3GPP TS 29.510. The 3GPP TS 29.510 defined how the oauth2 authorization should be implemented in the 5G network. This project mainly implements the interface defined in this standard.

The architecture is shown as below:

<img src="https://github.com/ochinchina/oauth5g/blob/master/5g-oauth-architecture.png" width="400x200">

# Compile the project

Install golang development environment in your enviroment and then compile the source code with following command:

```shell
# go get -u github.com/ochinchina/oauth5g
```

# Run oauth2 server

The script run-server.sh will start the oauth2 server. The server configuration is defined in server.yaml file. The following command is used to start the oauth server:

```shell
# ./run-server.sh
```

# Run oauth2 proxy

The script run-proxy.sh can be used to start the oauth2 proxy and its configuration can be found in proxy.yaml file.

```shell
# ./run-proxy.sh
```

After start the proxy, the user can get the token from the client by curl command:

```shell
# curl http://127.0.0.1:8082/reqtoken -d@token_request.json -H "Content-Type: application/json"
eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiQU1GIl0sImV4cCI6MTYwMzU0ODcwNywiaXNzIjoiNjg4ZDc1MGMtMTQzYy0xMWViLWFlMmUtNmZlMjZhOGVkODc4Iiwic3ViIjoiTE1GLTEiLCJzY29wZSI6Im5hbWYtY29tbSJ9.fI3cNo4ahxRYQGwapIDo50g-erxMPXpRe8i1OXhJmDz1oX2X4cOJKLHHzD57OeSNze9k_ymAHfQRql7pEvlmSS3spPNIKOcJ09sylX9kl1xv3wu7WvIDFrS9RitFkIVfX_GznKwYvQXy-w9aH4jMHeUHZjz4MSrtM1z6NYQR1OlNagjtAapdEaMSBKqzNLHns3aTkH_pD2XPDCRpozkq7FtMA-QjLrl2rObRKWYVgO7sgpVI-q4M4dHdWeM2wuwoCBBKHNXm-eBP15qwl_fruFiycp6tYIZ3qVGOrBc2Dkj5B3HIy99KxS-kxnDUCsbm0UmUzxcN7973Dq9Phuwkxw
```
