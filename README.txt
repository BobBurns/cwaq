CLOUDWATCH ALARM QUERY - cwaq
==============================
Nagios plugin to check AWS Cloudwatch Alarm and return status. Useful for when
 you just want to check if a predefined alarm is present on a AWS instance and
 utilize nagios alerting capability.

USAGE
=====
$ ./cwaq -i <instance-id> -m <metric> [-n <namespace>]


INSTALL
=======
Go get the latest version of GO here: https://golang.org/doc/install
Then in cwaq directory,
$ go get
$ go build

Make sure the server you are running this program on has the keys for the hosts
it is checking, or for an added security benefit, run on an AWS instance and
give the hosts a policy to allow AWS Cloudwatch access (eg. "Action" : "cloudwatch:DescribeAlarmsForMetric").

NOTES
=====
You might want to play with the time period parameter to get better granularity.

EXAMPLE OUTPUT
==============
bob@b79ce:~/godev/src/github.com/BobBurns/cwaq$ ./cwaq -i i-009614dxxxxxxx -m CPUUtilization
OK | state=0
bob@b79ce:~/godev/src/github.com/BobBurns/cwaq$ echo $?
0
bob@b79ce:~/godev/src/github.com/BobBurns/cwaq$ ./cwaq -i i-009614dxxxxxxx -m CPUUtilization
ALARM - Threshold Crossed: 1 out of the last 1 datapoints [99.636612021858 (08/09/18 19:10:00)] was greater than or equal to the threshold (10.0) (minimum 1 datapoint for OK -> ALARM transition). | state=1
bob@b79ce:~/godev/src/github.com/BobBurns/cwaq$ echo $?
2
bob@b79ce:~/godev/src/github.com/BobBurns/cwaq$
