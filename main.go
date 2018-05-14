package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

func main() {
  // get some metadata
  region, InstanceID := RegionInstanceID()
	sessionInfo, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		log.Fatalln("Could not create new AWS session, please check AWS Credentials in use")
	}
  action := strings.ToLower(os.Args[2])
	ELBsvc := elb.New(sessionInfo)
	if len(os.Args) > 1 && os.Args[1] != "" {
		//do all the things in here
		elbs := strings.Split(os.Args[1], ",")
		for _, elb := range elbs {
      // check if action is "register" or "deregister". Fail if neither
      if action == "register" {
        log.Printf("Registering instance with ELB Name: %s", elb)
  			//make call to ELB API
  			_, err := RegisterInstanceWithELB(InstanceID, elb, ELBsvc)
  			if err != nil {
  				log.Fatalf("Encountered error registering instance with ELB: %s\n Error: %s", elb, err)
  			}
  			log.Printf("Instance %s registered with ELB %s successfully", InstanceID, elb)
      } else if action == "deregister" {
        log.Printf("Deregistering instance with ELB Name: %s", elb)
  			//make call to ELB API
  			_, err := DeregisterInstanceFromELB(InstanceID, elb, ELBsvc)
  			if err != nil {
  				log.Fatalf("Encountered error deregistering instance with ELB: %s\n Error: %s", elb, err)
  			}
  			log.Printf("Instance %s deregistered with ELB %s successfully", InstanceID, elb)
      } else {
        log.Fatalf("Usage: ./ELB_connector ELB,LIST,DELIMITED,WITH,COMMA register|deregister")
      }

		}
	} else {
		log.Fatalf("Usage: ./ELB_connector ELB,LIST,DELIMITED,WITH,COMMA register|deregister")
	}
}

// RegionInstanceID Returns the metadata of the ec2 instance and returns region and InstanceID
func RegionInstanceID() (region string, InstanceID string) {
	meta := ec2metadata.New(session.New())
	ec2InstanceIDentifyDocument, _ := meta.GetInstanceIdentityDocument()
	region = ec2InstanceIDentifyDocument.Region
	InstanceID = ec2InstanceIDentifyDocument.InstanceID
	//fmt.Println(InstanceID)
	return region, InstanceID
}

// RegisterInstanceWithELB Register this instance with its ELB
func RegisterInstanceWithELB(InstanceID string, ELBName string, elbClient *elb.ELB) (bool, error) {
	registerArgs := &elb.RegisterInstancesWithLoadBalancerInput{
		Instances: []*elb.Instance{
			{
				InstanceId: aws.String(InstanceID),
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

// DeregisterInstanceFromELB Removes this instance with its ELB
func DeregisterInstanceFromELB(InstanceID string, ELBName string, elbClient *elb.ELB) (bool, error) {
	deregisterArgs := &elb.DeregisterInstancesFromLoadBalancerInput{
    Instances: []*elb.Instance{
        {
            InstanceId: aws.String(InstanceID),
        },
    },
    LoadBalancerName: aws.String(ELBName),
}

	_, err := elbClient.DeregisterInstancesFromLoadBalancer(deregisterArgs)
	if err != nil {
		return false, errors.New(err.Error())
	}
	return true, nil
}
