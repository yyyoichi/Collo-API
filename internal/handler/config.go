package handler

import (
	"strings"
	"time"
	"yyyoichi/Collo-API/internal/analyzer"
	apiv3 "yyyoichi/Collo-API/internal/api/v3"
	"yyyoichi/Collo-API/internal/matrix"
	"yyyoichi/Collo-API/internal/ndl"
)

type Config struct {
	ndlConfig      ndl.Config
	analyzerConfig analyzer.Config
	matrixConfig   matrix.Config
}

func NewConfig(v3Config *apiv3.RequestConfig) Config {
	ndlConfig := ndl.Config{}
	l, _ := time.LoadLocation("Asia/Tokyo")
	ndlConfig.Search.Any = v3Config.Keyword
	ndlConfig.Search.From = v3Config.From.AsTime().In(l)
	ndlConfig.Search.Until = v3Config.Until.AsTime().In(l)
	ndlConfig.NDLAPI = ndl.NDLAPI(v3Config.PickGroupType)
	if v3Config.NdlApiType == 0 {
		ndlConfig.NDLAPI = ndl.SpeechAPI
	}
	ndlConfig.UseCache = v3Config.UseNdlCache
	ndlConfig.CreateCache = v3Config.CreateNdlCache
	ndlConfig.CacheDir = v3Config.NdlCacheDir

	analyzerConfig := analyzer.Config{}
	analyzerConfig.Includes = make([]analyzer.PartOfSpeechType, len(v3Config.PartOfSpeechTypes))
	for i, t := range v3Config.PartOfSpeechTypes {
		analyzerConfig.Includes[i] = analyzer.PartOfSpeechType(t)
	}
	analyzerConfig.StopWords = v3Config.Stopwords

	matrixConfig := matrix.Config{}
	switch v3Config.PickGroupType {
	case apiv3.RequestConfig_PICK_GROUP_TYPE_ISSUEID:
		matrixConfig.GroupingFuncType = matrix.PickByKey
	case apiv3.RequestConfig_PICK_GROUP_TYPE_MONTH:
		matrixConfig.GroupingFuncType = matrix.PickByMonth
	default:
		matrixConfig.GroupingFuncType = matrix.PickByKey
	}
	matrixConfig.ReduceThreshold = 0.05 // 5%の単語利用
	matrixConfig.AtGroupID = v3Config.ForcusGroupId

	config := Config{
		ndlConfig:      ndlConfig,
		analyzerConfig: analyzerConfig,
		matrixConfig:   matrixConfig,
	}
	return config
}

const (
	strNdlConfig      = "n!?:"
	strAnalyzerConfig = "a!?:"
	strMatrixConfig   = "m!?:"
)

func (c *Config) ToString() string {
	var buf strings.Builder

	buf.WriteString(strNdlConfig)
	buf.WriteString(c.ndlConfig.ToString())

	buf.WriteString(strAnalyzerConfig)
	buf.WriteString(c.analyzerConfig.ToString())

	buf.WriteString(strMatrixConfig)
	buf.WriteString(c.matrixConfig.ToString())

	return buf.String()
}
