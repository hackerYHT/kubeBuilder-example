# kubeBuilder-example
自定义CRD-controllerOperator

```shell
sudo go mod init apiserver
sudo chmod 777 go.mod   //给go.mod权限
export GOPROXY=https://goproxy.io     	会报错，走代理
kubebuilder init --domain ppdapi.com --owner yht
kubebuilder create api --group apps --version v1 --kind Podsbook
go env -w GOPROXY=https://goproxy.cn
make manifests
```

[参考链接1](https://www.cnblogs.com/alisystemsoftware/p/11580202.html)

[参考链接2](https://podsbook.com/posts/kubernetes/operator/#创建一个operator项目)

[官方快速入门](https://book.kubebuilder.io/getting-started)
