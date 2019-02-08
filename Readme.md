# IAM EKS User Mapper

The general overview for what this tool does can be found here: https://ygrene.tech/mapping-iam-groups-to-eks-user-access-66fd745a6b77

## Setting up in your environment:
1) Have an AWS IAM Group with users that you want to have access to your EKS cluster
(https://console.aws.amazon.com/iam/home?#/groups)
2) Create a new IAM User with an IAM ReadOnly policy
3) Replace the ACCESS_KEY_ID environment variable in `kubernetes/deployment.yaml` with your new generated user's access key id
4) Replace the `awsKey:` variable in `deployment/secret.yaml` with the base64 contents of your generated user's secret access key
```bash
$ echo -n "secretkey" | base64
```
5) Update the `AWS_REGION` environment variable in `kubernetes/deployment.yaml` if you aren't running in `us-west-2` with your EKS cluster
6) Edit the `kubernetes/deployment.yaml` `command:` with both the IAM group name you want to provide access to, and the Kubernetes group each user in the group should be mapped to.
(there is an example in the manifest already)
7) Finally:
```bash
$ kubectl apply -f kubernetes/
```
8) Rejoice, now user management will be a bit easier.

## Have suggestions or want to contribute?
Raise a PR or file an issue, I'd love to help!