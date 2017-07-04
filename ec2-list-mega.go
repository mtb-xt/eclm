package main

// Loosely based on http://www.heystephenwood.com/2015/02/list-running-ec2-instances-with-golang.html
// No support for multi-account configurations (yet?)

import (
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    //"github.com/olekukonko/tablewriter"
    "github.com/mtb-xt/tablewriter"
    "net/url"
    "os"
    "runtime"
    //"strings"
    //"sort"
    terminal "github.com/wayneashleyberry/terminal-dimensions"
    "sync"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

//func printIds(region string, wg *sync.WaitGroup) {
func printIds(sess *session.Session, wg *sync.WaitGroup) {
    defer wg.Done()

    svc := ec2.New(sess)

    // Here we create an input that will filter any instances that aren't either
    // of these two states. This is generally what we want
    params := &ec2.DescribeInstancesInput{
        Filters: []*ec2.Filter{
            &ec2.Filter{
                Name: aws.String("instance-state-name"),
                Values: []*string{
                    aws.String("running"),
                    aws.String("pending"),
                },
            },
        },
    }

    resp, err := svc.DescribeInstances(params)
    if err != nil {
        fmt.Println("there was an error listing instances in", err.Error())
        // i threw away the log package >_>
        //log.Fatal(err.Error())
    }

    data := make([][]string, 0)

    // Loop through the instances. They don't always have a name-tag so set it
    // to None if we can't find anything.
    for idx, _ := range resp.Reservations {
        for _, inst := range resp.Reservations[idx].Instances {
            // We need to see if the Name is one of the tags. It's not always
            // present and not required in Ec2.
            name := "None"
            for _, keys := range inst.Tags {
                if *keys.Key == "Name" {
                    name = url.QueryEscape(*keys.Value)
                }
            }

            important_vals := []*string{
                inst.InstanceId,
                &name,
                inst.PrivateIpAddress,
                inst.InstanceType,
                inst.PublicIpAddress,
            }

            // Convert any nil value to a printable string in case it doesn't
            // doesn't exist, which is the case with certain values
            output_vals := []string{}
            for _, val := range important_vals {
                if val != nil {
                    output_vals = append(output_vals, *val)
                } else {
                    output_vals = append(output_vals, "None")
                }
            }
            // The values that we care about, in the order we want to print them
            data = append(data, output_vals)
        }
    }
    // Don't output anything if the region is empty
    if len(data) > 0 {
        table := tablewriter.NewWriter(os.Stdout)
        uinttermwidth, _ := terminal.Width()
        var termwidth int = int(uinttermwidth)

        // Maximum field lengths of all fields except Name is 72, so get the maximum Name size with the current term width:
        maxname := termwidth - 73
        table.SetHeader([]string{"Instance Id", "Name", "Private IP", "Type", "Public IP"})
        table.SetBorder(false)
        table.SetAutoWrapText(false)
        table.AppendBulk(data)
        // Set Name column to calculated size
        table.SetColMinWidth(1,maxname)
        // Set IP and type columns to their maximum possible length
        table.SetColMinWidth(2,15)
        table.SetColMinWidth(3,10)

        table.Render()

    }

}

func main() {
    // Go for it!
    runtime.GOMAXPROCS(runtime.NumCPU())
    // Create a session to share configuration, and load external configuration.
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    // Make sure the config file exists
    config := os.Getenv("HOME") + "/.aws/config"
    if _, err := os.Stat(config); os.IsNotExist(err) {
        fmt.Println("No config file found at: %s", config)
        os.Exit(1)
    }

    var wg sync.WaitGroup

    svc := ec2.New(sess)

    // Iterate over every single stinking region to get a list of available
    // ec2 instances
    regions, err := svc.DescribeRegions(&ec2.DescribeRegionsInput{})
    check(err)
    //fmt.Println(regions)
    for _, region := range regions.Regions {
        wg.Add(1)
        //fmt.Printf("Calling region %s\n", *region.RegionName)
        //go printIds(*region.RegionName, &wg)
        go printIds(sess.Copy(aws.NewConfig().WithRegion(*region.RegionName)), &wg)
        //sess.Copy(&aws.Config(Region: region.RegionName))
    }
    //  }

    // Allow the goroutines to finish printing
    wg.Wait()

}
