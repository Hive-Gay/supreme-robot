package webapp

var HeadFrameworkCSSTemplate = []templateHeadLink{
	{
		HRef: "https://cdn.hive.gay/hive-bootstrap-develop.css",
		Rel: "stylesheet",
		CrossOrigin: "anonymous",
	},
	{
		HRef: "/static/vendor/fontawesome-free-5.15.1-web/css/all.min.css",
		Rel: "stylesheet",
		CrossOrigin: "anonymous",
	},
}

var HeadCSSTemplate = []templateHeadLink{
	{
		HRef: "/static/css/default.css",
		Rel: "stylesheet",
	},
}
