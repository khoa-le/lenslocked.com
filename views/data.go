package views

const (
	AlertLevelError = "danger"
	AlertLevelWarning = "warning"
	AlertLevelInfo = "info"
	AlertLevelSuccess = "success"

	//AlertMessageGeneric is display any user when error is encountered
	AlertMessageGeneric = "Something went wrong, please try again and contact us if the problem persits"
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