package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// Nagios exit codes
const (
	UNKNOWN  = 3
	CRITICAL = 2
	OK       = 0
)

func usage() {
	fmt.Println("Usage: ./cwaq -i <instance-id> -m <metric> [-n <namespace]")
}

func main() {

	// get cli args
	idFlag := flag.String("i", "", "Instance ID of the host to check")
	metricFlag := flag.String("m", "", "Metric to check for alarm")
	namespaceFlag := flag.String("n", "AWS/EC2", "Namespace")

	flag.Parse()

	if *idFlag == "" || *metricFlag == "" {
		usage()
		os.Exit(-1)
	}

	// init aws session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	checkerrF(err)

	svc := cloudwatch.New(sess)

	// build parameters for DescribeAlarmsForMetric call

	// []*Dimension
	var dims []*cloudwatch.Dimension
	dims = append(dims, &cloudwatch.Dimension{
		Name:  aws.String("InstanceId"),
		Value: aws.String(*idFlag),
	})

	params := cloudwatch.DescribeAlarmsForMetricInput{
		Dimensions: dims,
		MetricName: aws.String(*metricFlag),
		Namespace:  aws.String(*namespaceFlag),
		Period:     aws.Int64(300),
		Statistic:  aws.String("Average"),
	}

	// make call
	resp, err := svc.DescribeAlarmsForMetric(&params)
	checkerrP(err)

	// Check result and output nagios friendly exit codes
	// for more detail about MetricAlarm type see
	// https://docs.aws.amazon.com/sdk-for-go/api/service/cloudwatch/#MetricAlarm
	n := len(resp.MetricAlarms) - 1
	if n >= 0 {
		a := *resp.MetricAlarms[n].StateValue
		if a == "ALARM" {
			fmt.Println(a, "-", *resp.MetricAlarms[n].StateReason, "| state=1")
			os.Exit(CRITICAL)
		} else if a == "OK" {
			fmt.Println(a, "| state=0")
			os.Exit(OK)
		} else {
			fmt.Println(a, "-", *resp.MetricAlarms[n].StateReason, "| state=0")
			os.Exit(UNKNOWN)
		}
	}

}

// helpers
func checkerrF(e error) {
	if e != nil {
		panic(e)
	}
}

func checkerrP(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, e.Error())
		os.Exit(-1)
	}
}
