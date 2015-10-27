package ecu

// Instance encapsulates representative information regarding AWS EC2
// instances.
type Instance struct {
	APIName         string
	Memory          float64
	ComputeUnitsx10 int64 // x10 because a lot of them are fractional
	Cores           int
	ECUPerCore      float64
	Burstable       bool
}

// Data from http://www.ec2instances.info
var instances = []*Instance{
	{"c1.medium", 1.7, 50, 2, 2.5, false},
	{"c1.xlarge", 7.0, 200, 8, 2.5, false},
	{"c3.2xlarge", 15.0, 280, 8, 3.5, false},
	{"c3.4xlarge", 30.0, 550, 16, 3.438, false},
	{"c3.8xlarge", 60.0, 1080, 32, 3.375, false},
	{"c3.large", 3.75, 70, 2, 3.5, false},
	{"c3.xlarge", 7.5, 140, 4, 3.5, false},
	{"c4.2xlarge", 15.0, 310, 8, 3.875, false},
	{"c4.4xlarge", 30.0, 620, 16, 3.875, false},
	{"c4.8xlarge", 60.0, 1320, 36, 3.667, false},
	{"c4.large", 3.75, 80, 2, 4, false},
	{"c4.xlarge", 7.5, 160, 4, 4, false},
	{"cc2.8xlarge", 60.5, 880, 32, 2.75, false},
	{"cg1.4xlarge", 22.5, 335, 16, 2.094, false},
	{"cr1.8xlarge", 244.0, 880, 32, 2.75, false},
	{"d2.2xlarge", 61.0, 280, 8, 3.5, false},
	{"d2.4xlarge", 122.0, 560, 16, 3.5, false},
	{"d2.8xlarge", 244.0, 1160, 36, 3.222, false},
	{"d2.xlarge", 30.5, 140, 4, 3.5, false},
	{"g2.2xlarge", 15.0, 260, 8, 3.25, false},
	{"g2.8xlarge", 60.0, 1040, 32, 3.25, false},
	{"hi1.4xlarge", 60.5, 350, 16, 2.188, false},
	{"hs1.8xlarge", 117.0, 350, 17, 2.059, false},
	{"i2.2xlarge", 61.0, 270, 8, 3.375, false},
	{"i2.4xlarge", 122.0, 530, 16, 3.312, false},
	{"i2.8xlarge", 244.0, 1040, 32, 3.25, false},
	{"i2.xlarge", 30.5, 140, 4, 3.5, false},
	{"m1.large", 7.5, 40, 2, 2, false},
	{"m1.medium", 3.75, 20, 1, 2, false},
	{"m1.small", 1.7, 10, 1, 1, false},
	{"m1.xlarge", 15.0, 80, 4, 2, false},
	{"m2.2xlarge", 34.2, 130, 4, 3.25, false},
	{"m2.4xlarge", 68.4, 260, 8, 3.25, false},
	{"m2.xlarge", 17.1, 65, 2, 3.25, false},
	{"m3.2xlarge", 30.0, 260, 8, 3.25, false},
	{"m3.large", 7.5, 65, 2, 3.25, false},
	{"m3.medium", 3.75, 30, 1, 3, false},
	{"m3.xlarge", 15.0, 130, 4, 3.25, false},
	{"m4.10xlarge", 160.0, 1245, 40, 3.112, false},
	{"m4.2xlarge", 32.0, 260, 8, 3.25, false},
	{"m4.4xlarge", 64.0, 535, 16, 3.344, false},
	{"m4.large", 8.0, 65, 2, 3.25, false},
	{"m4.xlarge", 16.0, 130, 4, 3.25, false},
	{"r3.2xlarge", 61.0, 260, 8, 3.25, false},
	{"r3.4xlarge", 122.0, 520, 16, 3.25, false},
	{"r3.8xlarge", 244.0, 1040, 32, 3.25, false},
	{"r3.large", 15.25, 65, 2, 3.25, false},
	{"r3.xlarge", 30.5, 130, 4, 3.25, false},
	{"t1.micro", 0.613, 1, 1, 1, true},
	{"t2.large", 8.0, 2, 2, 1, true},
	{"t2.medium", 4.0, 2, 2, 1, true},
	{"t2.micro", 1.0, 1, 1, 1, true},
	{"t2.small", 2.0, 1, 1, 1, true},
}

var instanceLookup map[string]*Instance

func init() {
	instanceLookup = make(map[string]*Instance)
	for _, inst := range instances {
		instanceLookup[inst.APIName] = inst
	}
}

// LookupName finds an instance of the type given.  If one exists, it
// returns that instance and true; otherwise it returns nil and false.
func LookupName(name string) (*Instance, bool) {
	i, ok := instanceLookup[name]
	return i, ok
}
