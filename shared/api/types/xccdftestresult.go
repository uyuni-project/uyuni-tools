package types

type XccdfTestResult struct {
	Xid              int    `mapstructure:"xid"`
	Sid              int    `mapstructure:"sid"`
	ActionId         int    `mapstructure:"action_id"`
	Path             string `mapstructure:"path"`
	Ovalfiles        string `mapstructure:"ovalfiles"`
	OscapParameters  string `mapstructure:"oscap_parameters"`
	TestResult       string `mapstructure:"test_result"`
	Benchmark        string `mapstructure:"benchmark"`
	BenchmarkVersion string `mapstructure:"benchmark_version"`
	Profile          string `mapstructure:"profile"`
	ProfileTitle     string `mapstructure:"profile_title"`
	StartTime        string `mapstructure:"start_time"`
	EndTime          string `mapstructure:"end_time"`
	Errors           string `mapstructure:"errors"`
	Deletable        bool   `mapstructure:"deletable"`
}
