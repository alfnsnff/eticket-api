package templates

import "fmt"

func BookingSuccessEmail(customerName, orderID string, ticketCount int, year int) string {
	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head><meta charset="UTF-8"><title>Booking Confirmed</title></head>
	<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px;">
		<div style="max-width: 600px; background: white; padding: 30px; border-radius: 8px;">
			<h2 style="color: #333;">Hi %s,</h2>
			<p>Thank you for your booking!</p>
			<p><strong>Booking ID:</strong> %s</p>
			<p><strong>Total Tickets:</strong> %d</p>
			<p><strong>Status:</strong> Paid</p>
			<p>You can now check in at the port using your ticket.</p>
			<p style="margin-top: 30px;">Regards,<br/>eTicket System</p>
			<hr/>
			<small>&copy; %d eTicket. All rights reserved.</small>
		</div>
	</body>
	</html>
	`, customerName, orderID, ticketCount, year)
}
