package log

type ITestLogger interface {
	Logf(format string, args ...interface{})
}

type TestWriter struct {
	TestLogger ITestLogger
}

func NewTestWriter(testLogger ITestLogger) *TestWriter {
	return &TestWriter{TestLogger: testLogger}
}

func (t *TestWriter) Write(p []byte) (n int, err error) {
	t.TestLogger.Logf(string(p))
	return len(p), nil
}
