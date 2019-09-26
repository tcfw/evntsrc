package emails

const (
	globalTemplate = `
<div style="background: #fff; font-family: proxima nova, helvetica, arial, sans-serif; color: #334; font-size: 14px;">
	<div style="max-width: 620px; box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.07), 0px 5px 15px rgba(0, 0, 0, 0.15); border-radius: 4px; background: #FDFDFE; margin: 0 auto;">
		<img alt="Evntsrc.io Logo and email header" src="https://staging.evntsrc.io/static/email_header.png" style="max-width: calc(100% + 34px); height: auto; margin: -12px -17px; image-rendering: auto;"/>
		<div style="padding: 35px;">
			{{.}}
		</div>
		<div style="text-align: center; margin-top: 10px; padding-bottom: 20px; color: #aaa; font-size: 11px;">
			&copy; 2019 evntsrc.io &nbsp;|&nbsp; <a href="https://evntsrc.io/privacy" style="color: #aaa; text-decoration: none;">Privacy</a>
			&nbsp;|&nbsp; <a href="https://evntsrc.io/terms" style="color: #aaa; text-decoration: none;">Terms</a>
		</div>
	</div>
</div>
	`
)
