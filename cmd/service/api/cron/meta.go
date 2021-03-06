package cron

import (
	"strconv"

	"github.com/fidelfly/gox/cronx"
)

const (
	MetaJobId        = "meta.job.id"
	MetaJobCode      = "meta.job.code"
	MetaJobType      = "meta.job.type"
	MetaJobRunWay    = "meta.job.run.way"
	MetaJobTraceable = "meta.job.traceable"
)

func GetJobId(md *cronx.Metadata) int64 {
	if v, ok := md.Get(MetaJobId); ok {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			return id
		}
	}
	return 0
}

func GetJobCode(md *cronx.Metadata) string {
	if v, ok := md.Get(MetaJobCode); ok {
		return v
	}
	return ""
}

func GetJobType(md *cronx.Metadata) string {
	if v, ok := md.Get(MetaJobType); ok {
		return v
	}
	return ""
}

func IsTraceable(md *cronx.Metadata) bool {
	if v, ok := md.Get(MetaJobTraceable); ok {
		if tv, err := strconv.ParseBool(v); err == nil {
			return tv
		}
	}
	return false
}
