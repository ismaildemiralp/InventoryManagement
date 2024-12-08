📋 Envanter Yönetimi

    Google Sheets üzerinden envanter bilgilerini (Varlık Numarası, Tür, Kullanıcı, Garanti Bilgileri vb.) okur.
    Envanter bilgilerindeki güncellemeleri algılar ve dinamik olarak QR kod içeriklerini günceller.

🖨 QR Kod Üretimi

    Her cihaz için benzersiz bir QR kod oluşturur.
    QR kodlar cihaz bilgilerini URL üzerinden erişilebilir hale getirir.
    Güncellenen envanter bilgileri, QR kodun URL'sinde otomatik olarak güncellenir.

🌐 Web Sunucusu

    Envanter bilgilerini görüntülemek için bir web arayüzü sağlar.
    QR kod tarandığında cihaz detayları, kullanıcı dostu bir HTML sayfasında gösterilir.

🖥 Masaüstü Uygulaması

    Qt ile geliştirilmiş bir arayüz sunar.
    İki temel işlev:
        QR Kod Üret: Tüm cihazlar için QR kodları toplu olarak oluşturur.
        Envanter Verilerini Güncelle: Google Sheets'teki değişiklikleri senkronize eder ve QR kodları yeniden oluşturur.