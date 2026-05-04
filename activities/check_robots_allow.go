package activities

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"temporal-playground/specs"

	"github.com/gobwas/glob"
	"go.temporal.io/sdk/temporal"
)

func (a *Activities) CheckRobotsAllow(ctx context.Context, robots []specs.RobotRules, requestUrl string) error {
	urlData, err := url.Parse(requestUrl)
	if err != nil {
		return err
	}

	if !isAllowedByRobots(robots, urlData) {
		return temporal.NewNonRetryableApplicationError(
			"Disallowed by robots.txt",
			"Disallowed",
			errors.New("Disallowed by robots.txt"),
		)
	}

	return nil
}

func isAllowedByRobots(robots []specs.RobotRules, url *url.URL) bool {
	var relevantRobotsTxt *specs.RobotRules
	for _, currRobots := range robots {
		if currRobots.UserAgent == "*" {
			relevantRobotsTxt = &currRobots
		}
	}

	largestDisallow := 0
	largestAllow := 0
	allowed := true
	if relevantRobotsTxt != nil {
		for _, disallow := range relevantRobotsTxt.Disallow {
			if checkAllowPattern(disallow, url.Path) {
				largestDisallow = max(len(disallow), largestDisallow)
			}
		}

		for _, allow := range relevantRobotsTxt.Allow {
			if checkAllowPattern(allow, url.Path) {
				largestAllow = max(len(allow), largestDisallow)
			}
		}

		if largestAllow < largestDisallow {
			allowed = false
		}
	}

	return allowed
}

func checkAllowPattern(allowRule string, path string) bool {
	if strings.HasPrefix(path, allowRule) {
		return true
	}

	if strings.HasSuffix(allowRule, "/") {
		allowRule = allowRule[0 : len(allowRule)-1]
	}

	g := glob.MustCompile(allowRule)
	return g.Match(path)
}
