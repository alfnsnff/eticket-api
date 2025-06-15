package templates

import (
	"fmt"
)

func PasswordResetEmail(username, link string, time int) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
		<meta charset="UTF-8">
		<title>Password Reset</title>
		<style>
			body {
			font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
			background-color: #f4f4f4;
			margin: 0;
			padding: 0;
			}
			.container {
			max-width: 600px;
			margin: 40px auto;
			background-color: #ffffff;
			padding: 30px;
			border-radius: 8px;
			box-shadow: 0 0 10px rgba(0,0,0,0.05);
			}
			.header {
			text-align: center;
			padding-bottom: 20px;
			}
			.header img {
			width: 120px;
			}
			.title {
			font-size: 24px;
			color: #333333;
			margin-bottom: 10px;
			}
			.content {
			font-size: 16px;
			color: #555555;
			line-height: 1.5;
			margin-bottom: 30px;
			}
			.button {
			display: inline-block;
			padding: 12px 20px;
			background-color: #6441a5;
			color: #ffffff;
			text-decoration: none;
			border-radius: 6px;
			font-weight: bold;
			}
			.footer {
			font-size: 13px;
			color: #888888;
			text-align: center;
			margin-top: 30px;
			}
		</style>
		</head>
		<body>
		<div class="container">
			<div class="header">
			<img src="https://upload.wikimedia.org/wikipedia/commons/1/13/Ticket_emoji.png" alt="Logo" />
			</div>
			<div class="title">Reset Your Password</div>
			Hello %s,<br><br>
			We received a request to reset your password. If you made this request, click the button below. This link will expire in 15 minutes.<br><br>
			<a href="%s" class="button">Reset Password</a><br><br>
			If you did not request a password reset, please ignore this email or contact support if you have concerns.
			</div>
			<div class="footer">
			&copy; %d eTicket. All rights reserved.
			</div>
		</div>
		</body>
		</html>
		`, username, link, time)
}
