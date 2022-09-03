# pulumi-iaac-aws
Repo to setup  AWS platform and EKS using Pulumi 


## Steps to provision Stack
- Navigate to VPC sub-path 

  > pulumi up 

  > Specify the stack name and region

- Navigate to vpc-networking and update CIDR in code

  > pulumi up

- Navigate to EKS sub-path

  > pulumi up

- Connect to EKS cluster 

  > aws eks update-kubeconfig --name <provisioned-cluster-name> --region us-west-2