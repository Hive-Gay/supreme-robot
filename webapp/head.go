package webapp

var HeadFrameworkCSSTemplate = []templateHeadLink{
	//{
	//	HRef: "/static/vendor/bootstrap-4.5.3-cyborg.min.css",
	//	Rel: "stylesheet",
	//	CrossOrigin: "anonymous",
	//},
	{
		HRef: "/static/vendor/hive-bootstrap.css",
		Rel: "stylesheet",
		CrossOrigin: "anonymous",
	},
	{
		HRef: "/static/vendor/fontawesome-free-5.15.1-web/css/all.min.css",
		Rel: "stylesheet",
		Integrity: "sha384-vp86vTRFVJgpjF9jiIGPEEqYqlDwgyBgEF109VFjmqGmIY/Y4HV4d3Gp2irVfcrp",
		CrossOrigin: "anonymous",
	},
}

var HeadCSSTemplate = []templateHeadLink{
	{
		HRef: "/static/css/default.css",
		Rel: "stylesheet",
	},
}
