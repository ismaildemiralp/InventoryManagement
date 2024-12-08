# Envanter Yönetim Sistemi

Bu proje, **Google Sheets** entegrasyonu ile **dinamik QR kod oluşturma** ve **envanter takibi** sağlayan bir sistemdir.

## Özellikler

### 📋 Envanter Yönetimi
- Google Sheets üzerinden envanter bilgilerini (Varlık Numarası, Tür, Kullanıcı, Garanti Bilgileri vb.) okur.
- Envanter bilgilerindeki güncellemeleri algılar ve dinamik olarak QR kod içeriklerini günceller.

### 🖨 QR Kod Üretimi
- Her cihaz için benzersiz bir QR kod oluşturur.
- QR kodlar cihaz bilgilerini URL üzerinden erişilebilir hale getirir.
- Güncellenen envanter bilgileri, QR kodun URL'sinde otomatik olarak güncellenir.

### 🌐 Web Sunucusu
- Envanter bilgilerini görüntülemek için bir web arayüzü sağlar.
- QR kod tarandığında cihaz detayları, kullanıcı dostu bir HTML sayfasında gösterilir.

### 🖥 Masaüstü Uygulaması
- Qt ile geliştirilmiş bir arayüz sunar.
- İki temel işlev:
    - **QR Kod Üret:** Tüm cihazlar için QR kodları toplu olarak oluşturur.
    - **Envanter Verilerini Güncelle:** Google Sheets'teki değişiklikleri senkronize eder ve QR kodları yeniden oluşturur.

## Nasıl Çalışır?

1. **Google Sheets ile Entegrasyon:**  
   Google Cloud Console üzerinden bir kimlik doğrulama JSON dosyası oluşturun ve projeye entegre edin.

2. **QR Kod Üretimi:**  
   Her cihaz için QR kodlar oluşturulur ve kaydedilir. QR kodlar tarandığında cihaz bilgilerine URL üzerinden erişilebilir olur.

3. **Web Sunucusu:**  
   Web sunucusu başlatılır ve cihaz bilgileri **http://0.0.0.0:8080/devices/{Varlık Numarası}** adresinde görüntülenebilir.

## Gereksinimler

- Go 1.19+
- Qt kütüphanesi (`github.com/therecipe/qt`)
- Google Sheets API (`google.golang.org/api/sheets/v4`)
- QR kod kütüphanesi (`github.com/skip2/go-qrcode`)

## Katkıda Bulunma

- Projeye katkıda bulunmak için bir **pull request** oluşturabilirsiniz.
- Geri bildirimleriniz için teşekkür ederiz! 😊
