package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func goodStream(stream chan interface{}) bool {
	for rsp := range stream {
		switch rsp.(type) {
		case error:
			return false
		}
	}
	return true
}

func TestExecutorTrainAndPredict(t *testing.T) {
	a := assert.New(t)
	a.NotPanics(func() {
		pr, e := newParser().Parse(testTrainSelectIris)
		a.NoError(e)
		stream := runExtendedSQL(testTrainSelectIris, testDB, testCfg, pr)
		a.True(goodStream(stream))

		pr, e = newParser().Parse(testPredictSelectIris)
		a.NoError(e)
		stream = runExtendedSQL(testPredictSelectIris, testDB, testCfg, pr)
		a.True(goodStream(stream))
	})
}

func TestStandardSQL(t *testing.T) {
	a := assert.New(t)
	a.NotPanics(func() {
		stream := runStandardSQL(testSelectIris, testDB)
		a.True(goodStream(stream))
	})
	a.NotPanics(func() {
		stream := runStandardSQL(testStandardExecutiveSQLStatement, testDB)
		a.True(goodStream(stream))
	})
}

func TestCreatePredictionTable(t *testing.T) {
	a := assert.New(t)
	trainParsed, e := newParser().Parse(testTrainSelectIris)
	a.NoError(e)
	predParsed, e := newParser().Parse(testPredictSelectIris)
	a.NoError(e)
	a.NoError(createPredictionTable(trainParsed, predParsed, testDB))
}

func TestLogChanWriter_Write(t *testing.T) {
	a := assert.New(t)

	c := make(chan interface{})

	go func() {
		defer close(c)
		cw := &logChanWriter{c: c}
		cw.Write([]byte("hello\n世界"))
		cw.Write([]byte("hello\n世界"))
		cw.Write([]byte("\n"))
		cw.Write([]byte("世界\n世界\n世界\n"))
	}()

	a.Equal("hello\n", <-c)
	a.Equal("世界hello\n", <-c)
	a.Equal("世界\n", <-c)
	a.Equal("世界\n", <-c)
	a.Equal("世界\n", <-c)
	a.Equal("世界\n", <-c)
	_, more := <-c
	a.False(more)
}
