package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/skip2/go-qrcode"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Computer struct {
	AssetNumber     string
	AssetType       string
	ComputerDetails string
	AssignedUser    string
	WarrantyDetails string
	IP_Address      string
	Location        string
	Host            string
}

var computers map[string]Computer
var mu sync.RWMutex

func readGoogleSheet(sheetID string, sheetRange string) (map[string]Computer, error) {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile("{JSON file contains the authentication credentials required to access the Google Sheets API.}")) // Create on Google Cloud Console
	if err != nil {
		log.Fatalf("Error connecting to Google Sheets: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(sheetID, sheetRange).Do()
	if err != nil {
		return nil, fmt.Errorf("Error reading data from sheet: %v", err)
	}

	computers := make(map[string]Computer)
	for _, row := range resp.Values {
		if len(row) < 8 {
			continue
		}
		computers[row[0].(string)] = Computer{
			AssetNumber:     row[0].(string),
			AssetType:       row[1].(string),
			ComputerDetails: row[2].(string),
			AssignedUser:    row[3].(string),
			WarrantyDetails: row[4].(string),
			IP_Address:      row[5].(string),
			Location:        row[6].(string),
			Host:            row[7].(string),
		}
	}
	return computers, nil
}

func updateQRCode(sheetID string, sheetRange string) error {
	computers, err := readGoogleSheet(sheetID, sheetRange)
	if err != nil {
		return fmt.Errorf("Error reading data from Google Sheets: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	for id, computer := range computers {
		qrData := fmt.Sprintf("http://{HOST IP}:8080/devices/%s", computer.AssetNumber)
		qrCodeFile := filepath.Join("QR Codes", fmt.Sprintf("qrcode_%s.png", id))

		err := generateQRCode(qrData, qrCodeFile)
		if err != nil {
			return fmt.Errorf("Error updating QR code for %s: %v", id, err)
		}
		fmt.Printf("QR code saved for %s as %s\n", id, qrCodeFile)
	}

	return nil
}

func generateQRCode(data string, filename string) error {
	qrSizeInPixels := 90

	qrCode, err := qrcode.Encode(data, qrcode.Medium, qrSizeInPixels)
	if err != nil {
		return fmt.Errorf("Error generating QR code: %v", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Error creating QR code file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(qrCode)
	if err != nil {
		return fmt.Errorf("Error writing QR code to file: %v", err)
	}

	return nil
}

func generateAllQRCodes(sheetID string, sheetRange string) error {
	computers, err := readGoogleSheet(sheetID, sheetRange)
	if err != nil {
		return err
	}

	outputDir := "QR Codes"
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error creating QR Codes directory: %v", err)
	}

	for id, computer := range computers {
		qrData := fmt.Sprintf("http://{HOSTIP}:8080/devices/%s", computer.AssetNumber)
		qrCodeFile := filepath.Join(outputDir, fmt.Sprintf("qrcode_%s.png", id))

		err := generateQRCode(qrData, qrCodeFile)
		if err != nil {
			return fmt.Errorf("Error generating QR code for %s: %v", id, err)
		}
		fmt.Printf("QR code saved for %s as %s\n", id, qrCodeFile)
	}

	return nil
}

func startWebServer() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	http.HandleFunc("/devices/", func(w http.ResponseWriter, r *http.Request) {
		deviceID := r.URL.Path[len("/devices/"):]

		mu.RLock()
		computer, exists := computers[deviceID]
		mu.RUnlock()

		if !exists {
			http.NotFound(w, r)
			return
		}

		var sb strings.Builder
		sb.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n")
		sb.WriteString("<meta charset=\"UTF-8\">\n")
		sb.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
		sb.WriteString("<title>Device Information</title>\n")
		sb.WriteString("<style>\n")
		sb.WriteString("body { font-family: Arial, sans-serif; background-color: #003366; color: white; margin: 0; padding: 0; }\n")
		sb.WriteString(".logo-container {\n  display: flex;\n  justify-content: center;\n  align-items: center;\n  margin-top: 20px;\n}\n")
		sb.WriteString(".logo {\n  max-width: 200px;\n  max-height: 200px;\n  width: auto;\n  height: auto;\n}\n")
		sb.WriteString("table {\n  width: 100%;\n  border-collapse: collapse;\n  margin-top: 35px;\n}\n")
		sb.WriteString("th, td {\n  padding: 8px;\n  text-align: left;\n  border: 1px solid #ddd;\n}\n")
		sb.WriteString("th {\n  background-color: #00509E;\n  color: white;\n}\n")
		sb.WriteString("td {\n  background-color: transparent;\n  color: white;\n}\n")
		sb.WriteString("</style>\n</head>\n<body>\n")

		sb.WriteString("<div class=\"logo-container\">\n")
		sb.WriteString("<img src=\"company-logo\" alt=\"Logo\" class=\"logo\">\n")
		sb.WriteString("</div>\n")

		sb.WriteString("<table>\n")
		sb.WriteString(fmt.Sprintf("<tr><th>Inventory ID</th><td>%s</td></tr>\n", html.EscapeString(computer.AssetNumber)))
		sb.WriteString(fmt.Sprintf("<tr><th>Type</th><td>%s</td></tr>\n", html.EscapeString(computer.AssetType)))
		sb.WriteString(fmt.Sprintf("<tr><th>Device Information</th><td>%s</td></tr>\n", html.EscapeString(computer.ComputerDetails)))
		sb.WriteString(fmt.Sprintf("<tr><th>Assigned User</th><td>%s</td></tr>\n", html.EscapeString(computer.AssignedUser)))
		sb.WriteString(fmt.Sprintf("<tr><th>Warranty Expiration Date</th><td>%s</td></tr>\n", html.EscapeString(computer.WarrantyDetails)))
		sb.WriteString(fmt.Sprintf("<tr><th>IP Address</th><td>%s</td></tr>\n", html.EscapeString(computer.IP_Address)))
		sb.WriteString(fmt.Sprintf("<tr><th>Location</th><td>%s</td></tr>\n", html.EscapeString(computer.Location)))
		sb.WriteString(fmt.Sprintf("<tr><th>Host</th><td>%s</td></tr>\n", html.EscapeString(computer.Host)))
		sb.WriteString("</table>\n</body>\n</html>\n")

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(sb.String()))
	})

	go func() {
		log.Println("Web server is starting: http://0.0.0.0:8080")
		log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
	}()
}

func main() {
	var err error
	computers, err = readGoogleSheet("{Google Sheets ID}", "{Sheets Name}!{Range of Columns}")
	if err != nil {
		log.Fatalf("Failed to read data: %v", err)
	}

	startWebServer()

	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("QR Code Generator")
	window.SetFixedSize2(400, 150)

	produceButton := widgets.NewQPushButton2("Generate QR Codes", nil)
	produceButton.SetFixedWidth(200)

	updateButton := widgets.NewQPushButton2("Update Inventory Data", nil)
	updateButton.SetFixedWidth(200)

	produceButton.ConnectClicked(func(bool) {
		err := generateAllQRCodes("{Google Sheets ID}", "{Sheets Name}!{Range of Columns}")
		if err != nil {
			msgBox := widgets.NewQMessageBox2(
				widgets.QMessageBox__Critical,
				"Error",
				"Error generating QR codes: "+err.Error(),
				widgets.QMessageBox__Ok,
				window,
				core.Qt__WindowType(core.Qt__Dialog),
			)
			msgBox.Exec()
		} else {
			msgBox := widgets.NewQMessageBox2(
				widgets.QMessageBox__Information,
				"Success",
				"QR Codes Generated!",
				widgets.QMessageBox__Ok,
				window,
				core.Qt__WindowType(core.Qt__Dialog),
			)
			msgBox.Exec()
		}
	})

	updateButton.ConnectClicked(func(bool) {
		err := updateQRCode("{Google Sheets ID}", "{Sheets Name}!{Range of Columns}") // Example : "1gsdgqweqw-dsafasgh-exx" ,	"Assets!E:K"
		if err != nil {
			msgBox := widgets.NewQMessageBox2(
				widgets.QMessageBox__Critical,
				"Error",
				"Error updating inventory data: "+err.Error(),
				widgets.QMessageBox__Ok,
				window,
				core.Qt__WindowType(core.Qt__Dialog),
			)
			msgBox.Exec()
		} else {
			msgBox := widgets.NewQMessageBox2(
				widgets.QMessageBox__Information,
				"Success",
				"Inventory Data is Updated!",
				widgets.QMessageBox__Ok,
				window,
				core.Qt__WindowType(core.Qt__Dialog),
			)
			msgBox.Exec()
		}
	})

	layout := widgets.NewQVBoxLayout()
	layout.AddWidget(produceButton, 0, core.Qt__AlignCenter)
	layout.AddWidget(updateButton, 0, core.Qt__AlignCenter)

	centralWidget := widgets.NewQWidget(nil, 0)
	centralWidget.SetLayout(layout)
	window.SetCentralWidget(centralWidget)

	window.Show()

	app.Exec()
}
