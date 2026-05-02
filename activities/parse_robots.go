package activities

import (
	"bufio"
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"temporal-playground/specs"
)

func (a *Activities) ParseRobots(ctx context.Context, robotsContent string) ([]specs.RobotRules, error) {
	reader := strings.NewReader(robotsContent)
	scanner := bufio.NewScanner(reader)
	ruleRegex := regexp.MustCompile(`^\s*(?P<name>[a-zA-Z0-9\-]+)\s*:\s*(?P<value>.+)\s*$`)
	signalRegex := regexp.MustCompile(`^\s*(?P<name>[a-zA-Z0-9\-]+)\s*=\s*(?P<value>.+)\s*$`)

	rules := []specs.RobotRules{}
	var currRule specs.RobotRules = newRule()
	for scanner.Scan() {
		line := scanner.Text()
		match := ruleRegex.FindStringSubmatch(line)
		if match != nil {
			ruleName := strings.ToLower(match[1])
			if ruleName != "user-agent" && currRule.UserAgent == "" {
				currRule.UserAgent = "*"
			}

			switch ruleName {
			case "user-agent":
				if currRule.UserAgent != "" {
					rules = append(rules, currRule)
					currRule = newRule()
				}

				currRule.UserAgent = match[2]

			case "disallow":
				currRule.Disallow = append(currRule.Disallow, match[2])

			case "allow":
				currRule.Allow = append(currRule.Allow, match[2])

			case "sitemap":
				currRule.Sitemaps = append(currRule.Sitemaps, match[2])

			case "license":
				currRule.License = match[2]

			case "content-signal":
				signals := strings.Split(match[2], ",")
				for _, signal := range signals {
					signalMatch := signalRegex.FindStringSubmatch(signal)
					if signalMatch != nil {
						switch signalMatch[1] {
						case "ai-input":
							currRule.ContentSignal.AiInput = match[2] == "yes"
						case "ai-train":
							currRule.ContentSignal.AiTrain = match[2] == "yes"
						case "search":
							currRule.ContentSignal.Search = match[2] == "yes"
						}

						currRule.ContentSignal.AiTrain = match[1] == "train"
					}
				}

				currRule.Sitemaps = append(currRule.Sitemaps, match[2])

			case "crawl-delay":
				delay, err := strconv.Atoi(match[2])
				if err != nil {
					return nil, errors.New("Crawl-Delay must be an integer. Got: " + match[2])
				}
				currRule.CrawlDelay = uint8(delay)
			}
		}
	}

	if currRule.UserAgent != "" {
		rules = append(rules, currRule)
	}

	return rules, nil
}

func newRule() specs.RobotRules {
	return specs.RobotRules{
		ContentSignal: specs.RobotContentSignal{
			AiTrain: true,
			Search:  true,
			AiInput: true,
		},
	}
}
