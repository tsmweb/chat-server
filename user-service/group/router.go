package group

import "fmt"

const version string = "v1"

var (
	resource string
	resourceMember string
)

func init() {
	resource = fmt.Sprintf("/%s/group", version)
	resourceMember = fmt.Sprintf("/%s/group/member", version)
}
