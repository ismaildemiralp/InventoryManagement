
# Inventory Management System

This project is a system that enables **dynamic QR code generation** and **inventory tracking** with **Google Sheets** integration.

## Features

### 📋 Inventory Management
- Reads inventory data (Asset Number, Type, User, Warranty Info, etc.) from Google Sheets.
- Detects updates in inventory data and dynamically updates the QR code content accordingly.

### 🖨 QR Code Generation
- Generates a unique QR code for each device.
- QR codes make device information accessible via a URL.
- Updated inventory information is automatically reflected in the QR code's URL.

### 🌐 Web Server
- Provides a web interface to view inventory information.
- When a QR code is scanned, device details are displayed on a user-friendly HTML page.

### 🖥 Desktop Application
- Provides an interface developed with Qt.
- Two main functions:
    - **Generate QR Codes:** Creates QR codes for all devices in bulk.
    - **Update Inventory Data:** Syncs changes from Google Sheets and regenerates QR codes.

## How It Works

1. **Google Sheets Integration:**  
   Create an authentication JSON file via Google Cloud Console and integrate it into the project.

2. **QR Code Generation:**  
   Generates and saves QR codes for each device. When scanned, the QR code redirects to the device’s information page via URL.

3. **Web Server:**  
   Starts a web server and displays device information at  
   **http://hostIP:8080/devices/{Asset Number}**.

## Requirements

- Go 1.19+
- Qt library (`github.com/therecipe/qt`)
- Google Sheets API (`google.golang.org/api/sheets/v4`)
- QR code library (`github.com/skip2/go-qrcode`)
