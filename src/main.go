package main

import (
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/kataras/golog"
	"gopkg.in/yaml.v2"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	roleMapping := flag.String("role-mappings", "", "--role-mappings=devs:system:masters;contractors:read-only,contractor-partial-admin-role")
	flag.Parse()

	//get all mapping tuples
	mt := strings.Split(*roleMapping, ";")

	im := make(InputMap)
	for _, v := range mt {
		rm := strings.SplitN(v, ":", 2)
		if len(rm) != 2 {
			//incorrect mapping
			golog.Error("Cannot parse string into role mapping: ", v)
		} else {
			golog.Info("Loading parsed config into map: ", v)
			//enumerate the k8s roles
			roleArr := strings.Split(rm[1], ",")
			im[rm[0]] = roleArr
		}
	}
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		for iamGroup, k8sRoles := range im {
			users := getAwsIamGroup(iamGroup)
			cf, err := clientset.CoreV1().ConfigMaps("kube-system").Get("aws-auth", metav1.GetOptions{})
			if err != nil {
				panic(err.Error())
			}
			var newConfig []MapUserConfig

			for _, user := range users.Users {
				newConfig = append(newConfig, MapUserConfig{
					UserArn:  *user.Arn,
					Username: *user.UserName,
					Groups:   k8sRoles,
				})
			}
			roleStr, err := yaml.Marshal(newConfig)
			if err != nil {
				golog.Error(err)
			}
			cf.Data["mapUsers"] = string(roleStr)

			newCF, err := clientset.CoreV1().ConfigMaps("kube-system").Update(cf)
			if err != nil {
				golog.Error(err)
			} else {
				golog.Info("successfully updated user roles")
				golog.Info(newCF)
			}
			time.Sleep(10 * time.Second)
		}
	}
}

func getAwsIamGroup(groupName string) *iam.GetGroupOutput {
	sess := session.Must(session.NewSession(&aws.Config{}))
	iamClient := iam.New(sess)
	group, err := iamClient.GetGroup(&iam.GetGroupInput{
		GroupName: aws.String(groupName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				golog.Error(iam.ErrCodeNoSuchEntityException, aerr.Error())
			case iam.ErrCodeServiceFailureException:
				golog.Error(iam.ErrCodeServiceFailureException, aerr.Error())
			default:
				golog.Error(aerr.Error())
			}
		}
	}
	return group
}
