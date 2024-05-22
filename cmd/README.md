# kubectl接口



目前支持的接口如下

kubectl apply -f \<filename\> (可带-a或--all标志)

kubectl get pod 

kubectl get deployment

kubectl get service

kubectl get node

kebectl describe pod \<pod_name\>

kubectl describe deployment \<deployment_name\>

kubectl describe service \<service_name\>

kubectl describe node \<node_name\>

kebectl delete pod \<pod_name\>

kubectl delete deployment \<deployment_name\>

kubectl delete service \<service_name\>



如要使用上述命令，需要先构建所有文件

```
cd cmd
go run kubectl.go
```

