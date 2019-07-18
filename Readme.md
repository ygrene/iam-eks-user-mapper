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
6) Edit the `kubernetes/deployment.yaml` `command:` with both the IAM groups (by name) you want to provide access to, and the Kubernetes groups each user in the group should be mapped to.

The model for AWS IAM Group to Kubernetes Role Mapping is as follows:

| Symbol | Description
| --- | --- | 
| `;` | Separates a mapping tuple 
| `:` | Separates an AWS IAM Group from a list of Kubernetes Roles 
| `,` | Separates a list of Kubernetes Roles for which an IAM Group maps to

An Example argument: 

`devs:cluster-admin;contractors:partial-admin,no-prod-access-role`

7) Finally:
```bash
$ kubectl apply -f kubernetes/
```
8) Rejoice, now user management will be a bit easier.

## Have suggestions or want to contribute?
Raise a PR or file an issue, I'd love to help!

## Changelog:

### 0.0.1
* Initial Release

### 0.0.2
**BREAKING CHANGE**
* Modified command-line arguments to accept a new flag `--role-mappings`
* Removes the flags `--k8s-cap` and `--aws-iam-group`
* Functionality change to now allow mapping of multiple AWS IAM Groups to multiple Kubernetes Groups