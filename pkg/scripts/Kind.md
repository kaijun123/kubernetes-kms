## Introduction to Kind

- When using Kind to run Kubernetes, there is an additional level of virtualization added. 
- In Kind, Kubernetes clusters are run as Docker containers on the host OS (ie your laptop)
- Within those Docker containers run the pods (Eg etcd, apiserver, or user-deployed pods)
