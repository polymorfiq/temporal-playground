package specs

type RobotRules struct {
	UserAgent     string
	License       string
	Disallow      []string
	Allow         []string
	Sitemaps      []string
	CrawlDelay    uint8
	ContentSignal RobotContentSignal
}

type RobotContentSignal struct {
	AiTrain bool
	Search  bool
	AiInput bool
}
