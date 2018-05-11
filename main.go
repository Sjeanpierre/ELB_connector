package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

func main() {
	region := os.Getenv("EC2REGION")
	sessionInfo, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		log.Fatalln("Could not create new AWS session, please check AWS Credentials in use")
	}
	ELBsvc := elb.New(sessionInfo)
	InstanceID := os.Getenv("EC2_INSTANCE_ID")
	if InstanceID == "" {
		log.Fatal("EC2_INSTANCE_ID ENV var is not set")
	}
	if len(os.Args) > 1 && os.Args[1] != "" {
		//do all the things in here
		elbs := strings.Split(os.Args[1], ",")
		for _, elb := range elbs {
			log.Printf("Registering instance with ELB Name: %s", elb)
			//make call to ELB API
			_, err := RegisterInstanceWithELB(InstanceID, elb, ELBsvc)
			if err != nil {
				log.Fatalf("Encountered error registering instance with ELB: %s\n Error: %s", elb, err)
			}
			log.Printf("Instance %s registered with ELB %s successfully", InstanceID, elb)
		}
	} else {
		log.Fatalf("Usage: ./ELB_connector ELB,LIST,DELIMITED,WITH,COMMA")
	}
}

func RegisterInstanceWithELB(instanceID string, ELBName string, elbClient *elb.ELB) (bool, error) {
	registerArgs := &elb.RegisterInstancesWithLoadBalancerInput{
		Instances: []*elb.Instance{
			{
				InstanceId: aws.String(instanceID),
			},
		},
		LoadBalancerName: aws.String(ELBName),
	}

	_, err := elbClient.RegisterInstancesWithLoadBalancer(registerArgs)
	if err != nil {
		return false, errors.New(err.Error())
	}
	return true, nil
}
