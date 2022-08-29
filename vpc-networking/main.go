package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a NAT and associate EIP
		ngwEIP, err := ec2.NewEip(ctx, "NGW-eip", &ec2.EipArgs{
			Vpc: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		// lookup subnet via CIDR
		pubSub, err := ec2.LookupSubnet(ctx, &ec2.LookupSubnetArgs{
			Filters: []ec2.GetSubnetFilter{
				ec2.GetSubnetFilter{
					Name: "tag:Name",
					Values: []string{
						"pub-sub-1",
					},
				},
			},
		}, nil)
		if err != nil {
			return err
		}

		ngw, err := ec2.NewNatGateway(ctx, "dev-natgw", &ec2.NatGatewayArgs{
			AllocationId: ngwEIP.AllocationId,
			SubnetId:     pulumi.String(pubSub.Id),
			Tags: pulumi.StringMap{
				"Name": pulumi.String("dev-natgw"),
			},
		},
		)
		if err != nil {
			return err
		}

		vpcCIDR := "10.0.0.0/16"
		// Get VPC with CIDR block
		lookupVPC, err := ec2.LookupVpc(ctx, &ec2.LookupVpcArgs{
			CidrBlock: &vpcCIDR,
		})

		// lookup IGW
		lookupIGW, err := ec2.LookupInternetGateway(ctx, &ec2.LookupInternetGatewayArgs{
			Tags: map[string]string{"Name": "dev-igw"},
		})

		// create route tables
		pubRT, err := ec2.NewRouteTable(ctx, "pub-routes", &ec2.RouteTableArgs{
			VpcId: pulumi.String(lookupVPC.Id),
			Routes: ec2.RouteTableRouteArray{
				&ec2.RouteTableRouteArgs{
					CidrBlock: pulumi.String("0.0.0.0/0"),
					GatewayId: pulumi.String(lookupIGW.InternetGatewayId),
				},
			},
			Tags: pulumi.StringMap{
				"Name": pulumi.String("pub-routetable"),
			},
		})
		if err != nil {
			return err
		}

		privRT, err := ec2.NewRouteTable(ctx, "priv-routes", &ec2.RouteTableArgs{
			VpcId: pulumi.String(lookupVPC.Id),
			Routes: ec2.RouteTableRouteArray{
				&ec2.RouteTableRouteArgs{
					CidrBlock: pulumi.String("0.0.0.0/0"),
					GatewayId: ngw.ID(),
				},
			},
			Tags: pulumi.StringMap{
				"Name": pulumi.String("priv-routetable"),
			},
		})
		if err != nil {
			return err
		}

		// Export the name of the bucket
		ctx.Export("NAT gateway", ngw.ID())
		ctx.Export("Public RT", pubRT.ID())
		ctx.Export("Private RT", privRT.ID())
		return nil
	})
}
