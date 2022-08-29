package main

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		slug := fmt.Sprintf("%v", ctx.Stack())
		subs := map[string]string{"pub-sub-1": "10.0.1.0/24", "pub-sub-2": "10.0.3.0/24", "priv-sub-1": "10.0.10.0/24", "priv-sub-2": "10.0.12.0/24"}
		azs := []string{"us-west-2b", "us-west-2c", "us-west-2b", "us-west-2c"}
		vpcName := strings.Join([]string{"vpc-eks", slug}, "-")
		vpc, err := ec2.NewVpc(ctx, vpcName, &ec2.VpcArgs{
			CidrBlock:          pulumi.String("10.0.0.0/16"),
			EnableDnsHostnames: pulumi.BoolPtr(true),
		})
		if err != nil {
			return err
		}
		subCount := 0
		for sub, cidr := range subs {
			_, err = ec2.NewSubnet(ctx, sub, &ec2.SubnetArgs{
				VpcId:                                vpc.ID(),
				CidrBlock:                            pulumi.String(cidr),
				AvailabilityZone:                     pulumi.String(azs[subCount]),
				EnableResourceNameDnsARecordOnLaunch: pulumi.BoolPtr(true),
				MapPublicIpOnLaunch:                  pulumi.BoolPtr(true),
				Tags: pulumi.StringMap{
					"Name": pulumi.String(sub),
				},
			},
				pulumi.DependsOn([]pulumi.Resource{vpc}),
			)
			subCount++
			if err != nil {
				return err
			}
		}

		intGW, err := ec2.NewInternetGateway(ctx, "dev-igw", &ec2.InternetGatewayArgs{
			VpcId: vpc.ID(),
			Tags: pulumi.StringMap{
				"Name": pulumi.String("dev-igw"),
			},
		})
		if err != nil {
			return err
		}

		// Export the name and ID of VPC
		ctx.Export("VPC ID ", vpc.ID())
		ctx.Export("Igw is ", intGW.ID())
		return nil
	})

}
