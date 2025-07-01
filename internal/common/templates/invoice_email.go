package templates

import (
	"eticket-api/internal/domain"
	"fmt"
	"strings"
	"time"
)

func BookingInvoiceEmail(booking *domain.Booking, payment *domain.Transaction) string {
	// QR logic: show QR from URL if available, else from base64 string, else show unavailable
	qrImgHTML := ""
	if payment.QrUrl != nil {
		qrImgHTML = fmt.Sprintf(`<img src="%s" alt="QRIS Pembayaran">`, *payment.QrUrl)
	} else if payment.QrString != nil {
		qrImgHTML = fmt.Sprintf(`<img src="data:image/png;base64,%s" alt="QRIS Pembayaran">`, *payment.QrString)
	} else {
		qrImgHTML = `<div style="color:#dc3545;">QR tidak tersedia</div>`
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Tagihan Pembayaran Tiket - Tiket Hebat</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background: #f4f7fa; padding: 20px; }
        .email-container { max-width: 650px; margin: 0 auto; background: #fff; border-radius: 20px; box-shadow: 0 8px 24px rgba(0,0,0,0.07); overflow: hidden; }
        .header { background: linear-gradient(135deg, #007bff 0%%, #00c6ff 100%%); color: white; padding: 40px 30px; text-align: center; }
        .header h1 { font-size: 28px; font-weight: 700; margin-bottom: 10px; }
        .header p { font-size: 16px; opacity: 0.95; }
        .content { padding: 40px 30px; }
        .greeting { font-size: 20px; color: #333; margin-bottom: 20px; font-weight: 600; }
        .invoice-summary { background: #f1f8ff; border-radius: 15px; padding: 25px; margin: 25px 0; border-left: 5px solid #007bff; }
        .info-label { font-size: 14px; color: #666; text-transform: uppercase; margin-bottom: 5px; font-weight: 500; }
        .info-value { font-size: 18px; color: #333; font-weight: 700; }
        .payment-instruction { background: #fff3cd; border: 1px solid #ffd43b; border-radius: 12px; padding: 20px; margin: 25px 0; }
        .payment-instruction h3 { color: #856404; margin-bottom: 15px; }
        .qr-section { text-align: center; margin: 30px 0; }
        .qr-section img { width: 180px; height: 180px; }
        .footer { background: #343a40; color: white; padding: 30px; text-align: center; }
        @media (max-width: 600px) {
            .content { padding: 20px 10px; }
            .header { padding: 30px 10px; }
        }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="header">
            <h1>Tagihan Pembayaran Tiket</h1>
            <p>Segera selesaikan pembayaran Anda untuk mengamankan tiket</p>
        </div>
        <div class="content">
            <div class="greeting">
                Halo %s! 👋
            </div>
            <p style="margin-bottom: 25px; font-size: 16px; color: #555;">
                Terima kasih telah melakukan pemesanan tiket kapal di <strong>Tiket Hebat</strong>.<br>
                Berikut adalah detail tagihan dan instruksi pembayaran Anda:
            </p>
            <div class="invoice-summary">
                <div class="info-label">ID Pemesanan</div>
                <div class="info-value">%s</div>
                <div class="info-label" style="margin-top:15px;">Status</div>
                <div class="info-value" style="color: #007bff;">Menunggu Pembayaran</div>
                <div class="info-label" style="margin-top:15px;">Total Tagihan</div>
                <div class="info-value" style="color: #28a745;">Rp %s</div>
                <div class="info-label" style="margin-top:15px;">Batas Waktu Pembayaran</div>
                <div class="info-value" style="color: #dc3545;">%s WIB</div>
            </div>
            <div class="payment-instruction">
                <h3>📋 Cara Pembayaran:</h3>
                <ul style="color: #856404; padding-left: 20px;">
                    <li>Gunakan QRIS di bawah ini atau transfer ke rekening yang tertera pada halaman pembayaran.</li>
                    <li>Pastikan pembayaran dilakukan sebelum batas waktu di atas.</li>
                    <li>Setelah pembayaran, tiket akan dikirim otomatis ke email Anda.</li>
                </ul>
            </div>
            <div class="qr-section">
                <div style="margin-bottom:10px;">Scan QRIS untuk membayar:</div>
                %s
                <div style="margin-top:10px; font-size:14px; color:#555;">Atau gunakan kode pembayaran: <b>%s</b></div>
            </div>
            <div style="margin-top: 30px; color: #555; font-size: 16px;">
                Jika Anda mengalami kendala pembayaran, silakan hubungi customer service kami.<br>
                <strong>Tim Tiket Hebat</strong>
            </div>
        </div>
        <div class="footer">
            <div style="max-width: 500px; margin: 0 auto;">
                <h3 style="margin-bottom: 15px;">Tiket Hebat</h3>
                <p style="margin-bottom: 15px; color: #adb5bd;">
                    Platform pemesanan tiket kapal terpercaya di Indonesia
                </p>
                <div style="margin: 20px 0;">
                    <a href="#" style="margin: 0 10px; color: #6c757d; text-decoration: none; font-size: 14px;">Website</a> |
                    <a href="#" style="margin: 0 10px; color: #6c757d; text-decoration: none; font-size: 14px;">Facebook</a> |
                    <a href="#" style="margin: 0 10px; color: #6c757d; text-decoration: none; font-size: 14px;">Instagram</a>
                </div>
                <div style="font-size: 12px; color: #6c757d; margin-top: 20px;">
                    &copy; %d Tiket Hebat. Semua hak dilindungi.
                </div>
            </div>
        </div>
    </div>
</body>
</html>
`,
		booking.CustomerName,
		booking.OrderID,
		formatPrice(float64(payment.Amount)),
		time.Unix(payment.ExpiredTime, 0).Format("2 January 2006 15:04"),
		qrImgHTML,
		payment.PayCode,
		time.Now().Year(),
	)
}

func BookingSuccessEmail(booking *domain.Booking, tickets []*domain.Ticket) string {
	// Format departure date
	departureDate := booking.Schedule.DepartureDatetime.Format("Monday, 2 January 2006")
	departureTime := booking.Schedule.DepartureDatetime.Format("15:04")

	// Count passenger and vehicle tickets
	passengerCount := 0
	vehicleCount := 0
	totalPrice := 0.0

	for _, ticket := range tickets {
		switch ticket.Type {
		case "passenger":
			passengerCount++
		case "vehicle":
			vehicleCount++
		}
		totalPrice += ticket.Price
	}

	// Build ticket details HTML
	ticketDetailsHTML := buildTicketDetailsHTML(tickets)

	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Tiket Berhasil Dipesan - Tiket Hebat</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            padding: 20px;
            line-height: 1.6;
        }
        
        .email-container {
            max-width: 650px;
            margin: 0 auto;
            background: #ffffff;
            border-radius: 20px;
            overflow: hidden;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
        }
        
        .header {
            background: linear-gradient(135deg, #4CAF50 0%%, #45a049 100%%);
            color: white;
            padding: 40px 30px;
            text-align: center;
            position: relative;
        }
        
        .header::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><circle cx="20" cy="20" r="2" fill="rgba(255,255,255,0.1)"/><circle cx="80" cy="30" r="1.5" fill="rgba(255,255,255,0.1)"/><circle cx="40" cy="70" r="1" fill="rgba(255,255,255,0.1)"/><circle cx="90" cy="80" r="2.5" fill="rgba(255,255,255,0.1)"/></svg>');
        }
        
        .success-icon {
            width: 80px;
            height: 80px;
            background: rgba(255, 255, 255, 0.2);
            border-radius: 50%%;
            display: flex;
            align-items: center;
            justify-content: center;
            margin: 0 auto 20px;
            font-size: 40px;
        }
        
        .header h1 {
            font-size: 32px;
            font-weight: 700;
            margin-bottom: 10px;
            position: relative;
            z-index: 1;
        }
        
        .header p {
            font-size: 18px;
            opacity: 0.9;
            position: relative;
            z-index: 1;
        }
        
        .content {
            padding: 40px 30px;
        }
        
        .greeting {
            font-size: 20px;
            color: #333;
            margin-bottom: 20px;
            font-weight: 600;
        }
        
        .booking-summary {
            background: linear-gradient(135deg, #f8f9ff 0%%, #e8f2ff 100%%);
            border-radius: 15px;
            padding: 25px;
            margin: 25px 0;
            border-left: 5px solid #4CAF50;
            position: relative;
            overflow: hidden;
        }
        
        .booking-summary::before {
            content: '';
            position: absolute;
            top: -50%%;
            right: -20px;
            width: 100px;
            height: 100px;
            background: radial-gradient(circle, rgba(76, 175, 80, 0.1) 0%%, transparent 70%%);
            border-radius: 50%%;
        }
        
        .booking-info {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin-bottom: 25px;
        }
        
        .info-item {
            position: relative;
            z-index: 1;
        }
        
        .info-label {
            font-size: 14px;
            color: #666;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            margin-bottom: 5px;
            font-weight: 500;
        }
        
        .info-value {
            font-size: 18px;
            color: #333;
            font-weight: 700;
        }
        
        .route-info {
            background: white;
            border-radius: 12px;
            padding: 20px;
            margin: 20px 0;
            box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
            border: 1px solid #e8f2ff;
        }
        
        .route-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 15px;
            flex-wrap: wrap;
            gap: 10px;
        }
        
        .route-path {
            font-size: 20px;
            font-weight: 700;
            color: #333;
        }
        
        .ship-name {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 8px 15px;
            border-radius: 25px;
            font-size: 14px;
            font-weight: 600;
        }
        
        .datetime-info {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 15px;
        }
        
        .datetime-item {
            text-align: center;
            padding: 10px;
            background: #f8f9ff;
            border-radius: 8px;
        }
        
        .tickets-section {
            margin: 30px 0;
        }
        
        .section-title {
            font-size: 22px;
            color: #333;
            margin-bottom: 20px;
            font-weight: 700;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        .section-title::before {
            content: '🎫';
            font-size: 24px;
        }
        
        .ticket-card {
            background: white;
            border: 2px solid #e8f2ff;
            border-radius: 12px;
            padding: 20px;
            margin-bottom: 15px;
            position: relative;
            overflow: hidden;
            transition: all 0.3s ease;
        }
        
        .ticket-card::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            width: 5px;
            height: 100%%;
            background: linear-gradient(135deg, #4CAF50 0%%, #45a049 100%%);
        }
        
        .ticket-card.vehicle::before {
            background: linear-gradient(135deg, #FF9800 0%%, #F57C00 100%%);
        }
        
        .ticket-header {
            display: flex;
            justify-content: between;
            align-items: center;
            margin-bottom: 15px;
            flex-wrap: wrap;
            gap: 10px;
        }
        
        .ticket-type {
            background: #4CAF50;
            color: white;
            padding: 6px 12px;
            border-radius: 20px;
            font-size: 12px;
            font-weight: 600;
            text-transform: uppercase;
        }
        
        .ticket-type.vehicle {
            background: #FF9800;
        }
        
        .ticket-code {
            font-family: 'Courier New', monospace;
            background: #f5f5f5;
            padding: 5px 10px;
            border-radius: 5px;
            font-size: 14px;
            font-weight: bold;
            color: #333;
        }
        
        .ticket-details {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 15px;
        }
        
        .price-summary {
            background: linear-gradient(135deg, #1e3c72 0%%, #2a5298 100%%);
            color: white;
            border-radius: 15px;
            padding: 25px;
            margin: 25px 0;
            text-align: center;
        }
        
        .total-price {
            font-size: 32px;
            font-weight: 700;
            margin-bottom: 10px;
        }
        
        .price-breakdown {
            font-size: 16px;
            opacity: 0.9;
        }
        
        .next-steps {
            background: linear-gradient(135deg, #fff3cd 0%%, #ffeaa7 100%%);
            border: 1px solid #ffd43b;
            border-radius: 12px;
            padding: 20px;
            margin: 25px 0;
        }
        
        .next-steps h3 {
            color: #856404;
            margin-bottom: 15px;
            font-size: 18px;
        }
        
        .steps-list {
            list-style: none;
            padding: 0;
        }
        
        .steps-list li {
            color: #856404;
            margin-bottom: 8px;
            padding-left: 25px;
            position: relative;
        }
        
        .steps-list li::before {
            content: '✓';
            position: absolute;
            left: 0;
            color: #28a745;
            font-weight: bold;
        }
        
        .contact-info {
            background: #f8f9fa;
            border-radius: 12px;
            padding: 20px;
            margin: 25px 0;
            text-align: center;
        }
        
        .footer {
            background: #343a40;
            color: white;
            padding: 30px;
            text-align: center;
        }
        
        .footer-content {
            max-width: 500px;
            margin: 0 auto;
        }
        
        .social-links {
            margin: 20px 0;
        }
        
        .social-links a {
            display: inline-block;
            margin: 0 10px;
            color: #6c757d;
            text-decoration: none;
            font-size: 14px;
        }
        
        .copyright {
            font-size: 12px;
            color: #6c757d;
            margin-top: 20px;
        }
        
        @media (max-width: 600px) {
            body {
                padding: 10px;
            }
            
            .content {
                padding: 20px 15px;
            }
            
            .header {
                padding: 30px 20px;
            }
            
            .header h1 {
                font-size: 24px;
            }
            
            .booking-info {
                grid-template-columns: 1fr;
            }
            
            .route-header {
                flex-direction: column;
                align-items: flex-start;
            }
            
            .datetime-info {
                grid-template-columns: 1fr;
            }
            
            .ticket-details {
                grid-template-columns: 1fr;
            }
            
            .total-price {
                font-size: 24px;
            }
        }
    </style>
</head>
<body>
    <div class="email-container">
        <!-- Header -->
        <div class="header">
            <div class="success-icon">✅</div>
            <h1>Tiket Berhasil Dipesan!</h1>
            <p>Pembayaran telah berhasil diproses</p>
        </div>
        
        <!-- Content -->
        <div class="content">
            <div class="greeting">
                Halo %s! 👋
            </div>
            
            <p style="margin-bottom: 25px; font-size: 16px; color: #555;">
                Terima kasih telah memesan tiket kapal dengan <strong>Tiket Hebat</strong>. 
                Pembayaran Anda telah berhasil diproses dan tiket Anda sudah siap digunakan.
            </p>
            
            <!-- Booking Summary -->
            <div class="booking-summary">
                <div class="booking-info">
                    <div class="info-item">
                        <div class="info-label">ID Pemesanan</div>
                        <div class="info-value">%s</div>
                    </div>
                    <div class="info-item">
                        <div class="info-label">Status</div>
                        <div class="info-value" style="color: #4CAF50;">✅ BERHASIL</div>
                    </div>
                    <div class="info-item">
                        <div class="info-label">Tanggal Pemesanan</div>
                        <div class="info-value">%s</div>
                    </div>
                    <div class="info-item">
                        <div class="info-label">Email</div>
                        <div class="info-value">%s</div>
                    </div>
                </div>
            </div>
            
            <!-- Route Information -->
            <div class="route-info">
                <div class="route-header">
                    <div class="route-path">%s → %s</div>
                    <div class="ship-name">🚢 %s</div>
                </div>
                <div class="datetime-info">
                    <div class="datetime-item">
                        <div class="info-label">📅 Tanggal Keberangkatan</div>
                        <div class="info-value">%s</div>
                    </div>
                    <div class="datetime-item">
                        <div class="info-label">🕐 Waktu Keberangkatan</div>
                        <div class="info-value">%s WIB</div>
                    </div>
                </div>
            </div>
            
            <!-- Tickets Details -->
            <div class="tickets-section">
                <div class="section-title">Detail Tiket</div>
                %s
            </div>
            
            <!-- Price Summary -->
            <div class="price-summary">
                <div class="total-price">Rp %s</div>
                <div class="price-breakdown">
                    Total: %d Penumpang + %d Kendaraan
                </div>
            </div>
            
            <!-- Next Steps -->
            <div class="next-steps">
                <h3>📋 Langkah Selanjutnya:</h3>
                <ul class="steps-list">
                    <li>Datang ke pelabuhan minimal 30 menit sebelum keberangkatan</li>
                    <li>Bawa dokumen identitas yang valid</li>
                    <li>Tunjukkan email ini atau kode tiket untuk check-in</li>
                    <li>Untuk kendaraan, siapkan STNK dan dokumen kendaraan</li>
                </ul>
            </div>
            
            <!-- Contact Info -->
            <div class="contact-info">
                <h3 style="color: #333; margin-bottom: 15px;">📞 Butuh Bantuan?</h3>
                <p style="margin-bottom: 10px; color: #555;">Tim customer service kami siap membantu Anda</p>
                <div style="margin-top: 15px;">
                    <strong style="color: #333;">📧 Email:</strong> support@tikethebat.live<br>
                    <strong style="color: #333;">📱 WhatsApp:</strong> +62 811-1234-5678<br>
                    <strong style="color: #333;">🕐 Jam Operasional:</strong> 24/7
                </div>
            </div>
            
            <p style="margin-top: 30px; color: #555; font-size: 16px;">
                Selamat menikmati perjalanan Anda! ⛵<br>
                <strong>Tim Tiket Hebat</strong>
            </p>
        </div>
        
        <!-- Footer -->
        <div class="footer">
            <div class="footer-content">
                <h3 style="margin-bottom: 15px;">Tiket Hebat</h3>
                <p style="margin-bottom: 15px; color: #adb5bd;">
                    Platform pemesanan tiket kapal terpercaya di Indonesia
                </p>
                <div class="social-links">
                    <a href="#">Website</a> |
                    <a href="#">Facebook</a> |
                    <a href="#">Instagram</a> |
                    <a href="#">Twitter</a>
                </div>
                <div class="copyright">
                    &copy; %d Tiket Hebat. Semua hak dilindungi.
                </div>
            </div>
        </div>
    </div>
</body>
</html>`,
		booking.CustomerName,
		booking.OrderID,
		booking.CreatedAt.Format("2 January 2006"),
		booking.Email,
		booking.Schedule.DepartureHarbor.HarborName,
		booking.Schedule.ArrivalHarbor.HarborName,
		booking.Schedule.Ship.ShipName,
		departureDate,
		departureTime,
		ticketDetailsHTML,
		formatPrice(totalPrice),
		passengerCount,
		vehicleCount,
		time.Now().Year(),
	)
}

func BookingFailedEmail(booking *domain.Booking, reason string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Pembayaran Gagal - Tiket Hebat</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            padding: 20px;
            line-height: 1.6;
        }
        
        .email-container {
            max-width: 650px;
            margin: 0 auto;
            background: #ffffff;
            border-radius: 20px;
            overflow: hidden;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
        }
        
        .header {
            background: linear-gradient(135deg, #dc3545 0%%, #c82333 100%%);
            color: white;
            padding: 40px 30px;
            text-align: center;
            position: relative;
        }
        
        .error-icon {
            width: 80px;
            height: 80px;
            background: rgba(255, 255, 255, 0.2);
            border-radius: 50%%;
            display: flex;
            align-items: center;
            justify-content: center;
            margin: 0 auto 20px;
            font-size: 40px;
        }
        
        .header h1 {
            font-size: 32px;
            font-weight: 700;
            margin-bottom: 10px;
        }
        
        .content {
            padding: 40px 30px;
        }
        
        .booking-summary {
            background: linear-gradient(135deg, #fff5f5 0%%, #fed7d7 100%%);
            border-radius: 15px;
            padding: 25px;
            margin: 25px 0;
            border-left: 5px solid #dc3545;
        }
        
        .retry-button {
            background: linear-gradient(135deg, #007bff 0%%, #0056b3 100%%);
            color: white;
            padding: 15px 30px;
            text-decoration: none;
            border-radius: 25px;
            font-weight: bold;
            display: inline-block;
            margin: 20px 0;
            text-align: center;
            transition: all 0.3s ease;
        }
        
        .retry-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(0, 123, 255, 0.3);
        }
        
        .support-info {
            background: #f8f9fa;
            border-radius: 12px;
            padding: 20px;
            margin: 25px 0;
            text-align: center;
        }
        
        .footer {
            background: #343a40;
            color: white;
            padding: 30px;
            text-align: center;
        }
        
        @media (max-width: 600px) {
            body { padding: 10px; }
            .content { padding: 20px 15px; }
            .header { padding: 30px 20px; }
            .header h1 { font-size: 24px; }
        }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="header">
            <div class="error-icon">❌</div>
            <h1>Pembayaran Gagal</h1>
            <p>Maaf, terjadi kendala dalam proses pembayaran</p>
        </div>
        
        <div class="content">
            <div style="font-size: 20px; color: #333; margin-bottom: 20px; font-weight: 600;">
                Halo %s! 👋
            </div>
            
            <p style="margin-bottom: 25px; font-size: 16px; color: #555;">
                Maaf, pemesanan tiket Anda tidak dapat diproses karena terjadi masalah dengan pembayaran.
            </p>
            
            <div class="booking-summary">
                <h3 style="color: #dc3545; margin-bottom: 15px;">📋 Detail Pemesanan</h3>
                <p><strong>ID Pemesanan:</strong> %s</p>
                <p><strong>Alasan Gagal:</strong> %s</p>
                <p><strong>Status:</strong> <span style="color: #dc3545; font-weight: bold;">❌ GAGAL</span></p>
            </div>
            
            <div style="text-align: center; margin: 30px 0;">
                <a href="https://tikethebat.live" class="retry-button">
                    🔄 Coba Pesan Lagi
                </a>
            </div>
            
            <div style="background: #fff3cd; border: 1px solid #ffd43b; border-radius: 12px; padding: 20px; margin: 25px 0;">
                <h3 style="color: #856404; margin-bottom: 15px;">ℹ Yang Perlu Anda Ketahui:</h3>
                <ul style="color: #856404; padding-left: 20px;">
                    <li style="margin-bottom: 8px;">Tidak ada biaya yang dikenakan ke rekening Anda</li>
                    <li style="margin-bottom: 8px;">Anda dapat mencoba memesan kembali kapan saja</li>
                    <li style="margin-bottom: 8px;">Pastikan saldo atau limit kartu mencukupi</li>
                    <li>Hubungi bank Anda jika masalah berlanjut</li>
                </ul>
            </div>
            
            <div class="support-info">
                <h3 style="color: #333; margin-bottom: 15px;">📞 Butuh Bantuan?</h3>
                <p style="margin-bottom: 15px; color: #555;">Tim customer service kami siap membantu Anda</p>
                <div style="margin-top: 15px;">
                    <strong style="color: #333;">📧 Email:</strong> support@tikethebat.live<br>
                    <strong style="color: #333;">📱 WhatsApp:</strong> +62 811-1234-5678<br>
                    <strong style="color: #333;">🕐 Jam Operasional:</strong> 24/7
                </div>
            </div>
            
            <p style="margin-top: 30px; color: #555; font-size: 16px;">
                Kami mohon maaf atas ketidaknyamanan ini. Silakan coba lagi atau hubungi tim support kami.<br>
                <strong>Tim Tiket Hebat</strong>
            </p>
        </div>
        
        <div class="footer">
            <div style="max-width: 500px; margin: 0 auto;">
                <h3 style="margin-bottom: 15px;">Tiket Hebat</h3>
                <p style="margin-bottom: 15px; color: #adb5bd;">
                    Platform pemesanan tiket kapal terpercaya di Indonesia
                </p>
                <div style="margin: 20px 0;">
                    <a href="#" style="margin: 0 10px; color: #6c757d; text-decoration: none; font-size: 14px;">Website</a> |
                    <a href="#" style="margin: 0 10px; color: #6c757d; text-decoration: none; font-size: 14px;">Facebook</a> |
                    <a href="#" style="margin: 0 10px; color: #6c757d; text-decoration: none; font-size: 14px;">Instagram</a>
                </div>
                <div style="font-size: 12px; color: #6c757d; margin-top: 20px;">
                    &copy; %d Tiket Hebat. Semua hak dilindungi.
                </div>
            </div>
        </div>
    </div>
</body>
</html>`,
		booking.CustomerName,
		booking.OrderID,
		reason,
		time.Now().Year(),
	)
}

// Helper function to build ticket details HTML
func buildTicketDetailsHTML(tickets []*domain.Ticket) string {
	var html strings.Builder

	for _, ticket := range tickets {
		typeClass := "passenger"
		typeIcon := "👤"
		typeName := "PENUMPANG"

		if ticket.Type == "vehicle" {
			typeClass = "vehicle"
			typeIcon = "🚗"
			typeName = "KENDARAAN"
		}

		html.WriteString(fmt.Sprintf(`
        <div class="ticket-card %s">
            <div class="ticket-header">
                <div class="ticket-type %s">%s %s</div>
                <div class="ticket-code">%s</div>
            </div>
            <div class="ticket-details">
                <div class="info-item">
                    <div class="info-label">Nama</div>
                    <div class="info-value">%s</div>
                </div>`,
			typeClass,
			typeClass,
			typeIcon,
			typeName,
			ticket.TicketCode,
			ticket.PassengerName,
		))

		switch ticket.Type {
		case "passenger":
			html.WriteString(fmt.Sprintf(`
                <div class="info-item">
                    <div class="info-label">Umur</div>
                    <div class="info-value">%d tahun</div>
                </div>
                <div class="info-item">
                    <div class="info-label">Jenis Kelamin</div>
                    <div class="info-value">%s</div>
                </div>`,
				ticket.PassengerAge,
				getGenderDisplay(*ticket.PassengerGender),
			))

			if ticket.SeatNumber != nil {
				html.WriteString(fmt.Sprintf(`
                <div class="info-item">
                    <div class="info-label">Nomor Kursi</div>
                    <div class="info-value">%s</div>
                </div>`,
					*ticket.SeatNumber,
				))
			}
		case "vehicle":
			if ticket.LicensePlate != nil {
				html.WriteString(fmt.Sprintf(`
                <div class="info-item">
                    <div class="info-label">Plat Nomor</div>
                    <div class="info-value">%s</div>
                </div>`,
					*ticket.LicensePlate,
				))
			}
		}

		html.WriteString(fmt.Sprintf(`
                <div class="info-item">
                    <div class="info-label">Kelas</div>
                    <div class="info-value">%s</div>
                </div>
                <div class="info-item">
                    <div class="info-label">Harga</div>
                    <div class="info-value">Rp %s</div>
                </div>
            </div>
        </div>`,
			ticket.Class.ClassName,
			formatPrice(ticket.Price),
		))
	}

	return html.String()
}

// Helper functions
func formatPrice(price float64) string {
	// Format price with thousand separators
	return fmt.Sprintf("%.0f", price)
}

func getGenderDisplay(gender string) string {
	switch strings.ToLower(gender) {
	case "l", "laki-laki", "male":
		return "Laki-laki"
	case "p", "perempuan", "female":
		return "Perempuan"
	default:
		return gender
	}
}
