package lib

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

func ResetPasswordTemplate(token string) string {
	tpl := fmt.Sprintf(`
<!doctype html>
<html>
  <body>
    <div
      style='background-color:#F2F5F7;color:#242424;font-family:"Helvetica Neue", "Arial Nova", "Nimbus Sans", Arial, sans-serif;font-size:16px;font-weight:400;letter-spacing:0.15008px;line-height:1.5;margin:0;padding:32px 0;min-height:100%%;width:100%%'
    >
      <table
        align="center"
        width="100%%"
        style="margin:0 auto;max-width:600px;background-color:#FFFFFF"
        role="presentation"
        cellspacing="0"
        cellpadding="0"
        border="0"
      >
        <tbody>
          <tr style="width:100%%">
            <td>
              <div style="padding:24px 24px 8px 24px;text-align:left">
                <img
                  alt=""
                  src="https://d1iiu589g39o6c.cloudfront.net/live/platforms/platform_A9wwKSL6EV6orh6f/images/wptemplateimage_Xh1R23U9ziyct9nd/codoc.png"
                  height="24"
                  style="height:24px;outline:none;border:none;text-decoration:none;vertical-align:middle;display:inline-block;max-width:100%%"
                />
              </div>
			  <div style="font-weight:bold;text-align:left;margin:0;font-size:20px;padding:32px 24px 0px 24px">
				<h3>
					Reset your password?
				</h3>
			  </div>
              <div
                style="color:#474849;font-size:14px;font-weight:normal;text-align:left;padding:8px 24px 16px 24px"
              >
                If you didn&#x27;t request a reset, don&#x27;t worry. You can
                safely ignore this email.
              </div>
              <div style="text-align:left;padding:12px 24px 32px 24px">
                <a
                  href="http://localhost:5173/reset-password?token=%s"
                  style="color:#FFFFFF;font-size:14px;font-weight:bold;background-color:#0068FF;display:inline-block;padding:12px 20px;text-decoration:none"
                  target="_blank"
                  ><span
                    ><!--[if mso
                      ]><i
                        style="letter-spacing: 20px;mso-font-width:-100%%;mso-text-raise:30"
                        hidden
                        >&nbsp;</i
                      ><!
                    [endif]--></span
                  ><span>Change my password</span
                  ><span
                    ><!--[if mso
                      ]><i
                        style="letter-spacing: 20px;mso-font-width:-100%%"
                        hidden
                        >&nbsp;</i
                      ><!
                    [endif]--></span
                  ></a
                >
              </div>
              <div style="padding:16px 24px 16px 24px">
                <hr
                  style="width:100%%;border:none;border-top:1px solid #EEEEEE;margin:0"
                />
              </div>
              <div
                style="color:#474849;font-size:12px;font-weight:normal;text-align:left;padding:4px 24px 24px 24px"
              >
                Do not reply to this email.
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </body>
</html>
	`, token)
	return tpl
}

func SendMail(recipient string, title string, body string) {
	d := gomail.NewDialer("mailslurp", 2500, "user", "123456")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", "noreply@forgotpassword.local")
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", fmt.Sprintf("[NOREPLY] %s", title))
	m.SetBody("text/html", body)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
