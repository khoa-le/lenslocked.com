package views

const (
	AlertLevelError = "danger"
	AlertLevelWarning = "warning"
	AlertLevelInfo = "info"
	AlertLevelSuccess = "success"
)

//Alert is used to render Bootstrap Alert message in template
//bootstrap.gohtml template
type Alert struct{
	Level string
	Message string
}


//Data is the top level structure that view expect data 
//to come in
type Data struct{
	Alert *Alert
	Yield interface{}
}