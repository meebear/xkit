package main

import (
    "xkit/mytest/clip"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/aws/awserr"

    "fmt"
)

func awsCmdInit(c *clip.Command) {
    awsCmd := c.SubCommand("aws-ec2", "commands for AWS EC2 operations", nil)
    awsCmd.SubCommand("describe-instances", "describe instances", describeInstances)
    awsCmd.SubCommand("assign-secondary-ips", "assign secondary private addresses", assignSecPrivIPs)
}

func newSession() (*session.Session, error) {
    // Load session from shared config
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    return sess, nil
}

func describeInstances(c *clip.Command) error {
    sess, _ := newSession()

    // Create new EC2 client
    ec2Svc := ec2.New(sess)

    // Call to get detailed information on each instance
    result, err := ec2Svc.DescribeInstances(nil)
    if err != nil {
        fmt.Println("Error", err)
    } else {
        fmt.Println("Success", result)
    }
    return nil
}

// ha-t0: 172.31.32.35                   eni-0c83fa17955f79ab5
// ha-t1: 172.31.37.154  172.31.41.157   eni-076154949556e9b89

func assignSecPrivIPs(c *clip.Command) error {
    sess, _ := newSession()

    svc := ec2.New(sess)
    input := &ec2.AssignPrivateIpAddressesInput{
        NetworkInterfaceId: aws.String("eni-0c83fa17955f79ab5"),
        PrivateIpAddresses: []*string{ aws.String("172.31.41.157"), },
        AllowReassignment: aws.Bool(true),
    }

    result, err := svc.AssignPrivateIpAddresses(input)
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            default:
                fmt.Println(aerr.Error())
            }
        } else {
            // Print the error, cast err to awserr.Error to get the Code and
            // Message from an error.
            fmt.Println(err.Error())
        }
        return nil
    }

    fmt.Println(result)
    return nil
}
