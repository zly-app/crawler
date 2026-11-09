package config

import "testing"

func TestFrameConfigCheckFillsDefaults(t *testing.T) {
	conf := newFrameConfig()

	if err := conf.Check(); err != nil {
		t.Fatalf("Check() 返回错误: %v", err)
	}

	if len(conf.QueueSuffixes) == 0 {
		t.Fatal("QueueSuffixes 未填充默认值")
	}
	if conf.SeedQueueSuffix == "" {
		t.Fatal("SeedQueueSuffix 未填充默认值")
	}
	if conf.ErrorSeedQueueSuffix == "" {
		t.Fatal("ErrorSeedQueueSuffix 未填充默认值")
	}
	if conf.ParserErrorSeedQueueSuffix == "" {
		t.Fatal("ParserErrorSeedQueueSuffix 未填充默认值")
	}
	if conf.SetSuffix == "" {
		t.Fatal("SetSuffix 未填充默认值")
	}
	if conf.RequestTimeout != DefaultFrameRequestTimeout {
		t.Fatalf("RequestTimeout = %d, want %d", conf.RequestTimeout, DefaultFrameRequestTimeout)
	}
	if conf.SeedProcessTimeout != defSeedProcessTimeout {
		t.Fatalf("SeedProcessTimeout = %d, want %d", conf.SeedProcessTimeout, defSeedProcessTimeout)
	}
	if conf.RequestMaxAttemptCount != defaultFrameRequestMaxAttemptCount {
		t.Fatalf("RequestMaxAttemptCount = %d, want %d", conf.RequestMaxAttemptCount, defaultFrameRequestMaxAttemptCount)
	}
}

func TestFrameConfigCheckKeepsUserValues(t *testing.T) {
	conf := FrameConfig{
		QueueSuffixes:          []string{":a", ":b"},
		SeedQueueSuffix:        ":custom-seed",
		ErrorSeedQueueSuffix:   ":custom-error",
		RequestTimeout:         1234,
		SeedProcessTimeout:     5678,
		RequestMaxAttemptCount: 9,
	}

	if err := conf.Check(); err != nil {
		t.Fatalf("Check() 返回错误: %v", err)
	}

	if len(conf.QueueSuffixes) != 2 || conf.QueueSuffixes[0] != ":a" {
		t.Fatalf("用户配置的 QueueSuffixes 被覆盖: %v", conf.QueueSuffixes)
	}
	if conf.SeedQueueSuffix != ":custom-seed" {
		t.Fatalf("用户配置的 SeedQueueSuffix 被覆盖: %s", conf.SeedQueueSuffix)
	}
	if conf.ErrorSeedQueueSuffix != ":custom-error" {
		t.Fatalf("用户配置的 ErrorSeedQueueSuffix 被覆盖: %s", conf.ErrorSeedQueueSuffix)
	}
	if conf.RequestTimeout != 1234 {
		t.Fatalf("用户配置的 RequestTimeout 被覆盖: %d", conf.RequestTimeout)
	}
	if conf.SeedProcessTimeout != 5678 {
		t.Fatalf("用户配置的 SeedProcessTimeout 被覆盖: %d", conf.SeedProcessTimeout)
	}
	if conf.RequestMaxAttemptCount != 9 {
		t.Fatalf("用户配置的 RequestMaxAttemptCount 被覆盖: %d", conf.RequestMaxAttemptCount)
	}
}
