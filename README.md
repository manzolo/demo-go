# Go Demo
This is a simple service that print ip, hostname and datetime of server 

[docker image](https://hub.docker.com/r/manzolo/demo-go) 

Useful to test a [multipass microk8s cluster](https://github.com/manzolo/multipass-microk8s-cluster-demo.git)

## Example
```
curl http://localhost:8080
{
    "id": "6281838661429879825",
    "hostname": "go-deployment-c7f68fbd5-xjmrt",
    "ip": "10.1.36.76",
    "datetime": "2023.01.19 00:23:26"
}                                              
```