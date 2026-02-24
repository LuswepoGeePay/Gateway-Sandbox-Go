package services

const (
	VerificationEmailTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Verification Code</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      background-color: #f7f7f7;
      margin: 0;
      padding: 0;
      color: #333;
      line-height: 1.6;
    }
    .container {
      max-width: 600px;
      margin: 20px auto;
      background-color: #ffffff;
      border-radius: 8px;
      padding: 30px;
      box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    }
    .header {
      text-align: center;
      margin-bottom: 25px;
    }
    .header h1 {
      color: #2c3e50;
      font-size: 24px;
      margin-bottom: 10px;
    }
    .code-container {
      background-color: #f2f6fc;
      border-radius: 6px;
      padding: 20px;
      text-align: center;
      margin: 25px 0;
    }
    .code {
      font-size: 32px;
      font-weight: bold;
      letter-spacing: 5px;
      color: #3498db;
    }
    .instructions {
      margin: 25px 0;
      font-size: 16px;
    }
    .footer {
      margin-top: 30px;
      font-size: 14px;
      color: #7f8c8d;
      text-align: center;
      border-top: 1px solid #eee;
      padding-top: 20px;
    }
    .secondary {
      color: #7f8c8d;
      font-size: 14px;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>Your Verification Code</h1>
      <p>Please use the code below to complete your verification</p>
    </div>
    
    <div class="code-container">
      <div class="code">{{ .Code }}</div>
    </div>
    
    <div class="instructions">
      <p>This code will expire in <strong>10 minutes</strong>.</p>
      <p>If you didn't request this code, please ignore this message or contact support if you're concerned about your account security.</p>
    </div>
    
    <div class="footer">
      <p class="secondary">This is an automated message, please do not reply.</p>
      <p>© 2025 GeePay Sandbox. All rights reserved.</p>
    </div>
  </div>
</body>
</html>`

	ResetPasswordEmailTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Reset Password</title>
    <style>
        body { font-family: Arial, sans-serif; background-color: #f7f7f7; margin: 0; padding: 0; color: #333; line-height: 1.6; }
        .container { max-width: 600px; margin: 20px auto; background-color: #ffffff; border-radius: 8px; padding: 30px; box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1); }
        .header { text-align: center; margin-bottom: 25px; }
        .code-container { background-color: #f2f6fc; border-radius: 6px; padding: 20px; text-align: center; margin: 25px 0; }
        .code { font-size: 32px; font-weight: bold; letter-spacing: 5px; color: #e74c3c; }
        .footer { margin-top: 30px; font-size: 14px; color: #7f8c8d; text-align: center; border-top: 1px solid #eee; padding-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset Request</h1>
            <p>You requested to reset your password. Use the code below:</p>
        </div>
        <div class="code-container">
            <div class="code">{{ .Code }}</div>
        </div>
        <p>If you did not request this, please ignore this email.</p>
        <div class="footer">
            <p>© 2025 GeePay Sandbox. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

	CardOTPTemplate = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Payment Verification</title>
    <style>
      body { font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0; }
      .container { max-width: 600px; margin: 20px auto; background-color: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
      h1 { color: #333; text-align: center; }
      .code { background: #eef2f7; padding: 15px; text-align: center; border-radius: 5px; margin: 20px 0; }
      h2 { color: #007bff; margin: 0; letter-spacing: 5px; font-size: 32px; }
      p { color: #666; line-height: 1.5; text-align: center; }
    </style>
  </head>
  <body>
    <div class="container">
      <h1>Payment OTP Code</h1>
      <p>This code is for simulation purposes only.</p>
      <div class="code">
        <h2>{{ .Code }}</h2>
      </div>
      <p>If you did not request this code, ignore this email.</p>
    </div>
  </body>
</html>`
)
