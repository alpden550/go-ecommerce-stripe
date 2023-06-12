package renders

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap        map[string]string
	IntMap           map[string]int
	FloatMap         map[string]float32
	Data             map[string]interface{}
	Widgets          map[string]interface{}
	CSRFToken        string
	Flash            string
	Warning          string
	Error            string
	IsAuthenticated  int
	UserID           int
	API              string
	CSSVersion       string
	StripeSecret     string
	StripePublishKey string
}
