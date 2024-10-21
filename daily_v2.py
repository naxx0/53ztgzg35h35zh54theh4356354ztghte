import subprocess
import smtplib
import os
import configparser
from email.mime.base import MIMEBase
from email import encoders
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText

# Load configuration from config.ini
config = configparser.ConfigParser()
config.read('config.ini')

# Get SMTP and system configuration from the file
smtp_server = config['smtp']['smtp_server']
smtp_port = config['smtp'].getint('smtp_port')
smtp_username = config['smtp']['smtp_username']
smtp_password = config['smtp']['smtp_password']
sender_email = config['smtp']['sender_email']
recipient_email = config['smtp']['recipient_email']
systemname = config['system']['systemname']

# Define the filename for the text file
filename = 'results.txt'

# Run 'ifconfig' and 'netstat -antp' commands and save the outputs to the text file
with open(filename, 'w') as f:
    # Run 'ifconfig' command
    ifconfig_result = subprocess.run(['ifconfig'], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    f.write('ifconfig output:\n')
    f.write(ifconfig_result.stdout)
    f.write('\n\n')

    # Run 'netstat -antp' command
    netstat_result = subprocess.run(['netstat', '-antp'], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    f.write('netstat -antp output:\n')
    f.write(netstat_result.stdout)

# Check if the string '3389' exists in the text file
with open(filename, 'r') as f:
    content = f.read()

string_exists = '3389' in content

# Email subject based on whether '3389' exists, using the systemname variable
if string_exists:
    subject = f'{systemname} is up and running'
    body = 'none'
else:
    subject = f'{systemname} is NOT running!!!'
    body = 'none'

# Create a multipart email message
msg = MIMEMultipart()
msg['From'] = sender_email
msg['To'] = recipient_email
msg['Subject'] = subject

# Attach the email body
msg.attach(MIMEText(body, 'plain'))

# Attach the text file
with open(filename, 'rb') as attachment:
    part = MIMEBase('application', 'octet-stream')
    part.set_payload(attachment.read())

# Encode the attachment in base64
encoders.encode_base64(part)

# Add header to the attachment
part.add_header('Content-Disposition', f'attachment; filename= {filename}')

# Attach the attachment to the email message
msg.attach(part)

# Send the email
try:
    # Establish a secure session with the SMTP server
    with smtplib.SMTP(smtp_server, smtp_port) as server:
        server.starttls()  # Secure the connection
        server.login(smtp_username, smtp_password)
        server.send_message(msg)
    print('Email sent successfully.')
except Exception as e:
    print(f'Failed to send email: {e}')

# Remove the text file
try:
    os.remove(filename)
    print(f'{filename} has been deleted.')
except Exception as e:
    print(f'Failed to delete {filename}: {e}')
