kubectl delete statefulset nodo
kubectl delete service raft
kubectl delete pod cliente
kubectl create -f statefulset_go.yaml

