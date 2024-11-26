package algorithm

import (
	"strings"
)

var keywords = [...]string{
    "graphql",
    "graphiql",
    "cosmo router",
    "dataloader",
    "schema stitching",
    "wundergraph",
    "n+1 problem",
    " gql ",
    "apollo client",
    "apollo engine",
    "apollo federation",
    "apollo gateway",
    "apollo server",
    "apollo studio",
    "apollo summit",
    "relay client",
    "grafbase",
    " hasura ",
    "appsync",
    "rover cli",
}

func PostFilterMatch(body *string) bool {
	for _, keyword := range keywords {
		if strings.Contains(*body, keyword) {
			return true
		}
	}
	return false
}
