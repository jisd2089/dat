package taifinance

/**
    Author: luzequan
    Created: 2018-09-03 14:34:52
*/
const (
	TAIFIN_URL_TEST                = "http://test.91zhengxin.com/jyzx/zxservice.do"
	TAIFIN_URL                     = "http://service.91zhengxin.com/qdDcenter/getConsume"
	TAIFIN_CLIKEY                  = "C99BC53C32584933AB080B86C5318CE7"
	TFAPI_KEY                      = "E7ZTNJSFQcyV11HhbbtJilzeZ3HOdC"
	TFAPI_UPLOAD_URL               = "http://timodeller.taifinance.cn:8080/api/file-upload"
	TFAPI_PROCESSED_DATASETS_URL   = "http://timodeller.taifinance.cn:8080/api/processed-datasets"
	TFAPI_PREDICT_CREDIT_SCORE_URL = "http://timodeller.taifinance.cn:8080/api/predict-credit-score"
	TFAPI_PREDICT_CREDIT_SCORE_CARD_URL = "http://timodeller.taifinance.cn:8080/api/ predict-credit-score-card"
)


type PredictCreditScoreReq struct {
	ModelUID        string `json:"modelUID"`
	InstancesAmount string `json:"instancesAmount"`
	InstancesArray  string `json:"instancesArray"`
}
