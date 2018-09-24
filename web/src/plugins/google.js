window.gapiLoaded = () => {
	gapi.load('auth2', () => {
		gapi.auth2.init({
			client_id: process.env.VUE_APP_GAPI_CLIENT_ID,
			fetch_basic_profile: true
		}).then(app.gapiCallback);
	});
};

(function (d, s, id) {
	var js, fjs = d.getElementsByTagName(s)[0];
	if (d.getElementById(id)) {
		return;
	}
	js = d.createElement(s);
	js.id = id;
	js.src = "https://apis.google.com/js/platform.js?onload=gapiLoaded";
	fjs.parentNode.insertBefore(js, fjs);
}(document, 'script', 'google-jssdk'));