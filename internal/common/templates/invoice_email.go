package templates

import "fmt"

func BookingSuccessEmail(customerName, orderID string, ticketCount int, year int) string {
	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <title>Booking Confirmed</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                background-color: #f9f9f9;
                padding: 20px;
                margin: 0;
            }
            .container {
                max-width: 600px;
                background: #ffffff;
                padding: 30px;
                border-radius: 8px;
                margin: 0 auto;
                box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
            }
            .header {
                text-align: center;
                margin-bottom: 30px;
            }
            .header h1 {
                color: #28a745;
                font-size: 28px;
                margin: 0;
            }
            .content {
                color: #333;
                line-height: 1.6;
            }
            .info-box {
                background-color: #f8f9fa;
                padding: 20px;
                border-radius: 6px;
                margin: 20px 0;
                border-left: 4px solid #28a745;
            }
            .info-box p {
                margin: 8px 0;
                color: #333;
            }
            .next-steps {
                background-color: #e7f3ff;
                padding: 15px;
                border-radius: 6px;
                margin: 20px 0;
                border: 1px solid #b8daff;
                color: #004085;
                font-size: 14px;
            }
            .footer {
                text-align: center;
                margin-top: 30px;
                color: #6c757d;
            }
            .footer small {
                font-size: 12px;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1>Booking Confirmed!</h1>
            </div>
            <div class="content">
                <h2>Hi %s,</h2>
                <p>Thank you for your booking! Your payment has been successfully processed.</p>
                <div class="info-box">
                    <p><strong>Booking ID:</strong> <span style="color: #007bff;">%s</span></p>
                    <p><strong>Total Tickets:</strong> %d</p>
                    <p><strong>Status:</strong> <span style="color: #28a745; font-weight: bold;">Paid</span></p>
                </div>
                <div class="next-steps">
                    <strong>Next Steps:</strong><br/>
                    You can now check in at the port using your ticket. Please arrive at least 30 minutes before departure.
                </div>
                <p>Best regards,<br/>
                <strong>eTicket System Team</strong></p>
            </div>
            <hr style="border: none; height: 1px; background-color: #e9ecef; margin: 30px 0;"/>
            <div class="footer">
                <small>&copy; %d eTicket. All rights reserved.</small>
            </div>
        </div>
    </body>
    </html>
    `, customerName, orderID, ticketCount, year)
}

func BookingFailedEmail(customerName, orderID, reason string) string {
	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <title>Payment Failed</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                background-color: #f9f9f9;
                padding: 20px;
                margin: 0;
            }
            .container {
                max-width: 600px;
                background: #ffffff;
                padding: 30px;
                border-radius: 8px;
                margin: 0 auto;
                box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
            }
            .header {
                text-align: center;
                margin-bottom: 30px;
            }
            .header h1 {
                color: #dc3545;
                font-size: 28px;
                margin: 0;
            }
            .content {
                color: #333;
                line-height: 1.6;
            }
            .info-box {
                background-color: #f8f9fa;
                padding: 20px;
                border-radius: 6px;
                margin: 20px 0;
                border-left: 4px solid #dc3545;
            }
            .info-box p {
                margin: 8px 0;
                color: #333;
            }
            .next-steps {
                background-color: #fff3cd;
                padding: 15px;
                border-radius: 6px;
                margin: 20px 0;
                border: 1px solid #ffeaa7;
                color: #856404;
                font-size: 14px;
            }
            .footer {
                text-align: center;
                margin-top: 30px;
                color: #6c757d;
            }
            .footer small {
                font-size: 12px;
            }
            .action-button {
                background-color: #007bff;
                color: white;
                padding: 12px 24px;
                text-decoration: none;
                border-radius: 6px;
                font-weight: bold;
                display: inline-block;
                margin-top: 20px;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1>Payment Failed</h1>
            </div>
            <div class="content">
                <h2>Hi %s,</h2>
                <p>Unfortunately, your booking could not be confirmed due to a payment issue.</p>
                <div class="info-box">
                    <p><strong>Booking ID:</strong> <span style="color: #007bff;">%s</span></p>
                    <p><strong>Reason:</strong> <span style="color: #dc3545;">%s</span></p>
                    <p><strong>Status:</strong> <span style="color: #dc3545; font-weight: bold;">Failed</span></p>
                </div>
                <div class="next-steps">
                    <strong>What's Next:</strong><br/>
                    • You can try booking again<br/>
                    • Contact our support team if you need assistance<br/>
                    • No charges have been made to your account
                </div>
                <div style="text-align: center;">
                    <a href="#" class="action-button">Try Booking Again</a>
                </div>
                <p>We apologize for the inconvenience. If you continue to experience issues, please contact our support team.</p>
                <p>Best regards,<br/>
                <strong>eTicket System Team</strong></p>
            </div>
            <hr style="border: none; height: 1px; background-color: #e9ecef; margin: 30px 0;"/>
            <div class="footer">
                <small>&copy; 2025 eTicket. All rights reserved.</small>
            </div>
        </div>
    </body>
    </html>
    `, customerName, orderID, reason)
}
