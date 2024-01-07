package handler

import (
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

	analyzerConfig := analyzer.Config{}
	analyzerConfig.Includes = make([]analyzer.PartOfSpeechType, len(v3Config.PartOfSpeechTypes))
	for i, t := range v3Config.PartOfSpeechTypes {
		analyzerConfig.Includes[i] = analyzer.PartOfSpeechType(t)
	}
	analyzerConfig.StopWords = v3Config.Stopwords

	matrixConfig := matrix.Config{}
	if v3Config.PickGroupType <= 2 {
		matrixConfig.PickDocGroupID = func(d *matrix.Document) string { return d.Key }
	} else {
		matrixConfig.PickDocGroupID = func(d *matrix.Document) string {
			return d.At.Format("2006-01-02")
		}
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
